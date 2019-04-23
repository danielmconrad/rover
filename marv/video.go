package marv

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	raspividArgs = []string{
		"-t", "0", "-o", "-", "-w", "320", "-h", "240", "-fps", "12", "-pf", "baseline",
	}
	nalSeparator = []byte{0x00, 0x00, 0x00, 0x01}
	width        = 320
	height       = 240
)

func handleVideo(ctx context.Context) handlerFunc {
	upgrader := websocket.Upgrader{}

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

		go handleWebsocket(wsCtx, ws)

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}
}

func handleWebsocket(ctx context.Context, ws *websocket.Conn) {
	pauseChan := make(chan bool)
	messageChan := make(chan []byte)

	handleStartStream := func() {
		for {
			select {
			case frame := <-startVideo(ctx):
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
			go handleStartStream()
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

		if err := cmd.Start(); err != nil {
			log.Println("command start error", err)
		}

		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitAtNALSeparator)

		for scanner.Scan() {
			select {
			case frameChan <- scanner.Bytes():
			case <-ctx.Done():
				return
			}
		}

		if err := cmd.Wait(); err != nil {
			log.Println("command wait error", err)
			return
		}
	}()

	return frameChan
}

func splitAtNALSeparator(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if idx := bytes.Index(data[4:], nalSeparator); idx >= 0 {
		return idx + 4, data[0 : idx+4], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
