// +build !linux

package i2c

import (
	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2CFunction uint

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *i2c) Open(bus gopi.I2CBus) (*device, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *i2c) Devices() []gopi.I2CBus {
	return nil
}

func (this *i2c) SetSlave(gopi.I2CBus, uint8) error {
	return gopi.ErrNotImplemented
}

func (this *i2c) DetectSlave(gopi.I2CBus, uint8) (bool, error) {
	return false, gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - READ

func (this *i2c) ReadUint8(bus gopi.I2CBus, reg uint8) (uint8, error) {
	return 0, gopi.ErrNotImplemented
}

func (this *i2c) ReadInt8(bus gopi.I2CBus, reg uint8) (int8, error) {
	return 0, gopi.ErrNotImplemented
}

func (this *i2c) ReadUint16(bus gopi.I2CBus, reg uint8) (uint16, error) {
	return 0, gopi.ErrNotImplemented

}

func (this *i2c) ReadInt16(bus gopi.I2CBus, reg uint8) (int16, error) {
	return 0, gopi.ErrNotImplemented

}

func (this *i2c) ReadBlock(bus gopi.I2CBus, reg, length uint8) ([]byte, error) {
	return nil, gopi.ErrNotImplemented

}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - WRITE

func (this *i2c) WriteUint8(bus gopi.I2CBus, reg, value uint8) error {
	return gopi.ErrNotImplemented

}

func (this *i2c) WriteInt8(bus gopi.I2CBus, reg uint8, value int8) error {
	return gopi.ErrNotImplemented

}

func (this *i2c) WriteUint16(bus gopi.I2CBus, reg uint8, value uint16) error {
	return gopi.ErrNotImplemented

}

func (this *i2c) WriteInt16(bus gopi.I2CBus, reg uint8, value int16) error {
	return gopi.ErrNotImplemented
}
