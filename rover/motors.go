package rover

import (
	"context"
	"math"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	leftSpeedPin    = rpio.Pin(12)
	leftForwardPin  = rpio.Pin(20)
	leftBackwardPin = rpio.Pin(16)

	rightSpeedPin    = rpio.Pin(13)
	rightForwardPin  = rpio.Pin(6)
	rightBackwardPin = rpio.Pin(5)

	maxSpeed         = uint32(100)
	minimumChange    = float64(3)
	sendMotorSignals = true
)

// MotorState NEEDSCOMMENT
type MotorState struct {
	Left  int32
	Right int32
}

// StartMotors NEEDSCOMMENT
func StartMotors(ctx context.Context) chan *MotorState {
	motorChan := make(chan *MotorState)
	prevState := &MotorState{}

	go func() {
		if err := rpio.Open(); err != nil {
			logWarning("Disabling motor control")
			sendMotorSignals = false
		}
		defer rpio.Close()

		initializeMotors()

		for {
			select {
			case nextState := <-motorChan:
				if differentEnough(prevState, nextState) {
					setVelocities(nextState)
					prevState = nextState
				}
			case <-ctx.Done():
				fullStopMotors()
				return
			}
		}
	}()

	return motorChan
}

func differentEnough(prevState, nextState *MotorState) bool {
	leftDiff := math.Abs(float64(nextState.Left - prevState.Left))
	rightDiff := math.Abs(float64(nextState.Right - prevState.Right))

	return leftDiff > minimumChange || rightDiff > minimumChange
}

func initializeMotors() {
	initializeMotor(leftSpeedPin, leftBackwardPin, leftForwardPin)
	initializeMotor(rightSpeedPin, rightBackwardPin, rightForwardPin)
}

func initializeMotor(speedPin, backwardPin, forwardPin rpio.Pin) {
	if !sendMotorSignals {
		return
	}

	speedPin.Pwm()
	speedPin.Freq(1000 * int(maxSpeed))
	backwardPin.Output()
	forwardPin.Output()
}

func setVelocities(motorState *MotorState) {
	setVelocity(leftSpeedPin, leftBackwardPin, leftForwardPin, motorState.Left)
	setVelocity(rightSpeedPin, rightBackwardPin, rightForwardPin, motorState.Right)
}

func setVelocity(speedPin, backwardPin, forwardPin rpio.Pin, speed int32) {

	if !sendMotorSignals {
		return
	}

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
		speedPin.DutyCycle(maxSpeed, maxSpeed)
	}
}

func fullStopMotors() {
	fullStopMotor(leftSpeedPin, leftBackwardPin, leftForwardPin)
	fullStopMotor(rightSpeedPin, rightBackwardPin, rightForwardPin)
}

func fullStopMotor(speedPin, backwardPin, forwardPin rpio.Pin) {
	if !sendMotorSignals {
		return
	}

	forwardPin.Low()
	backwardPin.Low()
	speedPin.DutyCycle(0, maxSpeed)
}
