package main

import (
	"context"
	"log"

	"github.com/danielmconrad/gomarv/marv/motors"
	"github.com/danielmconrad/gomarv/marv/server"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	messages := server.Listen(ctx, 3737)
	commands := motors.Start(ctx)

	for {
		select {
		case message := <-messages:
			commands <- messageToCommand(message)
		case <-ctx.Done():
			return
		}
	}
}

func messageToCommand(message *server.Message) *motors.Command {
	log.Println("Received message", message)
	log.Println("Will send command")
	return &motors.Command{}
}
