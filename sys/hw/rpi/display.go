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
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Display struct {
	Display       uint
	PixelsPerInch string
}

type display struct {
	log      gopi.Logger
	id       uint
	handle   dxDisplayHandle
	modeinfo *dxDisplayModeInfo
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Display) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.hw.rpi.Display.Open{ id=%v }", config.Display)

	this := new(display)
	this.log = logger
	this.id = config.Display
	this.handle = DX_DISPLAY_NONE

	// Open display
	var err error
	if this.handle, err = dxDisplayOpen(this.id); err != nil {
		return nil, err
	} else if this.modeinfo, err = dxDisplayGetInfo(this.handle); err != nil {
		return nil, err
	}

	//

	// Success
	return this, nil
}

// Close
func (this *display) Close() error {
	this.log.Debug("sys.hw.rpi.Display.Close{ id=%v }", this.id)

	if this.handle != DX_DISPLAY_NONE {
		if err := dxDisplayClose(this.handle); err != nil {
			return err
		} else {
			this.handle = DX_DISPLAY_NONE
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Display returns display number
func (this *display) Display() uint {
	return this.id
}

// Return size
func (this *display) Size() (uint32, uint32) {
	return this.modeinfo.Size.Width, this.modeinfo.Size.Height
}

// Return pixels-per-inch
func (this *display) PixelsPerInch() uint32 {
	// TODO
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *display) String() string {
	return fmt.Sprintf("sys.hw.rpi.Display{ id=%v (%v) info=%v }", dxDisplayId(this.id), this.id, this.modeinfo)
}
