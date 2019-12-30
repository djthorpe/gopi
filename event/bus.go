/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package event

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Bus struct{}

type bus struct {
	gopi.UnitBase
}

///////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Bus) Name() string { return "gopi.Bus" }

func (config Bus) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(bus)
	if err := this.Init(log); err != nil {
		return nil, err
	}
	return this, nil
}

///////////////////////////////////////////////////////////////////////////////
// gopi.EventBus

/*
func (this *bus) AddHandler(string, gopi.EventHandler) {
	// TODO
}

func (this *bus) Emit(gopi.Event) {
	// TODO
}

*/
