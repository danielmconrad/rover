package marv

import (
	"context"
	"fmt"
	"net/http"
)

type handlerFunc func(w http.ResponseWriter, req *http.Request)

// StartServer NEEDSCOMMENT
func StartServer(ctx context.Context, port int) <-chan *ControllerState {
	controllerChan := make(chan *ControllerState)
	mux := http.NewServeMux()

	mux.HandleFunc("/controller", handleControllerRequests(ctx, controllerChan))
	mux.HandleFunc("/video", handleVideoRequests(ctx))
	mux.Handle("/", handleStaticRequests(ctx))

	go func() {
		defer close(controllerChan)
		logSuccess("Listening on port", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	}()

	return controllerChan
}

func handleStaticRequests(ctx context.Context) http.Handler {
	return http.FileServer(http.Dir("static/"))
}
