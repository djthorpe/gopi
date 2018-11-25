/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
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

	// Return name of the display
	Name() string

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

// I2C implements the I2C interface for sensors, etc.
type I2C interface {
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

// SPI implements the SPI interface for sensors, etc.
type SPI interface {
	Driver

	// Get SPI mode
	Mode() SPIMode
	// Get SPI speed
	MaxSpeedHz() uint32
	// Get Bits Per Word
	BitsPerWord() uint8
	// Set SPI mode
	SetMode(SPIMode) error
	// Set SPI speed
	SetMaxSpeedHz(uint32) error
	// Set Bits Per Word
	SetBitsPerWord(uint8) error

	// Read/Write
	Transfer(send []byte) ([]byte, error)

	// Read
	Read(len uint32) ([]byte, error)

	// Write
	Write(send []byte) error
}

// PWM implements the PWM interface for actuators, motors, etc.
type PWM interface {
	Driver

	// Frequency in Hz
	Frequency(GPIOPin) (float32, error)
	SetFrequency(float32, GPIOPin) error

	// Duty Cycle between 0.0 and 1,0
	DutyCycle(GPIOPin) (float32, error)
	SetDutyCycle(float32, GPIOPin) error
}

// LIRC implements the IR send & receive interface
type LIRC interface {
	Driver
	Publisher

	// Get receive and send modes
	RcvMode() LIRCMode
	SendMode() LIRCMode
	SetRcvMode(mode LIRCMode) error
	SetSendMode(mode LIRCMode) error

	// Receive parameters
	GetRcvResolution() (uint32, error)
	SetRcvTimeout(micros uint32) error
	SetRcvTimeoutReports(enable bool) error
	SetRcvCarrierHz(value uint32) error
	SetRcvCarrierRangeHz(min uint32, max uint32) error

	// Send parameters
	SetSendCarrierHz(value uint32) error
	SetSendDutyCycle(value uint32) error

	// Send Pulse Mode, values are in milliseconds
	PulseSend(values []uint32) error
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

	// SPIMode
	SPIMode uint8

	// LIRCMode
	LIRCMode uint32

	// LIRCType
	LIRCType uint32
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

const (
	SPI_MODE_CPHA SPIMode = 0x01
	SPI_MODE_CPOL SPIMode = 0x02
	SPI_MODE_0    SPIMode = 0x00
	SPI_MODE_1    SPIMode = (0x00 | SPI_MODE_CPHA)
	SPI_MODE_2    SPIMode = (SPI_MODE_CPOL | 0x00)
	SPI_MODE_3    SPIMode = (SPI_MODE_CPOL | SPI_MODE_CPHA)
	SPI_MODE_NONE SPIMode = 0xFF
)

const (
	LIRC_MODE_NONE     LIRCMode = 0x00000000
	LIRC_MODE_RAW      LIRCMode = 0x00000001
	LIRC_MODE_PULSE    LIRCMode = 0x00000002 // send only
	LIRC_MODE_MODE2    LIRCMode = 0x00000004 // rcv only
	LIRC_MODE_LIRCCODE LIRCMode = 0x00000010 // rcv only
	LIRC_MODE_MAX      LIRCMode = LIRC_MODE_LIRCCODE
)

const (
	LIRC_TYPE_SPACE     LIRCType = 0x00000000
	LIRC_TYPE_PULSE     LIRCType = 0x01000000
	LIRC_TYPE_FREQUENCY LIRCType = 0x02000000
	LIRC_TYPE_TIMEOUT   LIRCType = 0x03000000
	LIRC_TYPE_MAX       LIRCType = LIRC_TYPE_TIMEOUT
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

func (m SPIMode) String() string {
	switch m {
	case SPI_MODE_0:
		return "SPI_MODE_0"
	case SPI_MODE_1:
		return "SPI_MODE_1"
	case SPI_MODE_2:
		return "SPI_MODE_2"
	case SPI_MODE_3:
		return "SPI_MODE_3"
	default:
		return "[?? Invalid SPIMode]"
	}
}

func (m LIRCMode) String() string {
	switch m {
	case LIRC_MODE_NONE:
		return "LIRC_MODE_NONE"
	case LIRC_MODE_RAW:
		return "LIRC_MODE_RAW"
	case LIRC_MODE_PULSE:
		return "LIRC_MODE_PULSE"
	case LIRC_MODE_MODE2:
		return "LIRC_MODE_MODE2"
	case LIRC_MODE_LIRCCODE:
		return "LIRC_MODE_LIRCCODE"
	default:
		return "[?? Invalid LIRCMode value]"
	}
}

func (t LIRCType) String() string {
	switch t {
	case LIRC_TYPE_SPACE:
		return "LIRC_TYPE_SPACE"
	case LIRC_TYPE_PULSE:
		return "LIRC_TYPE_PULSE"
	case LIRC_TYPE_FREQUENCY:
		return "LIRC_TYPE_FREQUENCY"
	case LIRC_TYPE_TIMEOUT:
		return "LIRC_TYPE_TIMEOUT"
	default:
		return "[?? Invalid LIRCType value]"
	}
}
