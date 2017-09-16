/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"

// HardwareDriver2 implements the hardware driver interface, which
// provides information about the hardware that the software is
// running on
type HardwareDriver2 interface {
	Driver

	// Return name of the hardware platform
	Name() string

	// Return unique serial number of this hardware
	SerialNumber() string

	// Return the number of possible displays for this hardware
	NumberOfDisplays() uint
}

// DisplayDriver2 implements a pixel-based display device description
type DisplayDriver2 interface {
	Driver

	// Return display number
	Display() uint

	// Return display size for nominated display number, or (0,0) if display
	// does not exist
	Size() (uint32, uint32)

	// Return the PPI (pixels-per-inch) for the display, or return zero if unknown
	PixelsPerInch() uint32
}

// I2CDriver implements the I2C interface for sensors, etc.
type I2CDriver interface {
	Driver

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
	WriteUint8(reg, value uint8) error
	WriteInt8(reg uint8, value int8) error
	WriteUint16(reg uint8, value uint16) error
	WriteInt16(reg uint8, value int16) error
}
