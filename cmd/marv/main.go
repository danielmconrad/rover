package main

import (
	"log"
	"os"
	"time"

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

func main() {
	if err := rpio.Open(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	log.Println("Initializing Left")
	leftSpeedPin.Pwm()
	leftSpeedPin.Freq(1000 * int(maxSpeed))
	leftBackwardPin.Output()
	leftForwardPin.Output()

	log.Println("Initializing Right")
	rightSpeedPin.Pwm()
	rightSpeedPin.Freq(1000 * int(maxSpeed))
	rightBackwardPin.Output()
	rightForwardPin.Output()

	log.Println("Setting all to high for 2 seconds")
	leftSpeedPin.DutyCycle(90, maxSpeed)
	leftBackwardPin.Low()
	leftForwardPin.High()
	rightSpeedPin.DutyCycle(90, maxSpeed)
	rightBackwardPin.Low()
	rightForwardPin.High()
	time.Sleep(2 * time.Second)

	log.Println("Stopping")
	leftSpeedPin.DutyCycle(maxSpeed, maxSpeed)
	leftBackwardPin.Low()
	leftForwardPin.Low()
	rightSpeedPin.DutyCycle(maxSpeed, maxSpeed)
	rightBackwardPin.Low()
	rightForwardPin.Low()
	time.Sleep(1 * time.Second)

	log.Println("Cleaning up")
	leftSpeedPin.DutyCycle(0, maxSpeed)
	rightSpeedPin.DutyCycle(0, maxSpeed)
}
