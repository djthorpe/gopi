/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package base

import (
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Unit is the base struct for any unit
type Unit struct {
	Log    gopi.Logger
	Closed bool
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (this *Unit) Init(log gopi.Logger) error {
	this.Log = log
	return nil
}

func (this *Unit) Close() error {
	if this.Closed {
		return gopi.ErrInternalAppError.WithPrefix("Close called twice")
	} else {
		this.Log = nil
		this.Closed = true
		return nil
	}
}

func (this *Unit) String() string {
	if this.Log != nil {
		return "<" + this.Log.Name() + ">"
	} else {
		return "<gopi.Unit>"
	}
}
