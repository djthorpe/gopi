// +build !rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package display

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type display struct {
	id       uint
	platform gopi.Platform

	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (this *display) String() string {
	return fmt.Sprintf("<gopi.Display id=%v>", this.id)
}

func (this *display) Init(config Display) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Display

// Return display number
func (this *display) DisplayId() uint {
	return this.id
}

// Return name of the display
func (this *display) Name() string {
	return ""
}

// Return display size for nominated display number, or (0,0) if display does not exist
func (this *display) Size() (uint32, uint32) {
	return 0, 0
}

// Return the PPI (pixels-per-inch) for the display, or return zero if unknown
func (this *display) PixelsPerInch() uint32 {
	return 0
}
