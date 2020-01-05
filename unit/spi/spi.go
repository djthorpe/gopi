/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package spi

import (
	"github.com/djthorpe/gopi/v2/base"
	// Frameworks
	"os"
	"sync"

	gopi "github.com/djthorpe/gopi/v2"
)

type SPI struct {
	// Bus number
	Bus uint

	// Slave number
	Slave uint

	// Transfer delay between blocks, in microseconds
	Delay uint
}

type spi struct {
	dev           *os.File     // device
	bus           uint         // bus number
	slave         uint         // slave number
	mode          gopi.SPIMode // mode
	speed_hz      uint32       // maximum speed in hertz
	bits_per_word uint8        // bits per word
	delay_usec    uint16       // Transfer delay

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (SPI) Name() string { return "gopi.SPI" }

func (config SPI) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(spi)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}
