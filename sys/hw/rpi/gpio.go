// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	// Frameworks
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO struct {
	Hardware gopi.Hardware
}

type gpio struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config GPIO) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.hw.rpi.GPIO.Open{ }")

	this := new(gpio)
	this.log = logger

	// Success
	return this, nil
}

// Close
func (this *gpio) Close() error {
	this.log.Debug("sys.hw.rpi.GPIO.Close{ }")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gpio) String() string {
	return fmt.Sprintf("sys.hw.rpi.GPIO{ }")
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENT INTERFACE

// NumberOfPhysicalPins returns number of physical pins or zero
// if the GPIO interface is not enabled
func (this *gpio) NumberOfPhysicalPins() uint {
	return 0
}

// Pins() returns array of available logical pins
func (this *gpio) Pins() []gopi.GPIOPin {
	return []gopi.GPIOPin{}
}

// PhysicalPin returns logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
func (this *gpio) PhysicalPin(pin uint) gopi.GPIOPin {
	return gopi.GPIO_PIN_NONE
}

// PhysicalPinForPin returns physical pin number for logical pin.
// Returns 0 where there is no physical pin for this logical pin
func (this *gpio) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	return 0
}

// ReadPin reads pin state or returns LOW otherwise
func (this *gpio) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	return gopi.GPIO_LOW
}

// Write pin state
func (this *gpio) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	// TODO
}

// Get pin mode
func (this *gpio) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	return gopi.GPIO_ALT0
}

// Set pin mode
func (this *gpio) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	// TODO
}

// Set pull mode
func (this *gpio) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) {
	// TODO
}
