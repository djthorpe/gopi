/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// GPIO
//
// The abstract GPIO hardware interface can be used for interfacing a
// variety of external devices which use simple digital inputs and outputs.
// In order to use, construct a GPIO driver object. For the Raspberry Pi,
// you can acheive this using a rpi.GPIO object. For example,
//
//   gpio := gopi.Open(rpi.GPIO{ ... })
//   defer gpio.Close()
//
// When you have finished using the driver, use the Close method which will
// free up any resources. The pins on your GPIO connector have a physical
// pin value and a logical pin name. In order to convert from the physical
// pin number and vice-versa, use the following methods:
//
//   logical := gpio.PhysicalPin(40)
//   physical := gpio.PhysicalPinForPin(logical) // should be 40
//
// This will return GPIO_PIN_NONE when no logical pin is available at this
// physical pin position. You can also get a list of all logical pins and
// the number of physical pins on your GPIO connector using the following
// methods:
//
//  pins := gpio.Pins() // returns an array of logical pins
//  number_of_physical_pins := gpio.NumberOfPhysicalPins()
//
// You can read or write a pin to LOW or HIGH state using the following
// methods:
//
//  state := GPIO_HIGH
//  gpio.WritePin(logical,state)
//  state = gpio.ReadPin(logical) // Should be GPIO_HIGH
//
// And set a pin to INPUT or OUTPUT, and set resistor pull-up or pull-down:
//
//  gpio.SetPinMode(logical,GPIO_INPUT)
//  gpio.SetPullMode(logical,GPIO_PULL_UP)
//
// On the Raspberry Pi, you can also set pins to "alternate" modes. For example,
//
//  gpio.SetPinMode(logical,GPIO_ALT0)
//
package hw // import "github.com/djthorpe/gopi/hw"

import (
	"fmt"
)

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract GPIO interface
type GPIODriver interface {
	// Enforces general driver
	gopi.Driver

	// Return number of physical pins
	NumberOfPhysicalPins() uint

	// Return array of available logical pins
	Pins() []GPIOPin

	// Return logical pin for physical pin number. Returns
	// GPIO_PIN_NONE where there is no logical pin at that position
	PhysicalPin(uint) GPIOPin

	// Return physical pin number for logical pin. Returns 0 where there
	// is no physical pin for this logical pin
	PhysicalPinForPin(GPIOPin) uint

	// Read pin state
	ReadPin(GPIOPin) GPIOState

	// Write pin state
	WritePin(GPIOPin, GPIOState)

	// Get pin mode
	GetPinMode(GPIOPin) GPIOMode

	// Set pin mode
	SetPinMode(GPIOPin, GPIOMode)

	// Set pull mode
	SetPullMode(GPIOPin, GPIOPull)
}

// GPIO types
type (
	// Logical GPIO pin
	GPIOPin uint8

	// GPIO Pin state
	GPIOState uint8

	// GPIO Pin mode
	GPIOMode uint8

	// GPIO Pin resistor configuration (pull up/down or floating)
	GPIOPull uint8
)

////////////////////////////////////////////////////////////////////////////////
// GPIO CONSTANTS

const (
	// Invalid pin constant
	GPIO_PIN_NONE GPIOPin = 0xFF
)

const (
	GPIO_LOW GPIOState = iota
	GPIO_HIGH
)

const (
	GPIO_INPUT GPIOMode = iota
	GPIO_OUTPUT
	GPIO_ALT5
	GPIO_ALT4
	GPIO_ALT0
	GPIO_ALT1
	GPIO_ALT2
	GPIO_ALT3
)

const (
	GPIO_PULL_OFF GPIOPull = iota
	GPIO_PULL_DOWN
	GPIO_PULL_UP
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (p GPIOPin) String() string {
	return fmt.Sprintf("GPIO%v", uint8(p))
}

func (s GPIOState) String() string {
	switch s {
	case GPIO_LOW:
		return "LOW"
	case GPIO_HIGH:
		return "HIGH"
	default:
		return "[??? Invalid GPIOState value]"
	}
}

func (m GPIOMode) String() string {
	switch m {
	case GPIO_INPUT:
		return "INPUT"
	case GPIO_OUTPUT:
		return "OUTPUT"
	case GPIO_ALT0:
		return "ALT0"
	case GPIO_ALT1:
		return "ALT1"
	case GPIO_ALT2:
		return "ALT2"
	case GPIO_ALT3:
		return "ALT3"
	case GPIO_ALT4:
		return "ALT4"
	case GPIO_ALT5:
		return "ALT5"
	default:
		return "[??? Invalid GPIOMode value]"
	}
}

func (p GPIOPull) String() string {
	switch p {
	case GPIO_PULL_OFF:
		return "PULL_OFF"
	case GPIO_PULL_DOWN:
		return "PULL_DOWN"
	case GPIO_PULL_UP:
		return "PULL_UP"
	default:
		return "[??? Invalid GPIOPull value]"
	}
}
