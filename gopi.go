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
	"github.com/djthorpe/gopi/util"
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

	// Return display size for nominated display number, or (0,0) if display
	// does not exist
	GetDisplaySize(display uint16) (uint32, uint32)

	// Return serial number of hardware as uint64 - hopefully unique for this device
	GetSerialNumber() (uint64, error)

	// Return system capabilities
	GetCapabilities() []Tuple
}

// Abstract display interface
type DisplayDriver interface {
	// Enforces general driver
	Driver
}

// Abstract configuration which is used to open and return the
// concrete driver
type Config interface {
	// Opens the driver from configuration, or returns error
	Open(*util.LoggerDevice) (Driver, error)
}

// Capability key
type Capability uint

// Abstract set of key/value pairs
type Tuple interface {
	GetKey() Capability
	String() string
}


////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Capability keys
	CAP_HW_SERIAL Capability = iota   // serial number
	CAP_HW_PLATFORM        // platform
	CAP_HW_MODEL           // hardware model number
	CAP_HW_REVISION        // hardware revision
	CAP_HW_PCB             // hardware PCB number
	CAP_HW_WARRANTY        // hardware warranty bit
	CAP_HW_PROCESSOR_NAME  // processor name
	CAP_HW_PROCESSOR_TEMP  // processor temperature
	CAP_MAX                // maximum capability number
)

/*
	GOPI_CAP_DISPLAY            // array of display numbers
	GOPI_CAP_DISPLAY_ID         // display id
	GOPI_CAP_DISPLAY_NAME       // display name
	GOPI_CAP_DISPLAY_WIDTH      // display width
	GOPI_CAP_DISPLAY_HEIGHT     // display height
	GOPI_CAP_DISPLAY_PPI        // display density
	GOPI_CAP_DISPLAY_COLORMODEL // display colormodel
	GOPI_CAP_CLOCK              // clock units
	GOPI_CAP_CLOCK_ID           // speed of each clock
	GOPI_CAP_CLOCK_SPEED        // speed of each clock
	GOPI_CAP_CODEC              // array of enabled codecs
	GOPI_CAP_CODEC_ID           // codec name
	GOPI_CAP_CODEC_ENABLED      // boolean value of whether a codec is enabled
	GOPI_CAP_TEMP               // temperature areas
	GOPI_CAP_TEMP_ID            // temperature name
	GOPI_CAP_TEMP_VALUE         // value of temperature
	GOPI_CAP_MEMORY             // memory units
	GOPI_CAP_MEMORY_VALUE       // value of memory units
*/

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Open a driver - opens the concrete version given the config method
func Open(config Config, log *util.LoggerDevice) (Driver, error) {
	var err error

	if log == nil {
		log, err = util.Logger(util.NullLogger{})
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

func (c Capability) String() string {
	switch(c) {
	case CAP_HW_SERIAL:
		return "CAP_HW_SERIAL"
	case CAP_HW_PLATFORM:
		return "CAP_HW_PLATFORM"
	case CAP_HW_MODEL:
		return "CAP_HW_MODEL"
	case CAP_HW_REVISION:
		return "CAP_HW_REVISION"
	case CAP_HW_PCB:
		return "CAP_HW_PCB"
	case CAP_HW_WARRANTY:
		return "CAP_HW_WARRANTY"
	case CAP_HW_PROCESSOR_NAME:
		return "CAP_HW_PROCESSOR_NAME"
	case CAP_HW_PROCESSOR_TEMP:
		return "CAP_HW_PROCESSOR_TEMP"
	default:
		return "[?? Unknown Capability type]"
	}
}

