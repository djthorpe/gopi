/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

type FSEvents struct {
	FilePoll gopi.FilePoll
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (FSEvents) Name() string { return "gopi/fsevents" }

func (config FSEvents) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(fsevents)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}
