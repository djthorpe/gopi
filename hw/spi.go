/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package hw // import "github.com/djthorpe/gopi/hw"

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract SPI interface
type SPIDriver interface {
	// Enforces general driver
	gopi.Driver

	// Set SPI mode
	SetMode(SPIMode) error

	// Get SPI mode
	GetMode() SPIMode

	// Set SPI speed
	SetMaxSpeedHz(uint32) error

	// Get SPI speed
	GetMaxSpeedHz() uint32

	// Set Bits Per Word
	SetBitsPerWord(uint8) error

	// Get Bits Per Word
	GetBitsPerWord() uint8

	// Read/Write
	Transfer(send []byte) ([]byte, error)

	// Read
	Read(len uint32) ([]byte, error)

	// Write
	Write(send []byte) error
}

type SPIMode uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SPI_MODE_CPHA SPIMode = 0x01
	SPI_MODE_CPOL SPIMode = 0x02
	SPI_MODE_0    SPIMode = 0x00
	SPI_MODE_1    SPIMode = (0x00 | SPI_MODE_CPHA)
	SPI_MODE_2    SPIMode = (SPI_MODE_CPOL | 0x00)
	SPI_MODE_3    SPIMode = (SPI_MODE_CPOL | SPI_MODE_CPHA)
	SPI_MODE_NONE SPIMode = 0xFF
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (m SPIMode) String() string {
	switch m {
	case SPI_MODE_0:
		return "SPI_MODE_0"
	case SPI_MODE_1:
		return "SPI_MODE_1"
	case SPI_MODE_2:
		return "SPI_MODE_2"
	case SPI_MODE_3:
		return "SPI_MODE_3"
	default:
		return "[?? Invalid SPIMode]"
	}
}
