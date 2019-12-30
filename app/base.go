/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import "github.com/djthorpe/gopi/v2"

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type base struct {
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.App

func (this *base) Log() gopi.Logger {
	if logger, ok := this.Unit("logger").(gopi.Logger); ok {
		return logger
	} else {
		return nil
	}
}

func (this *base) Timer() gopi.Timer {
	if timer, ok := this.Unit("timer").(gopi.Timer); ok {
		return timer
	} else {
		return nil
	}
}

func (this *base) Bus() gopi.Bus {
	if bus, ok := this.Unit("bus").(gopi.Bus); ok {
		return bus
	} else {
		return nil
	}
}

func (this *base) Run(main gopi.MainFunc) int {
	if main == nil {
		return -1
	} else if err := main(this); err != nil {
		return -1
	} else {
		return 0
	}
}

func (this *base) Unit(name string) gopi.Unit {
	if units := this.Units(name); len(units) == 0 {
		return nil
	} else {
		return units[0]
	}
}

func (this *base) Units(string) []gopi.Unit {
	// Return units with highest priority one top
	return nil
}
