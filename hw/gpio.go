/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This package file implements the abstract hardware interfaces for GPIO
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
