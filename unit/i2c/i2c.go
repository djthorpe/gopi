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
)

type I2C struct {
	Bus uint
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (I2C) Name() string { return "gopi.I2C" }

func (config I2C) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(i2c)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}
