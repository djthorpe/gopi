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
//   import rpi "./device/rpi"
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

// Store state for the non-abstract input driver
type State struct {
	driver Driver
}

// Abstract driver interface
type Driver interface {
	// Return the logging object
	Log() *util.LoggerDevice

	// Close closes the driver and frees the underlying resources
	Close() error
}

type HardwareDriver interface {
	// Enforces general driver
	Driver

	// Adds display
	Display(DisplayConfig) (DisplayDriver,error)
}

type DisplayDriver interface {
	// Enforces general driver
	Driver
}

// Abstract configuration which is used to open and return the
// concrete driver
type DeviceConfig interface {
	// Opens the driver from configuration, or returns error
	Open(*util.LoggerDevice) (Driver, error)
}

// Abstract display configuration
type DisplayConfig interface {

}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: Config interface implementation

// Open a driver - opens the concrete version given the config method
func Open(config DeviceConfig,log *util.LoggerDevice) (HardwareDriver, error) {
	if log==nil {
		log, err := util.Logger(util.NullLogger{ })
		if err != nil {
			return nil, err
		}
	}
	driver, err := config.Open(log)
	if err != nil {
		return nil, err
	}
	return &State{ driver }, nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: State interface implementation

// Provides human-readable version
func (this *State) String() string {
	return fmt.Sprintf("<gopi>{%v}",this.driver)
}

// Closes the device and frees the resources
func (this *State) Close() error {
	return this.driver.Close()
}

// Return the logging object
func (this *State) Log() *util.LoggerDevice {
	return this.driver.Log()
}

// Returns a display object
func (this *State) Display(config DisplayConfig) (DisplayDriver,error) {
	displaydriver, err := this.driver.(HardwareDriver).Display(config)
	if err != nil {
		return nil, err
	}
	return &State{ displaydriver }, nil
}

