package rover

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	width  = 640
	height = 400
	fps    = 48

	ffmpegArgsLinux = []string{}

	raspividArgs = []string{
		"raspivid",
		"-w", strconv.Itoa(width), "-h", strconv.Itoa(height), "-fps", strconv.Itoa(fps),
		"-t", "0", "-pf", "baseline", "-o", "-",
	}

	ffmpegArgsMac = []string{
		"ffmpeg",

		// Input
		"-f", "avfoundation", "-framerate", "30", "-pixel_format", "yuyv422",
		"-video_size", "640x480", "-i", "0",

		// Output
		"-vcodec", "libx264", "-profile:v", "baseline", "-pix_fmt", "yuv420p", "-level:v", "4.2",
		"-preset", "ultrafast", "-tune", "zerolatency", "-bufsize", "0", "-crf", "22",
		"-f", "rawvideo", "-",
	}

	initialFrameCount = 4
	nalSeparator      = []byte{0x00, 0x00, 0x00, 0x01}

	isRaspberry = runtime.GOOS == "linux" && runtime.GOARCH == "arm"
	isMac       = runtime.GOOS == "darwin"
)

func handleVideoRequests(ctx context.Context) handlerFunc {
	framesChan, initialFrames := startCamera(ctx)
	clients := NewClientMap(initialFrames)

	upgrader := websocket.Upgrader{WriteBufferSize: 16384}

	go sendFramesToClients(framesChan, clients)

	return func(w http.ResponseWriter, req *http.Request) {
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

		go handleVideoWebsocket(wsCtx, ws, clients, initialFrames)

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}
}

func handleVideoWebsocket(ctx context.Context, ws *websocket.Conn, clients *ClientMap, initialFrames [][]byte) {
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

		if strings.Contains(string(message), "start") {
			clients.Start(ws)
		}

		if strings.Contains(string(message), "pause") {
			clients.Pause(ws)
		}
	}
}

func sendFramesToClients(frames chan []byte, clients *ClientMap) {
	for frame := range frames {
		for client, isPlaying := range clients.Clients() {
			if isPlaying {
				client.WriteMessage(websocket.BinaryMessage, frame)
			}
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

func startCamera(ctx context.Context) (chan []byte, [][]byte) {
	frameChan := make(chan []byte)
	initialFrames := [][]byte{}

	args := ffmpegArgsLinux

	if isMac {
		args = ffmpegArgsMac
	}

	if isRaspberry {
		args = raspividArgs
	}

	logInfo("Starting video: ", args)

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logPanic("stdout error", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logPanic("stderr error", err)
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
			logPanic("command start error", err)
		}

		if err := cmd.Wait(); err != nil {
			stderrLog, _ := ioutil.ReadAll(stderr)
			logPanic("command wait error", err, stderrLog)
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
