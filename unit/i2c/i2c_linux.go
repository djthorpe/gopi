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
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type i2c struct {
	bus   uint
	slave uint8
	dev   *os.File
	funcs linux.I2CFunction

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	I2C_SLAVE_NONE uint8 = 0xFF
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (this *i2c) String() string {
	if this.Unit.Closed {
		return fmt.Sprintf("<gopi.I2C>")
	} else {
		return fmt.Sprintf("<gopi.I2C bus=%v funcs=%v>", this.bus, this.funcs)
	}
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *i2c) Init(config I2C) error {
	this.bus = config.Bus
	this.slave = I2C_SLAVE_NONE

	// Open the device
	if dev, err := linux.I2COpenDevice(config.Bus); err != nil {
		return err
	} else {
		this.dev = dev
	}

	// Get functionality
	if funcs, err := linux.I2CFunctions(this.dev.Fd()); err != nil {
		this.dev.Close()
		return err
	} else {
		this.funcs = funcs
	}

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
func (this *i2c) SetSlave(slave uint8) error {
	if err := linux.I2CSetSlave(this.dev.Fd(), slave); err != nil {
		return err
	} else {
		this.slave = slave
		return nil
	}
}

// GetSlave gets current slave address
func (this *i2c) GetSlave() uint8 {
	return this.slave
}

// DetectSlave returns true if a slave was detected at a particular address
func (this *i2c) DetectSlave(slave uint8) (bool, error) {
	detect, err := linux.I2CDetectSlave(this.dev.Fd(), slave, this.funcs)
	if err != nil {
		return false, err
	}
	if this.slave != I2C_SLAVE_NONE {
		if err := this.SetSlave(this.slave); err != nil {
			return false, err
		}
	}
	return detect, nil
}

func (this *i2c) ReadUint8(reg uint8) (uint8, error) {
	if this.slave == I2C_SLAVE_NONE {
		return 0, gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CReadUint8(this.dev.Fd(), reg, this.funcs)
	}
}

func (this *i2c) ReadInt8(reg uint8) (int8, error) {
	if this.slave == I2C_SLAVE_NONE {
		return 0, gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CReadInt8(this.dev.Fd(), reg, this.funcs)
	}
}

func (this *i2c) ReadUint16(reg uint8) (uint16, error) {
	if this.slave == I2C_SLAVE_NONE {
		return 0, gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CReadUint16(this.dev.Fd(), reg, this.funcs)
	}
}

func (this *i2c) ReadInt16(reg uint8) (int16, error) {
	if this.slave == I2C_SLAVE_NONE {
		return 0, gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CReadInt16(this.dev.Fd(), reg, this.funcs)
	}
}

func (this *i2c) ReadBlock(reg, length uint8) ([]byte, error) {
	if this.slave == I2C_SLAVE_NONE {
		return nil, gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CReadBlock(this.dev.Fd(), reg, length, this.funcs)
	}
}

func (this *i2c) WriteUint8(reg, value uint8) error {
	if this.slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CWriteUint8(this.dev.Fd(), reg, value, this.funcs)
	}
}

func (this *i2c) WriteInt8(reg uint8, value int8) error {
	if this.slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CWriteInt8(this.dev.Fd(), reg, value, this.funcs)
	}
}

func (this *i2c) WriteUint16(reg uint8, value uint16) error {
	if this.slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CWriteUint16(this.dev.Fd(), reg, value, this.funcs)
	}
}

func (this *i2c) WriteInt16(reg uint8, value int16) error {
	if this.slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter.WithPrefix("slave")
	} else {
		return linux.I2CWriteInt16(this.dev.Fd(), reg, value, this.funcs)
	}
}
