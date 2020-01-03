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
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Display struct {
	Id       uint
	Platform gopi.Platform
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
	} else if err := this.Init(config); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}
