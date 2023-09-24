package main

import (
	"machine"
)

var lookupTable = []int8{
	0, 1, -1, 0, // 00 -> 00, 01 -> 01, 10 -> 11, 11 -> 10
	-1, 0, 0, 1, // 00 -> 10, 01 -> 00, 10 -> 01, 11 -> 11
	1, 0, 0, -1, // 00 -> 11, 01 -> 10, 10 -> 00, 11 -> 01
	0, -1, 1, 0, // 00 -> 01, 01 -> 11, 10 -> 10, 11 -> 00
}

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

/**
 * RotaryEncoder struct
 */
type RotaryEncoder[T Integer] struct {
	PinA   machine.Pin
	PinB   machine.Pin
	Value  T
	step   T
	min    T
	max    T
	prevAB T
}

/**
 * Create a new rotary encoder
 * step: the step value to increment/decrement the value
 */
func NewRotaryEncoder[T Integer](pinA machine.Pin, pinB machine.Pin, initValue T, step T, min T, max T) *RotaryEncoder[T] {
	return &RotaryEncoder[T]{
		PinA:   pinA,
		PinB:   pinB,
		Value:  initValue,
		step:   step,
		min:    min,
		max:    max,
		prevAB: 0x03,
	}
}

/**
 * Start the rotary encoder
 */
func (encoder *RotaryEncoder[T]) Start() {
	encoder.PinA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	encoder.PinB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	encoder.PinA.SetInterrupt(machine.PinRising, encoder.interrupt)
	encoder.PinB.SetInterrupt(machine.PinRising, encoder.interrupt)
}

/**
 * PIN Interrupt handler for the rotary encoder
 */
func (encoder *RotaryEncoder[T]) interrupt(pin machine.Pin) {
	isPinAHigh := encoder.PinA.Get()
	isPinBHigh := encoder.PinB.Get()

	encoder.prevAB = (encoder.prevAB << 2)
	if isPinAHigh {
		encoder.prevAB |= 0x02
	}
	if isPinBHigh {
		encoder.prevAB |= 0x01
	}
	addition := T(lookupTable[encoder.prevAB&0x0f]) * encoder.step
	// overflow check
	next, overflow := add[T](encoder.Value, addition)
	if overflow {
		// set to max value or min value
		if next < 0 {
			next = encoder.max
		} else {
			// min value of T
			next = encoder.min
		}
	}

	if next > encoder.max {
		next = encoder.max
	} else if next < encoder.min {
		next = encoder.min
	}

	encoder.Value = next
}

func add[T Integer](a T, b T) (T, bool) {
	sum := a + b
	overflow := (a > 0 && b > 0 && sum < 0) || (a < 0 && b < 0 && sum > 0)
	return sum, overflow
}
