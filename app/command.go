/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import (
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type command struct {
	base
}

////////////////////////////////////////////////////////////////////////////////
// gopi.App implementation for command-line tool

func NewCommandLineTool(units ...string) (gopi.App, error) {
	this := new(command)

	// Create module instances
	if err := this.base.Init(units); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}
