// +build !linux

package spi

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Devices) Enumerate() []Device {
	return nil
}

// Open SPI device
func (this *Devices) Open(dev Device, delay uint16) (gopi.SPI, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Devices) Close(dev Device) error {
	return gopi.ErrNotImplemented
}
