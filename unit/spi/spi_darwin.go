// +build darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package spi

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *spi) Init(config SPI) error {
	// SPI not implemented on darwin
	return gopi.ErrNotImplemented
}
