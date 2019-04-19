package marv

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type handlerFunc func(w http.ResponseWriter, req *http.Request)

// StartServer NEEDSCOMMENT
func StartServer(ctx context.Context, port int) <-chan *Message {
	controllerChan := make(chan *Message)
	mux := http.NewServeMux()

	mux.HandleFunc("/controller", handleController(ctx, func(message *Message) *Message {
		controllerChan <- message
		return &Message{Event: "ack"}
	}))
	mux.Handle("/", handleStatic(ctx))

	go func() {
		defer close(controllerChan)
		log.Printf("Listening on port %d", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	}()

	return controllerChan
}

func handleStatic(ctx context.Context) http.Handler {
	return http.FileServer(http.Dir("marv/static/"))
}

func handleController(ctx context.Context, onMessage messageHandler) handlerFunc {
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
