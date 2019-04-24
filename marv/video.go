package marv

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	height       = 360
	width        = 640
	raspividArgs = []string{
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),
		// "-n", // No Preview
		// "-ev", "10", // Stabilization
		// "-ex", "antishake" // Antishake
		"-rot", "180",
		"-fps", "48",
		"-t", "0",
		"-pf", "baseline",
		"-o", "-",
	}

	initialFrameCount = 8
	nalSeparator      = []byte{0x00, 0x00, 0x00, 0x01}
)

type clientMap struct {
	clients       map[*websocket.Conn]bool
	mutex         sync.RWMutex
	initialFrames [][]byte
}

func newClientMap(initialFrames [][]byte) *clientMap {
	return &clientMap{
		clients:       map[*websocket.Conn]bool{},
		initialFrames: initialFrames,
	}
}

// Clients NEEDSCOMMENT
func (c *clientMap) Clients() map[*websocket.Conn]bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.clients
}

// Pause NEEDSCOMMENT
func (c *clientMap) Pause(client *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.clients[client] = false
}

// Start NEEDSCOMMENT
func (c *clientMap) Start(client *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, alreadyHasClient := c.clients[client]
	c.clients[client] = true

	if !alreadyHasClient {
		for _, frame := range c.initialFrames {
			client.WriteMessage(websocket.BinaryMessage, frame)
		}
	}
}

func handleVideoRequests(ctx context.Context) handlerFunc {
	framesChan, initialFrames := startCamera(ctx)
	clients := newClientMap(initialFrames)

	upgrader := websocket.Upgrader{WriteBufferSize: 16384}

	go sendFramesToClients(framesChan, clients)

	return func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)
		defer ws.Close()

		if err != nil {
			log.Println("upgrade error", err)
			return
		}

		wsCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ws.SetCloseHandler(func(code int, text string) error {
			cancel()
			return nil
		})

		go handleVideoWebsocket(wsCtx, ws, clients, initialFrames)

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}
}

func startCamera(ctx context.Context) (chan []byte, [][]byte) {
	frameChan := make(chan []byte)
	initialFrames := [][]byte{}

	cmd := exec.CommandContext(ctx, "raspivid", raspividArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panicln("stdout error", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Panicln("stderr error", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer([]byte{}, 1024*1024)
	scanner.Split(splitAtNALSeparator)

	go func() {
		for scanner.Scan() {
			select {
			case frameChan <- scanner.Bytes():
			case <-ctx.Done():
				return
			}
		}
	}()

	var initialFramesWait sync.WaitGroup
	initialFramesWait.Add(1)
	go func() {
		defer initialFramesWait.Done()
		for frame := range frameChan {
			if len(initialFrames) >= initialFrameCount {
				return
			}
			initialFrames = append(initialFrames, frame)
		}
	}()

	go func() {
		if err := cmd.Start(); err != nil {
			log.Panicln("command start error", err)
		}

		if err := cmd.Wait(); err != nil {
			stderrLog, _ := ioutil.ReadAll(stderr)
			log.Panicln("command wait error", err, stderrLog)
		}
	}()

	initialFramesWait.Wait()

	return frameChan, initialFrames
}

func splitAtNALSeparator(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if idx := bytes.Index(data, nalSeparator); idx >= 0 {
		if idx == 0 {
			return idx + 4, nil, nil
		}
		return idx + 1, append(nalSeparator, data[:idx]...), nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func sendFramesToClients(frames chan []byte, clients *clientMap) {
	for frame := range frames {
		for client, isPlaying := range clients.Clients() {
			if isPlaying {
				client.WriteMessage(websocket.BinaryMessage, frame)
			}
		}
	}
}

func handleVideoWebsocket(ctx context.Context, ws *websocket.Conn, clients *clientMap, initialFrames [][]byte) {
	ws.WriteJSON(map[string]interface{}{
		"action": "init",
		"width":  width,
		"height": height,
	})

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return
		}

		if strings.HasPrefix(string(message), "REQUESTSTREAM") {
			clients.Start(ws)
		}

		if strings.HasPrefix(string(message), "STOPSTREAM") {
			clients.Pause(ws)
		}
	}
}
