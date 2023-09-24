package main

import (
	"machine"
	"time"
)

var switchState = true

func main() {
	encoder := NewRotaryEncoder[int32](machine.A2, machine.A1, 4095, 10, 0, 4095)
	encoder.Start()

	i2c := machine.I2C1
	err := i2c.Configure(machine.I2CConfig{
		SCL: machine.D5,
		SDA: machine.D4,
	})

	encoderLED := machine.D7
	encoderLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	encoderLED.Low()

	encoderSwitch := machine.D8
	encoderSwitch.Configure(machine.PinConfig{Mode: machine.PinInput})

	encoderSwitch.SetInterrupt(machine.PinRising, func(machine.Pin) {
		switchState = !switchState
		if switchState {
			encoderLED.Low()
		} else {
			encoderLED.High()
		}
	})

	if err != nil {
		println(err.Error())
		return
	}

	println("Start")

	for {
		println(encoder.Value)
		value := encoder.Value
		if !switchState {
			value = 0
		}
		// I2C で 12bit の値を DAC に送る
		// 0x60: MCP4726 の I2C アドレス
		wErr := i2c.WriteRegister(0x60, 0, []byte{0x00, byte((value >> 8) & 0x0F), byte(value & 0xFF)})
		if wErr != nil {
			println(wErr.Error())
		}
		time.Sleep(time.Millisecond * 100)
	}
}
