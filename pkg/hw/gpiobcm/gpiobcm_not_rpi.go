// +build !rpi

package gpiobcm

import (
	"github.com/djthorpe/gopi/v3"
)

type GPIO struct {
	gopi.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *GPIO) New(gopi.Config) error {
	return nil
}

func (this *GPIO) Close() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GPIO) String() string {
	return "<gpiobcm>"
}

////////////////////////////////////////////////////////////////////////////////
// PINS

func (this *GPIO) NumberOfPhysicalPins() uint {
	return 0
}

func (this *GPIO) Pins() []gopi.GPIOPin {
	return nil
}

func (this *GPIO) PhysicalPin(pin uint) gopi.GPIOPin {
	return gopi.GPIO_PIN_NONE
}

func (this *GPIO) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	return 0
}

func (this *GPIO) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	return gopi.GPIO_LOW
}

func (this *GPIO) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {}

func (this *GPIO) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	return gopi.GPIO_NONE
}

func (this *GPIO) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {}

func (this *GPIO) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) error {
	return gopi.ErrNotImplemented
}

func (this *GPIO) Watch(gopi.GPIOPin, gopi.GPIOEdge) error {
	return gopi.ErrNotImplemented
}
