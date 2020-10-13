// +build !linux

package spi

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Devices) Enumerate() []gopi.SPI {
	return nil
}

// Open SPI device
func (this *Devices) Open(bus, slave uint, delay uint16) (gopi.SPI, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Devices) Close(gopi.SPI) error {
	return gopi.ErrNotImplemented
}
