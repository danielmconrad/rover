package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message NEEDSCOMMENT
type Message struct {
	Event string `json:"event"`
}

type handlerFunc func(w http.ResponseWriter, req *http.Request)
type messageHandler func(*Message) *Message

// Listen NEEDSCOMMENT
func Listen(ctx context.Context, port int, onMessage messageHandler) {
	mux := http.NewServeMux()

	mux.HandleFunc("/messages", handleMessages(ctx, onMessage))
	mux.Handle("/", handleStatic(ctx))

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func handleStatic(ctx context.Context) http.Handler {
	return http.FileServer(http.Dir("marv/static/"))
}

func handleMessages(ctx context.Context, onMessage messageHandler) handlerFunc {
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
				break
			}
		}
	}
}
