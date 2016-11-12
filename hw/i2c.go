/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// I2C
//
// The abstract I2C hardware interface can be used for interfacing a
// variety of external devices over the I2C interface. In order to use,
// construct an I2C driver object. For any Linux with an I2C driver,
// you can achieve this using a linux.I2C object. For example,
//
//   i2c, err := gopi.Open(linux.I2C{ Bus: 1 })
//   if err != nil { /* handle error */ }
//   defer i2c.Close()
//
package hw // import "github.com/djthorpe/gopi/hw"

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract I2C interface
type I2CDriver interface {
	// Enforces general driver
	gopi.Driver

	// Set current slave address
	SetSlave(uint8) error

	// Get current slave address
	GetSlave() uint8

	// Return true if a slave was detected at a particular address
	DetectSlave(uint8) (bool, error)

	// Read Byte (8-bits), Word (16-bits) & Block ([]byte) from registers
	ReadUint8(reg uint8) (uint8, error)
	ReadInt8(reg uint8) (int8, error)
	ReadUint16(reg uint8) (uint16, error)
	ReadInt16(reg uint8) (int16, error)
	ReadBlock(reg, length uint8) ([]byte, error)

	// Write Byte (8-bits) & Word (16-bits) to registers
	WriteUint8(reg,value uint8) error
	WriteInt8(reg uint8, value int8) error
	WriteUint16(reg uint8,value uint16) error
	WriteInt16(reg uint8,value int16) error
}
