package marv

import (
	"context"
	"log"
	"math"
	"os"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	leftSpeedPin    = rpio.Pin(12)
	leftForwardPin  = rpio.Pin(16)
	leftBackwardPin = rpio.Pin(20)

	rightSpeedPin    = rpio.Pin(13)
	rightForwardPin  = rpio.Pin(5)
	rightBackwardPin = rpio.Pin(6)

	maxSpeed = uint32(100)
)

// MotorState NEEDSCOMMENT
type MotorState struct {
	Left  int32
	Right int32
}

// StartMotors NEEDSCOMMENT
func StartMotors(ctx context.Context) chan *MotorState {
	motorChan := make(chan *MotorState)

	lastState := &MotorState{}

	go func() {
		if err := rpio.Open(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		defer rpio.Close()

		initializeMotor(leftSpeedPin, leftBackwardPin, leftForwardPin)
		initializeMotor(rightSpeedPin, rightBackwardPin, rightForwardPin)

		for {
			select {
			case motorState := <-motorChan:
				// log.Println(motorState)

				if motorState.Left != lastState.Left || lastState.Right != lastState.Right {
					setMotors(motorState)
					lastState = motorState
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return motorChan
}

func initializeMotor(speedPin, backwardPin, forwardPin rpio.Pin) {
	speedPin.Pwm()
	speedPin.Freq(1000 * int(maxSpeed))
	backwardPin.Output()
	forwardPin.Output()
}

func setMotors(motorState *MotorState) {
	setMotor(leftSpeedPin, leftBackwardPin, leftForwardPin, motorState.Left)
	setMotor(rightSpeedPin, rightBackwardPin, rightForwardPin, motorState.Right)
}

func setMotor(speedPin, backwardPin, forwardPin rpio.Pin, speed int32) {
	absSpeed := uint32(math.Abs(float64(speed)))

	if speed > 20 {
		backwardPin.Low()
		forwardPin.High()
		speedPin.DutyCycle(absSpeed, maxSpeed)
	} else if speed < -20 {
		forwardPin.Low()
		backwardPin.High()
		speedPin.DutyCycle(absSpeed, maxSpeed)
	} else {
		forwardPin.High()
		backwardPin.High()
		speedPin.DutyCycle(maxSpeed, maxSpeed) // Stop
	}
}
