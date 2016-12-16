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
}

// Abstract display interface
type DisplayDriver interface {
	// Enforces general driver
	Driver

	// Return the PPI (pixels-per-inch) for the display, or return zero if
	// display size is unknown
	GetPixelsPerInch() uint32

	// Returns the display size in pixels (width/height)
	GetDisplaySize() (uint32, uint32)

	// Returns the display number
	GetDisplay() uint16
}

// Abstract configuration which is used to open and return the
// concrete driver
type Config interface {
	// Opens the driver from configuration, or returns error
	Open(*util.LoggerDevice) (Driver, error)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// TODO

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
