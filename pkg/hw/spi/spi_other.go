// +build !linux

package spi

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *spi) Devices() []gopi.SPIBus {
	return nil
}

func (this *spi) Mode(gopi.SPIBus) gopi.SPIMode {
	return gopi.SPI_MODE_NONE
}

func (this *spi) SetMode(gopi.SPIBus, gopi.SPIMode) error {
	return gopi.ErrNotImplemented
}

func (this *spi) MaxSpeedHz(gopi.SPIBus) uint32 {
	return 0
}

func (this *spi) SetMaxSpeedHz(gopi.SPIBus, uint32) error {
	return gopi.ErrNotImplemented
}

func (this *spi) BitsPerWord(gopi.SPIBus) uint8 {
	return 0
}

func (this *spi) SetBitsPerWord(gopi.SPIBus, uint8) error {
	return gopi.ErrNotImplemented
}

func (this *spi) Transfer(gopi.SPIBus, []byte) ([]byte, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *spi) Read(gopi.SPIBus, []byte) error {
	return gopi.ErrNotImplemented
}

func (this *spi) Write(gopi.SPIBus, []byte) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// MOCK DEVICE

type device struct{}

func NewDevice(gopi.SPIBus, uint16) (*device, error) {
	return nil, gopi.ErrNotImplemented
}

func (*device) Close() error {
	return gopi.ErrNotImplemented
}
