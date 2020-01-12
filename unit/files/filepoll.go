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

type FilePoll struct {
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (FilePoll) Name() string { return "gopi.filepoll" }

func (config FilePoll) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(filepoll)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}
