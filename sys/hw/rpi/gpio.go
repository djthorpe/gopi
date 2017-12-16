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

type GPIO struct{}

type gpio struct {
	log    gopi.Logger
	pins   map[gopi.GPIOPin]*pinstate
	pinmax uint
}

type pinstate struct {
	physical uint
	state    gopi.GPIOState
	mode     gopi.GPIOMode
	pull     gopi.GPIOPull
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// Map logical pins to physical pins
	pinmap = map[uint]gopi.GPIOPin{
		3:  gopi.GPIOPin(2),  // I2C_SDA1      On Rpi1 A/B Rev1: GPIO0 otherwise GPIO2
		5:  gopi.GPIOPin(3),  // I2C_SCL1      On Rpi1 A/B Rev1: GPIO1 otherwise GPIO3
		7:  gopi.GPIOPin(4),  // GPIO_CLK
		8:  gopi.GPIOPin(14), // TXD0
		10: gopi.GPIOPin(15), // RXD0
		11: gopi.GPIOPin(17), // GPIO_GEN0
		12: gopi.GPIOPin(18), // GPIO_GEN1
		13: gopi.GPIOPin(27), // GPIO_GEN2     On Rpi1 A/B Rev1: GPIO21 otherwise GPIO27
		15: gopi.GPIOPin(22), // GPIO_GEN3
		16: gopi.GPIOPin(23), // GPIO_GEN4
		18: gopi.GPIOPin(24), // GPIO_GEN5
		19: gopi.GPIOPin(10), // SPI_MOSI
		21: gopi.GPIOPin(9),  // SPI_MOSO
		22: gopi.GPIOPin(25), // GPIO_GEN6
		23: gopi.GPIOPin(11), // SPI_CLK
		24: gopi.GPIOPin(8),  // SPI_CE0_N
		26: gopi.GPIOPin(7),  // SPI_CE1_N
		29: gopi.GPIOPin(5),  // Not on Rpi1 (all below here)
		31: gopi.GPIOPin(6),
		32: gopi.GPIOPin(12),
		33: gopi.GPIOPin(13),
		35: gopi.GPIOPin(19),
		36: gopi.GPIOPin(16),
		37: gopi.GPIOPin(26),
		38: gopi.GPIOPin(20),
		40: gopi.GPIOPin(21),
	}
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config GPIO) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.mock.GPIO.Open{  }")

	this := new(gpio)
	this.log = logger
	this.pins = make(map[gopi.GPIOPin]*pinstate, len(pinmap))
	this.pinmax = 0

	// Iterate through pins, setting initial state
	for k, v := range pinmap {
		this.pins[v] = &pinstate{
			physical: k,
			state:    gopi.GPIO_LOW,
			mode:     gopi.GPIO_ALT0,
			pull:     gopi.GPIO_PULL_OFF,
		}
		// set highest pin number
		if k > this.pinmax {
			this.pinmax = k
		}
	}

	// Success
	return this, nil
}

// Close
func (this *gpio) Close() error {
	this.log.Debug("sys.mock.GPIO.Close{ }")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gpio) String() string {
	return fmt.Sprintf("sys.mock.GPIO{ number_of_pins=%v }", this.pinmax)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENT INTERFACE

// NumberOfPhysicalPins returns number of physical pins or zero
// if the GPIO interface is not enabled
func (this *gpio) NumberOfPhysicalPins() uint {
	return this.pinmax
}

// Pins() returns array of available logical pins
func (this *gpio) Pins() []gopi.GPIOPin {
	pins := make([]gopi.GPIOPin, 0, len(this.pins))
	for k := range this.pins {
		pins = append(pins, k)
	}
	return pins
}

// PhysicalPin returns logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
func (this *gpio) PhysicalPin(pin uint) gopi.GPIOPin {
	for k, v := range this.pins {
		if v.physical == pin {
			return k
		}
	}
	return gopi.GPIO_PIN_NONE
}

// PhysicalPinForPin returns physical pin number for logical pin.
// Returns 0 where there is no physical pin for this logical pin
func (this *gpio) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	if pin, ok := this.pins[logical]; ok == false {
		return 0
	} else {
		return pin.physical
	}
}

// ReadPin reads pin state or returns LOW otherwise
func (this *gpio) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	if pin, ok := this.pins[logical]; ok {
		return pin.state
	} else {
		this.log.Error("sys.mock.GPIO: ReadPin on invalid logical pin %v", logical)
		return gopi.GPIO_LOW
	}
}

// Write pin state
func (this *gpio) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	if pin, ok := this.pins[logical]; ok {
		pin.state = state
	} else {
		this.log.Error("sys.mock.GPIO: ReadPin on invalid logical pin %v", logical)
	}
}

// Get pin mode
func (this *gpio) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	if pin, ok := this.pins[logical]; ok {
		return pin.mode
	} else {
		this.log.Error("sys.mock.GPIO: GetPinMode on invalid logical pin %v", logical)
		return gopi.GPIO_ALT0
	}
}

// Set pin mode
func (this *gpio) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	if pin, ok := this.pins[logical]; ok {
		pin.mode = mode
	} else {
		this.log.Error("sys.mock.GPIO: SetPinMode on invalid logical pin %v", logical)
	}
}

// Set pull mode
func (this *gpio) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) {
	if pin, ok := this.pins[logical]; ok {
		pin.pull = pull
	} else {
		this.log.Error("sys.mock.GPIO: SetPullMode on invalid logical pin %v", logical)
	}
}
