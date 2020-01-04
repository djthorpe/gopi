// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package i2c

import (
	"fmt"
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2CFunction uint32

type i2c struct {
	bus   uint
	slave uint8
	dev   *os.File
	funcs I2CFunction

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	I2C_DEV                   = "/dev/i2c"
	I2C_SLAVE_NONE      uint8 = 0xFF
	I2C_SMBUS_BLOCK_MAX       = 32 /* As specified in SMBus standard */
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (this *i2c) String() string {
	return fmt.Sprintf("<gopi.I2C bus=%v>", this.bus)
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *i2c) Init(config I2C) error {
	this.bus = config.Bus
	this.slave = I2C_SLAVE_NONE

	// Open the device
	if dev, err := i2c_open_device(config.Bus); err != nil {
		return err
	} else {
		this.dev = dev
	}

	// Get functionality
	/*if funcs, err := this.i2cFuncs(); err != nil {
		this.dev.Close()
		return nil, err
	} else {
		this.funcs = funcs
	}*/

	// success
	return nil
}

// Close I2C connection
func (this *i2c) Close() error {
	if err := this.dev.Close(); err != nil {
		return err
	}

	// Release resources
	this.dev = nil
	this.slave = I2C_SLAVE_NONE

	// Return error
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.I2C

// SetSlave sets current slave address
func (this *i2c) SetSlave(uint8) error {
	return gopi.ErrNotImplemented
}

// GetSlave gets current slave address
func (this *i2c) GetSlave() uint8 {
	return this.slave
}

// DetectSlave returns true if a slave was detected at a particular address
func (this *i2c) DetectSlave(uint8) (bool, error) {
	return false, gopi.ErrNotImplemented
}

func (this *i2c) ReadUint8(reg uint8) (uint8, error) {
	return 0, gopi.ErrNotImplemented
}
func (this *i2c) ReadInt8(reg uint8) (int8, error) {
	return 0, gopi.ErrNotImplemented
}
func (this *i2c) ReadUint16(reg uint8) (uint16, error) {
	return 0, gopi.ErrNotImplemented
}
func (this *i2c) ReadInt16(reg uint8) (int16, error) {
	return 0, gopi.ErrNotImplemented
}
func (this *i2c) ReadBlock(reg, length uint8) ([]byte, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *i2c) WriteUint8(reg, value uint8) error {
	return gopi.ErrNotImplemented
}
func (this *i2c) WriteInt8(reg uint8, value int8) error {
	return gopi.ErrNotImplemented
}
func (this *i2c) WriteUint16(reg uint8, value uint16) error {
	return gopi.ErrNotImplemented
}
func (this *i2c) WriteInt16(reg uint8, value int16) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func i2c_open_device(bus uint) (*os.File, error) {
	if file, err := os.OpenFile(fmt.Sprintf("%v-%v", I2C_DEV, bus), os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}
