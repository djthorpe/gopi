/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// Package gopi implements a Golang interface for the Raspberry Pi. It's
// a bag of interfaces to both the Broadcom ARM processor and to various
// devices which can be plugged in and interfaced via various interfaces.
//
// Start by creating a gopi object as follows:
//
//   import "github.com/djthorpe/gopi/device/rpi"
//
//   device, err := gopi.Open(rpi.Device{ /* configuration */ })
//   if err != nil { /* handle error */ }
//
// You should then have an object which can be used to retrieve information
// about the Raspberry Pi (serial number, memory configuration, temperature,
// and so forth). When you're done with the object you should release the
// resources using the following method:
//
//   if err := device.Close(); err != nil { /* handle error */ }
//
// You'll need a configuration from the "concrete device". In this case, it's
// the Raspberry Pi device driver. The use these abstrat interfaces will allow
// for other devices to be implemented a bit differently in the future.
//
package gopi // import "github.com/djthorpe/gopi"

import (
	"fmt"
)

import (
	"./util" // import "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract driver interface
type Driver interface {
	// Close closes the driver and frees the underlying resources
	Close() error
}

// Abstract hardware interface - this assumes the hardware has a display
type HardwareDriver interface {
	// Enforces general driver
	Driver

	// Return display size for nominated display number
	GetDisplaySize(display uint16) (uint32, uint32)

	// Return serial number of hardware
	GetSerialNumber() (uint64, error)
}

// Abstract display interface
type DisplayDriver interface {
	// Enforces general driver
	Driver
}

// Abstract GPIO interface
type GPIODriver interface {
	// Enforces general driver
	Driver

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
	WritePin(GPIOPin,GPIOState)

	// Get pin mode
	GetPinMode(GPIOPin) GPIOMode

	// Set pin mode
	SetPinMode(GPIOPin,GPIOMode)

}

// Abstract configuration which is used to open and return the
// concrete driver
type Config interface {
	// Opens the driver from configuration, or returns error
	Open(*util.LoggerDevice) (Driver, error)
}

// GPIO types
type (
	GPIOPin uint8
	GPIOState uint8
	GPIOMode uint8
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
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

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: Config interface implementation

// Open a driver - opens the concrete version given the config method
func Open(config Config,log *util.LoggerDevice) (Driver, error) {
	var err error
	
	if log==nil {
		log, err = util.Logger(util.NullLogger{ })
		if err != nil {
			return nil, err
		}
	}
	driver, err := config.Open(log)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (p GPIOPin) String() string {
	return fmt.Sprintf("GPIO%v",uint8(p))
}

func (s GPIOState) String() string {
	switch(s) {
	case GPIO_LOW:
		return "LOW"
	case GPIO_HIGH:
		return "HIGH"
	default:
		return "[??? Invalid GPIOState value]"
	}
}

func (m GPIOMode) String() string {
	switch(m) {
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

