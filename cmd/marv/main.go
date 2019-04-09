package main

import (
	"context"
	"log"

	"github.com/danielmconrad/gomarv/marv"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	messages := marv.StartServer(ctx, 3737)
	commands := marv.StartMotors(ctx)

	for {
		select {
		case message := <-messages:
			commands <- messageToCommand(message)
		case <-ctx.Done():
			return
		}
	}
}

func messageToCommand(message *marv.Message) *marv.Command {
	log.Println("Received message", message)
	return &marv.Command{}
}
