// +build darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package i2c

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type i2c struct {
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *i2c) Init(config I2C) error {
	// I2C not implemented on darwin
	return gopi.ErrNotImplemented
}
