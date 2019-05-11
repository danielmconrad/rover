package marv

import (
	"bufio"
	"context"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	width        = 640
	height       = 400
	raspividArgs = []string{
		"-w", strconv.Itoa(width), "-h", strconv.Itoa(height),
		"-fps", "48", "-t", "0", "-pf", "baseline", "-o", "-",
	}
	initialFrameBufferSize = 1024 * 24
	frameSize              = 4096
	scannerBufferSize      = 1024 * 1024
)

func handleVideoRequests(ctx context.Context) handlerFunc {
	if runtime.GOOS != "linux" || runtime.GOARCH != "arm" {
		logWarning("Video not supported")
		return handleUnsupportedVideoWebsocket
	}

	upgrader := websocket.Upgrader{WriteBufferSize: 16384}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ws, err := upgrader.Upgrade(w, req, nil)
		defer ws.Close()

		if err != nil {
			logError("upgrade error", err)
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
	// ws.WriteJSON(map[string]interface{}{
	// 	"action": "init",
	// 	"width":  width,
	// 	"height": height,
	// })

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return
		}

		logInfo("Incoming message", string(message))

		if parts := strings.Split(string(message), " "); len(parts) > 0 && parts[0] == "start" {
			go func() {
				for frame := range startCamera(ctx) {
					ws.WriteMessage(websocket.BinaryMessage, frame)
				}
			}()
		}
	}
}

func handleUnsupportedVideoWebsocket(w http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{WriteBufferSize: 16384}
	ws, err := upgrader.Upgrade(w, req, nil)
	defer ws.Close()

	if err != nil {
		logError("upgrade error", err)
		return
	}

	ws.WriteJSON(map[string]interface{}{
		"action": "error",
		"error":  "Not Supported",
	})
}

func startCamera(ctx context.Context) chan []byte {
	frameChan := make(chan []byte)

	cmd := exec.CommandContext(ctx, "raspivid", raspividArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logPanic("stdout error", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logPanic("stderr error", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer([]byte{}, scannerBufferSize)
	scanner.Split(splitAtLength(frameSize))

	go func() {
		for scanner.Scan() {
			select {
			case frameChan <- scanner.Bytes():
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		if err := cmd.Start(); err != nil {
			logPanic("command start error", err)
		}

		if err := cmd.Wait(); err != nil {
			stderrLog, _ := ioutil.ReadAll(stderr)
			logPanic("command wait error", err, stderrLog)
		}
	}()

	return frameChan
}

func splitAtLength(splitLength int) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if len(data) >= splitLength {
			dst := make([]byte, splitLength)
			copy(dst, data[:splitLength])
			return splitLength, dst, nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}
}
