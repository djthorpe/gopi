// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package display

import (
	"fmt"
	"strconv"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type display struct {
	id       uint
	platform gopi.Platform
	handle   rpi.DXDisplayHandle
	modeinfo rpi.DXDisplayModeInfo

	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (this *display) String() string {
	if this.handle == 0 {
		return fmt.Sprintf("<gopi.Display id=%v>", this.id)
	} else {
		return fmt.Sprintf("<gopi.Display id=%v name=%v info=%v>", this.id, strconv.Quote(this.Name()), this.modeinfo)
	}
}

func (this *display) Init(config Display) error {
	if handle, err := rpi.DXDisplayOpen(rpi.DXDisplayId(this.id)); err != nil {
		return err
	} else if modeinfo, err := rpi.DXDisplayGetInfo(handle); err != nil {
		rpi.DXDisplayClose(handle)
		return err
	} else {
		this.handle = handle
		this.modeinfo = modeinfo
	}
	// Success
	return nil
}

func (this *display) Close() error {
	if this.handle != 0 {
		if err := rpi.DXDisplayClose(this.handle); err != nil {
			return err
		}
	}

	// Release resources
	this.handle = 0

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Display

// Return display number
func (this *display) DisplayId() uint {
	return this.id
}

// Return name of the display
func (this *display) Name() string {
	return fmt.Sprint(rpi.DXDisplayId(this.id))
}

// Return display size for nominated display number, or (0,0) if display does not exist
func (this *display) Size() (uint32, uint32) {
	return this.modeinfo.Size.W, this.modeinfo.Size.H
}

// Return the PPI (pixels-per-inch) for the display, or return zero if unknown
func (this *display) PixelsPerInch() uint32 {
	return 0
}