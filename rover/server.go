package rover

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type handlerFunc func(w http.ResponseWriter, req *http.Request)

// StartServer NEEDSCOMMENT
func StartServer(ctx context.Context, port int, initialFrames []Frame) (<-chan *ControllerState, chan Frame) {
	controllerChan := make(chan *ControllerState)
	framesChan := make(chan Frame)
	mux := http.NewServeMux()

	mux.HandleFunc("/controller", handleControllerRequests(ctx, controllerChan))
	mux.HandleFunc("/video", handleVideoRequests(ctx, framesChan, initialFrames))
	mux.Handle("/", handleStaticRequests(ctx))

	go func() {
		defer close(controllerChan)
		logSuccess("Listening on port", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	}()

	return controllerChan, framesChan
}

func handleStaticRequests(ctx context.Context) http.Handler {
	return http.FileServer(http.Dir("static/"))
}

func handleVideoRequests(ctx context.Context, framesChan chan Frame, initialFrames []Frame) handlerFunc {
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

func handleVideoWebsocket(ctx context.Context, ws *websocket.Conn, clients *ClientMap, initialFrames []Frame) {
	ws.WriteJSON(map[string]interface{}{
		"action": "init",
		// "width":  width,
		// "height": height,
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

func sendFramesToClients(frames chan Frame, clients *ClientMap) {
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
