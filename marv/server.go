package marv

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type handlerFunc func(w http.ResponseWriter, req *http.Request)

// StartServer NEEDSCOMMENT
func StartServer(ctx context.Context, port int) <-chan *ControllerState {
	controllerChan := make(chan *ControllerState)
	mux := http.NewServeMux()

	mux.HandleFunc("/controller", handleControllerRequest(ctx, controllerChan))
	mux.HandleFunc("/video", handleVideoRequest(ctx))
	mux.Handle("/", handleStaticRequest(ctx))

	go func() {
		defer close(controllerChan)
		log.Printf("Listening on port %d", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	}()

	return controllerChan
}

func handleStaticRequest(ctx context.Context) http.Handler {
	return http.FileServer(http.Dir("static/"))
}
