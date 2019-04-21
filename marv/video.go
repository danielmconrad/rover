package marv

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

var (
	ffmpegArgs = []string{
		"-f", "v4l2", "-framerate", "25", "-video_size", "320x240", "-i", "/dev/video0",
		"-f", "mpegts", "-codec:v", "mpeg1video", "-s", "320x240", "-b:v", "1000k", "-bf", "0", "pipe:1",
	}
	raspividArgs = []string{
		"-v", "-w", "320", "-h", "240", "-fps", "12", "-n", "-md", "7", "-ih", "-t", "0", "-o", "-",
	}
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

		go func(c *websocket.Conn) {
			for {
				if _, _, err := c.NextReader(); err != nil {
					c.Close()
					break
				}
			}
		}(ws)

		clients[ws] = true
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

		s := bufio.NewScanner(io.MultiReader(stdout, stderr))

		oldBuff := []byte{}

		for s.Scan() {
			b := s.Bytes()

			if len(b) > 7 && b[4] == 'm' && b[5] == 'o' && b[6] == 'o' {
				dataChan <- oldBuff
				oldBuff = b
			} else {
				oldBuff = append(oldBuff, b...)
			}
		}

		if err := cmd.Wait(); err != nil {
			log.Println("command wait error", err)
			return
		}
	}()

	return dataChan
}
