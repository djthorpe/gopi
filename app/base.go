/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import (
	"strconv"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/config"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type base struct {
	flags *config.Flags
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.App

func (this *base) Init(name string, modules []string) error {
	// Make flags
	if flags := config.NewFlags(name); flags == nil {
		return nil
	} else {
		this.flags = flags
	}
	// Configure for logger
	if logger := gopi.UnitsByType(gopi.UNIT_LOGGER); logger == nil {
		return gopi.ErrNotFound.WithPrefix("Missing logger unit")
	} else if logger[0].Config != nil {
		if err := logger[0].Config(this); err != nil {
			return err
		} else if units := this.unitsWithDependencies(logger[0]); len(units) == 0 {

		} else {

		}
	}

	// Get modules and their dependendies
	for _, name := range modules {
		if units := gopi.UnitsByName(name); len(units) == 0 {
			return gopi.ErrNotFound.WithPrefix("Missing " + strconv.Quote(name) + " unit")
		} else {
			for _, unit := range units {
				if unit.Config != nil {
					if err := unit.Config(this); err != nil {
						return err
					}
				}
				if units := this.unitsWithDependencies(unit); len(units) == 0 {
					return gopi.ErrBadParameter.WithPrefix("Unable to satisfy dependencies for " + strconv.Quote(unit.Name) + " unit")
				}
			}
		}
	}

	// Success
	return nil
}

func (this *base) Flags() gopi.Flags {
	return this.flags
}

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

func (this *base) Unit(name string) gopi.Unit {
	if units := this.Units(name); len(units) == 0 {
		return nil
	} else {
		return units[0]
	}
}

func (this *base) Units(string) []gopi.Unit {
	// TODO: Return units with highest priority one top
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *base) unitsWithDependencies(unit gopi.UnitConfig) []gopi.UnitConfig {
	return nil
}
