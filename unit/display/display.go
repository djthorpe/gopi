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
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type Display struct {
	Id       uint
	Platform gopi.Platform
}

type display struct {
	id       uint
	platform gopi.Platform

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Display) Name() string { return "gopi.Display" }

func (config Display) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(display)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if config.Platform == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("Platform")
	} else if config.Platform.NumberOfDisplays() == 0 {
		return nil, fmt.Errorf("No displays available on platform")
	} else if config.Id >= config.Platform.NumberOfDisplays() {
		return nil, gopi.ErrBadParameter.WithPrefix("Id")
	} else {
		this.platform = config.Platform
		this.id = config.Id
	}
	return this, nil
}

func (this *display) String() string {
	return fmt.Sprintf("<gopi.Display id=%v name=%v>", this.id, strconv.Quote(this.Name()))
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Display

// Return display number
func (this *display) DisplayId() uint {
	return this.id
}

// Return name of the display
func (this *display) Name() string {
	return "UNKNOWN"
}

// Return display size for nominated display number, or (0,0) if display does not exist
func (this *display) Size() (uint32, uint32) {
	return 0, 0
}

// Return the PPI (pixels-per-inch) for the display, or return zero if unknown
func (this *display) PixelsPerInch() uint32 {
	return 0
}
