package main

import (
	"context"
	"log"

	"github.com/danielmconrad/gomarv/marv"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	controllerChan := marv.StartServer(ctx, 3737)
	motorChan := marv.StartMotors(ctx)

	for {
		select {
		case controllerState := <-controllerChan:
			log.Println("controllerState", controllerState)
			motorChan <- &marv.MotorState{Left: 100, Right: 100}
		case <-ctx.Done():
			return
		}
	}
}
