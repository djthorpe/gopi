// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiolinux

import (
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type GPIO struct {
	FilePoll gopi.FilePoll
	UnexportOnClose bool
}

type gpio struct {
	filepoll gopi.FilePoll

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GPIO_EXPORT   = "/sys/class/gpio/export"
	GPIO_UNEXPORT = "/sys/class/gpio/unexport"
	GPIO_PIN      = "/sys/class/gpio/gpio%v"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (GPIO) Name() string { return "gopi/gpio/linux" }

func (config GPIO) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(gpio)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *gpio) Init(config GPIO) error {
	this.Lock()
	defer this.Unlock()

	if config.FilePoll == nil {
		return gopi.ErrBadParameter.WithPrefix("FilePoll")
	} else {
		this.filepoll = config.FilePoll
	}

	// Return success
	return nil
}

func (this *gpio) Close() error {
	this.Lock()
	defer this.Unlock()

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gpio) String() string {
	if this.Unit.Closed {
		return fmt.Sprintf("<%v>", this.Log.Name())
	} else {
		return fmt.Sprintf("<%v>", this.Log.Name())
	}
}

////////////////////////////////////////////////////////////////////////////////
// PINS

// Return number of physical pins, or 0 if if cannot be returned
// or nothing is known about physical pins
func (this *gpio) NumberOfPhysicalPins() uint {
	return 0
}

// Return array of available logical pins or nil if nothing is
// known about pins
func (this *gpio) Pins() []gopi.GPIOPin {
	return nil
}

// Return logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
// or we don't know about the physical pins
func (this *gpio) PhysicalPin(pin uint) gopi.GPIOPin {
	return gopi.GPIO_PIN_NONE
}

// Return physical pin number for logical pin. Returns 0 where there
// is no physical pin for this logical pin, or we don't know anything
// about the layout
func (this *gpio) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	return 0
}

// Read pin state
func (this *gpio) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	this.Lock()
	defer this.Unlock()

	return 0
}

// Write pin state
func (this *gpio) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	this.Lock()
	defer this.Unlock()

}

// Get pin mode
func (this *gpio) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	this.Lock()
	defer this.Unlock()

	return 0
}

// Set pin mode
func (this *gpio) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	this.Lock()
	defer this.Unlock()

}

// Set pull mode to pull down or pull up - will
// return ErrNotImplemented if not supported
func (this *gpio) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) error {
	this.Lock()
	defer this.Unlock()

	return gopi.ErrNotImplemented
}

// Start watching for rising and/or falling edge,
// or stop watching when GPIO_EDGE_NONE is passed.
// Will return ErrNotImplemented if not supported
func (this *gpio) Watch(gopi.GPIOPin, gopi.GPIOEdge) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS
