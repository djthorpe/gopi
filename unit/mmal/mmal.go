// +build mmal

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mmal

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MMAL struct{}

type mmal struct {
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (MMAL) Name() string { return "gopi/mmal" }

func (config MMAL) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(mmal)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mmal

func (this *mmal) Init(config MMAL) error {
	// Success
	return nil
}

func (this *mmal) Close() error {

	// Success
	return this.Unit.Close()
}
