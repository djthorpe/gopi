/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Hardware implements the hardware driver interface, which
// provides information about the hardware that the software is
// running on
type Hardware interface {
	Driver

	// Return name of the hardware platform
	Name() string

	// Return unique serial number of this hardware
	SerialNumber() string

	// Return the number of possible displays for this hardware
	NumberOfDisplays() uint
}

// Display implements a pixel-based display device. Displays are always numbered
// from zero onwards
type Display interface {
	Driver

	// Return display number
	Display() uint

	// Return display size for nominated display number, or (0,0) if display
	// does not exist
	Size() (uint32, uint32)

	// Return the PPI (pixels-per-inch) for the display, or return zero if unknown
	PixelsPerInch() uint32
}

// GPIO implements the GPIO interface for simple input and output
type GPIO interface {
	// Enforces general driver and event publisher
	Driver
	Publisher

	// Return number of physical pins, or 0 if if cannot be returned
	// or nothing is known about physical pins
	NumberOfPhysicalPins() uint

	// Return array of available logical pins or nil if nothing is
	// known about pins
	Pins() []GPIOPin

	// Return logical pin for physical pin number. Returns
	// GPIO_PIN_NONE where there is no logical pin at that position
	// or we don't now about the physical pins
	PhysicalPin(uint) GPIOPin

	// Return physical pin number for logical pin. Returns 0 where there
	// is no physical pin for this logical pin, or we don't know anything
	// about the layout
	PhysicalPinForPin(GPIOPin) uint

	// Read pin state
	ReadPin(GPIOPin) GPIOState

	// Write pin state
	WritePin(GPIOPin, GPIOState)

	// Get pin mode
	GetPinMode(GPIOPin) GPIOMode

	// Set pin mode
	SetPinMode(GPIOPin, GPIOMode)

	// Set pull mode to pull down or pull up - will
	// return ErrNotImplemented if not supported
	SetPullMode(GPIOPin, GPIOPull) error

	// Start watching for rising and/or falling edge,
	// or stop watching when GPIO_EDGE_NONE is passed.
	// Will return ErrNotImplemented if not supported
	Watch(GPIOPin, GPIOEdge) error
}

// GPIOEvent implements an event from the GPIO driver
type GPIOEvent interface {
	Event

	// Pin returns the pin on which the event occurred
	Pin() GPIOPin

	// Edge returns whether the pin value is rising or falling
	// or will return NONE if not defined
	Edge() GPIOEdge
}

// I2CDriver implements the I2C interface for sensors, etc.
type I2CDriver interface {
	Driver

	// Set current slave address
	SetSlave(uint8) error

	// Get current slave address
	GetSlave() uint8

	// Return true if a slave was detected at a particular address
	DetectSlave(uint8) (bool, error)

	// Read Byte (8-bits), Word (16-bits) & Block ([]byte) from registers
	ReadUint8(reg uint8) (uint8, error)
	ReadInt8(reg uint8) (int8, error)
	ReadUint16(reg uint8) (uint16, error)
	ReadInt16(reg uint8) (int16, error)
	ReadBlock(reg, length uint8) ([]byte, error)

	// Write Byte (8-bits) & Word (16-bits) to registers
	WriteUint8(reg, value uint8) error
	WriteInt8(reg uint8, value int8) error
	WriteUint16(reg uint8, value uint16) error
	WriteInt16(reg uint8, value int16) error
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// Logical GPIO pin
	GPIOPin uint8

	// GPIO Pin state
	GPIOState uint8

	// GPIO Pin mode
	GPIOMode uint8

	// GPIO Pin resistor configuration (pull up/down or floating)
	GPIOPull uint8

	// GPIOEdge is a rising or falling edge
	GPIOEdge uint8
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Invalid pin constant
	GPIO_PIN_NONE GPIOPin = 0xFF
)

const (
	GPIO_LOW GPIOState = iota
	GPIO_HIGH
)

const (
	// Set pin mode and/or function
	GPIO_INPUT GPIOMode = iota
	GPIO_OUTPUT
	GPIO_ALT5
	GPIO_ALT4
	GPIO_ALT0
	GPIO_ALT1
	GPIO_ALT2
	GPIO_ALT3
	GPIO_NONE
)

const (
	GPIO_PULL_OFF GPIOPull = iota
	GPIO_PULL_DOWN
	GPIO_PULL_UP
)

const (
	GPIO_EDGE_NONE GPIOEdge = iota
	GPIO_EDGE_RISING
	GPIO_EDGE_FALLING
	GPIO_EDGE_BOTH
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p GPIOPin) String() string {
	return fmt.Sprintf("GPIO%v", uint8(p))
}

func (s GPIOState) String() string {
	switch s {
	case GPIO_LOW:
		return "GPIO_LOW"
	case GPIO_HIGH:
		return "GPIO_HIGH"
	default:
		return "[??? Invalid GPIOState value]"
	}
}

func (m GPIOMode) String() string {
	switch m {
	case GPIO_INPUT:
		return "GPIO_INPUT"
	case GPIO_OUTPUT:
		return "GPIO_OUTPUT"
	case GPIO_ALT0:
		return "GPIO_ALT0"
	case GPIO_ALT1:
		return "GPIO_ALT1"
	case GPIO_ALT2:
		return "GPIO_ALT2"
	case GPIO_ALT3:
		return "GPIO_ALT3"
	case GPIO_ALT4:
		return "GPIO_ALT4"
	case GPIO_ALT5:
		return "GPIO_ALT5"
	case GPIO_NONE:
		return "GPIO_NONE"
	default:
		return "[??? Invalid GPIOMode value]"
	}
}

func (p GPIOPull) String() string {
	switch p {
	case GPIO_PULL_OFF:
		return "GPIO_PULL_OFF"
	case GPIO_PULL_DOWN:
		return "GPIO_PULL_DOWN"
	case GPIO_PULL_UP:
		return "GPIO_PULL_UP"
	default:
		return "[??? Invalid GPIOPull value]"
	}
}

func (e GPIOEdge) String() string {
	switch e {
	case GPIO_EDGE_NONE:
		return "GPIO_EDGE_NONE"
	case GPIO_EDGE_RISING:
		return "GPIO_EDGE_RISING"
	case GPIO_EDGE_FALLING:
		return "GPIO_EDGE_FALLING"
	case GPIO_EDGE_BOTH:
		return "GPIO_EDGE_BOTH"
	default:
		return "[??? Invalid GPIOEdge value]"
	}
}
