package main

import (
	"time"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/megapi"
)

func main() {
	// use "/dev/ttyUSB0" if connecting with USB cable
	// use "/dev/ttyAMA0" on devices older than Raspberry Pi 3 Model B
	megaPiAdaptor := megapi.NewAdaptor("/dev/ttyS0")
	motor1 := megapi.NewMotorDriver(megaPiAdaptor, 1)
	motor2 := megapi.NewMotorDriver(megaPiAdaptor, 2)
	motor3 := megapi.NewMotorDriver(megaPiAdaptor, 3)
	motor4 := megapi.NewMotorDriver(megaPiAdaptor, 4)

	work := func() {
		speed := int16(0)
		fadeAmount := int16(30)

		gobot.Every(100*time.Millisecond, func() {
			motor1.Speed(speed)
			motor2.Speed(speed)
			motor3.Speed(speed)
			motor4.Speed(speed)

			if speed <= 300 {
				speed = speed + fadeAmount
			}

			// if speed == 0 || speed == 300 {
			// 	fadeAmount = -fadeAmount
			// }
		})
	}

	robot := gobot.NewRobot("megaPiBot",
		[]gobot.Connection{megaPiAdaptor},
		[]gobot.Device{motor1, motor2, motor3, motor4},
		work,
	)

	robot.Start()
}
