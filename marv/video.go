package marv

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	ffmpegArgs = []string{
		"-f", "v4l2", "-framerate", "25", "-video_size", "320x240", "-i", "/dev/video0",
		"-f", "mpegts", "-codec:v", "mpeg1video", "-s", "320x240", "-b:v", "1000k", "-bf", "0", "pipe:1",
	}
	raspividArgs = []string{
		"-t", "0", "-o", "-", "-w", "960", "-h", "540", "-fps", "12", "-pf", "baseline",
	}
	nalSeparator = []byte{0x00, 0x00, 0x00, 0x01}
)

func handleVideo(ctx context.Context) handlerFunc {
	upgrader := websocket.Upgrader{}

	clients := map[*websocket.Conn]bool{}

	go func() {
		for videoFragment := range startVideo() {
			for client := range clients {
				client.WriteMessage(websocket.BinaryMessage, videoFragment)
			}
		}
	}()

	return func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)

		if err != nil {
			log.Println("upgrade error", err)
			return
		}

		ws.WriteJSON(map[string]interface{}{
			"action": "init",
			"width":  960,
			"height": 540,
		})

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				ws.Close()
				return
			}

			if strings.HasPrefix(string(message), "REQUESTSTREAM") {
				clients[ws] = true

				ws.SetCloseHandler(func(code int, text string) error {
					delete(clients, ws)
					return nil
				})
			}
		}
	}
}

func startVideo() chan []byte {
	dataChan := make(chan []byte)

	go func() {
		// cmd := exec.Command("ffmpeg", ffmpegArgs...)
		cmd := exec.Command("raspivid", raspividArgs...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("stdout error", err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Println("stderr error", err)
		}

		if err := cmd.Start(); err != nil {
			log.Println("command start error", err)
		}

		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))

		prevBytes := []byte{}

		for scanner.Scan() {
			allBytes := append(prevBytes, scanner.Bytes()...)

			if idx := bytes.Index(allBytes[1:], nalSeparator); idx >= 0 {
				dataChan <- allBytes[:idx+1]
				prevBytes = allBytes[idx+1:]
			}
		}

		if err := cmd.Wait(); err != nil {
			log.Println("command wait error", err)
			return
		}
	}()

	return dataChan
}
