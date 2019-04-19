package marv

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

// Message NEEDSCOMMENT
type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// ControllerState NEEDSCOMMENT
type ControllerState struct {
	Left  float64 `json:"left"`
	Right float64 `json:"right"`
}

type messageHandler func(*Message) *Message

type handlerFunc func(w http.ResponseWriter, req *http.Request)

// StartServer NEEDSCOMMENT
func StartServer(ctx context.Context, port int) <-chan *ControllerState {
	controllerChan := make(chan *ControllerState)
	mux := http.NewServeMux()

	mux.HandleFunc("/controller", handleController(ctx, controllerChan))
	mux.HandleFunc("/video", handleVideo(ctx))
	mux.Handle("/", handleStatic(ctx))

	go func() {
		defer close(controllerChan)
		log.Printf("Listening on port %d", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	}()

	return controllerChan
}

func handleVideo(ctx context.Context) handlerFunc {
	upgrader := websocket.Upgrader{}

	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		defer conn.Close()

		// cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-ab", "128k", "-f", "mp3", "pipe:1")
		cmd := exec.Command(
			"ffmpeg",
			"-f", "v4l2",
			"-framerate", "25",
			"-video_size", "640x480",
			"-i", "/dev/video0",
			"-f", "mpegts",
			"-codec:v", "mpeg1video",
			"-s", "640x480",
			"-b:v", "1000k",
			"-bf", "0",
			"pipe:1")

		stdout, err := cmd.StdoutPipe()

		if err != nil {
			log.Panic(err)
		}
		err = cmd.Start()
		io.Copy(w, stdout)

		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		for {
			select {
			case err := <-done:
				log.Println("command error", err)
				return
			case <-ctx.Done():
				cmd.Process.Kill()
				return
			}
		}
	}
}

func handleController(ctx context.Context, controllerChan chan *ControllerState) handlerFunc {
	return handleMessageSocket(ctx, func(message *Message) *Message {
		if len(message.Data) == 0 {
			return &Message{Event: "nodata"}
		}

		controllerState := &ControllerState{}

		err := json.Unmarshal([]byte(message.Data), controllerState)
		if err != nil {
			log.Println("unmarshal error:", err)
			return &Message{Event: "error", Data: err.Error()}
		}

		controllerChan <- controllerState
		return &Message{Event: "ack"}
	})
}

func handleStatic(ctx context.Context) http.Handler {
	return http.FileServer(http.Dir("marv/static/"))
}

func handleMessageSocket(ctx context.Context, onMessage messageHandler) handlerFunc {
	upgrader := websocket.Upgrader{}

	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		defer conn.Close()

		for {
			messageType, messageBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				break
			}

			message := &Message{}
			err = json.Unmarshal(messageBytes, message)
			if err != nil {
				log.Println("unmarshal error:", err)
			}

			response := onMessage(message)

			responseBytes, err := json.Marshal(response)
			if err != nil {
				log.Println("marshal error:", err)
			}

			err = conn.WriteMessage(messageType, responseBytes)
			if err != nil {
				log.Println("write error:", err)
			}
		}
	}
}
