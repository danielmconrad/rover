package main

import (
	"context"
	"log"

	"github.com/danielmconrad/gomarv/marv/motors"
	"github.com/danielmconrad/gomarv/marv/server"
)

func main() {
	messages := make(chan *server.Message)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go runServer(ctx, messages)

	for {
		select {
		case message := <-messages:
			consumeMessage(message)
		case <-ctx.Done():
			return
		}
	}
}

func runServer(ctx context.Context, messages chan<- *server.Message) {
	server.Listen(ctx, 3737, func(message *server.Message) *server.Message {
		messages <- message
		return &server.Message{Event: "ack"}
	})
}

func consumeMessage(message *server.Message) {
	log.Println("Received message", message)
	command := motors.JoystickCommand{}
	log.Println("Will send command", command)
}
