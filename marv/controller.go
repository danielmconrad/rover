package marv

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ControllerState NEEDSCOMMENT
type ControllerState struct {
	Left  float64 `json:"left"`
	Right float64 `json:"right"`
}

// Message NEEDSCOMMENT
type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type messageHandler func(*Message) *Message

func handleControllerRequest(ctx context.Context, controllerChan chan *ControllerState) handlerFunc {
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

func handleMessageSocket(ctx context.Context, onMessage messageHandler) handlerFunc {
	upgrader := websocket.Upgrader{}

	return func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		defer ws.Close()

		for {
			messageType, messageBytes, err := ws.ReadMessage()
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

			err = ws.WriteMessage(messageType, responseBytes)
			if err != nil {
				log.Println("write error:", err)
			}
		}
	}
}
