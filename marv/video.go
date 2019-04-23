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

	"github.com/gorilla/websocket"
)

var (
	width        = 640
	height       = 360
	nalSeparator = []byte{0x00, 0x00, 0x00, 0x01}
	raspividArgs = []string{
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),

		"-rot", "180",
		"-fps", "24",
		"-t", "0",
		"-pf", "baseline",
		"-o", "-",
		// "-n", // No Preview
		// "-ev", "10", // Stabilization
		// "-g", "48", // Keyframes
		// "-ex", "antishake", // Antishake
	}
)

func handleVideoRequest(ctx context.Context) handlerFunc {
	upgrader := websocket.Upgrader{
		WriteBufferSize: 16384,
	}

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

		go handleVideoWebsocket(wsCtx, ws)

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}
}

func handleVideoWebsocket(ctx context.Context, ws *websocket.Conn) {
	pauseChan := make(chan bool)
	messageChan := make(chan []byte)

	handleStartStream := func(frameChan chan []byte) {
		for {
			select {
			case frame := <-frameChan:
				ws.WriteMessage(websocket.BinaryMessage, frame)
			case <-pauseChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}

	handleMessages := func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				return
			}
			messageChan <- message
		}
	}

	handleMessage := func(message []byte) {
		if strings.HasPrefix(string(message), "REQUESTSTREAM") {
			go handleStartStream(startVideo(ctx))
		}

		if strings.HasPrefix(string(message), "STOPSTREAM") {
			pauseChan <- true
		}
	}

	ws.WriteJSON(map[string]interface{}{
		"action": "init",
		"width":  width,
		"height": height,
	})

	go handleMessages()

	for {
		select {
		case message := <-messageChan:
			handleMessage(message)
		case <-ctx.Done():
			break
		}
	}
}

func startVideo(ctx context.Context) chan []byte {
	frameChan := make(chan []byte)

	go func() {
		cmd := exec.CommandContext(ctx, "raspivid", raspividArgs...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("stdout error", err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Println("stderr error", err)
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

		if err := cmd.Start(); err != nil {
			log.Println("command start error", err)
		}

		if err := cmd.Wait(); err != nil {
			stderrLog, _ := ioutil.ReadAll(stderr)
			log.Println("command wait error", err, stderrLog)
			return
		}
	}()

	return frameChan
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
