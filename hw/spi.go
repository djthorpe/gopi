/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/


// SPI
//
// The abstract SPI hardware interface can be used for interfacing a
// variety of external devices over the SPI interface. In order to use,
// construct an SPI driver object. For any Linux with an SPI driver,
// you can achieve this using a linux.SPI object. For example,
//
//   spi, err := gopi.Open(linux.SPI{ Bus: 1 })
//   if err != nil { /* handle error */ }
//   defer spi.Close()
//
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

	// TODO
}
