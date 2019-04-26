package marv

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gorilla/websocket"
)

// ControllerState NEEDSCOMMENT
type ControllerState struct {
	Left  float64 `json:"left"`
	Right float64 `json:"right"`
}

func handleControllerRequests(ctx context.Context, controllerChan chan *ControllerState) handlerFunc {
	upgrader := websocket.Upgrader{}

	return func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			logError("Upgrade error", err)
			return
		}
		defer ws.Close()

		wsCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ws.SetCloseHandler(func(code int, text string) error {
			cancel()
			return nil
		})

		go handleControllerWebsocket(wsCtx, ws, controllerChan)

		for {
			select {
			case <-ctx.Done():
				return
			}
		}

	}
}

func handleControllerWebsocket(ctx context.Context, ws *websocket.Conn, controllerChan chan *ControllerState) {
	ws.WriteJSON(map[string]interface{}{
		"action": "init",
	})

	previousControllerState := &ControllerState{}

	for {
		controllerState := &ControllerState{}

		err := ws.ReadJSON(controllerState)
		if err != nil {
			return
		}

		if !reflect.DeepEqual(*controllerState, *previousControllerState) {
			controllerChan <- controllerState
		}

		previousControllerState = controllerState
	}
}
