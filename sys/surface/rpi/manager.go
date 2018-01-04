// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type EGL struct {
	Display gopi.Display
}

type egl struct {
	log     gopi.Logger
	display gopi.Display
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config EGL) Open(logger gopi.Logger) (gopi.Driver, err) {
	return nil, gopi.ErrNotImplemented
}
