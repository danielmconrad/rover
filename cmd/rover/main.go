package main

import (
	"context"
	"flag"

	"github.com/danielmconrad/rover/rover"
)

var (
	width  = uint64(640)
	height = uint64(400)
	fps    = uint64(48)
)

func main() {
	parseVariables()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	cameraFramesChan, initialFrames := rover.StartCamera(ctx, width, height, fps)
	controllerChan, serverFramesChan := rover.StartServer(ctx, 3737, initialFrames)
	motorChan := rover.StartMotors(ctx)

	for {
		select {
		case frame := <-cameraFramesChan:
			serverFramesChan <- frame
		case controllerState := <-controllerChan:
			motorChan <- &rover.MotorState{
				Left:  int32(controllerState.Left * 100),
				Right: int32(controllerState.Right * 100),
			}
		case <-ctx.Done():
			return
		}
	}
}

func parseVariables() {
	flag.Uint64Var(&width, "width", width, "Video width")
	flag.Uint64Var(&width, "w", width, "Video width (short)")

	flag.Uint64Var(&height, "height", height, "Video height")
	flag.Uint64Var(&height, "h", height, "Video height (short)")

	flag.Uint64Var(&fps, "fps", fps, "Video framerate")
	flag.Uint64Var(&fps, "f", fps, "Video framerate (short)")

	flag.Parse()
}
