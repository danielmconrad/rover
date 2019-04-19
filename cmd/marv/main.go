package main

import (
	"context"

	"github.com/danielmconrad/gomarv/marv"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// create video instance and send into the server

	controllerChan := marv.StartServer(ctx, 3737)
	motorChan := marv.StartMotors(ctx)

	for {
		select {
		case controllerState := <-controllerChan:
			// log.Println("controllerState", controllerState)

			motorState := &marv.MotorState{
				Left:  int32(controllerState.Left * 100),
				Right: int32(controllerState.Right * 100),
			}

			// log.Println("motorState", motorState)

			motorChan <- motorState
		case <-ctx.Done():
			return
		}
	}
}
