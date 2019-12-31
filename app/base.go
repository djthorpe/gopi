/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/config"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type base struct {
	flags gopi.Flags
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.App

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *base) Init(name string, units []string) error {
	// Make flags
	if flags := config.NewFlags(filepath.Base(name)); flags == nil {
		return nil
	} else {
		this.flags = flags
	}

	// Get units and dependendies
	units = append([]string{"logger"}, units...)
	if units_, err := gopi.UnitWithDependencies(units...); err != nil {
		return err
	} else {
		// Call configuration for units
		for _, unit := range units_ {
			if unit.Config != nil {
				if err := unit.Config(this); err != nil {
					return fmt.Errorf("%s: %w", unit.Name, err)
				}
			}
		}
	}

	// Success
	return nil
}

func (this *base) Run() int {
	if err := this.flags.Parse(os.Args[1:]); errors.Is(err, gopi.ErrHelp) {
		this.flags.Usage(os.Stderr)
		return -1
	} else if err != nil {
		fmt.Fprintln(os.Stderr, this.flags.Name()+":", err)
		return -1
	} else if this.flags.HasFlag("version", gopi.FLAG_NS_DEFAULT) && this.flags.GetBool("version", gopi.FLAG_NS_DEFAULT) {
		this.flags.Version(os.Stderr)
		return -1
	}

	// Success
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// RETURN PROPERTIES

func (this *base) Flags() gopi.Flags {
	return this.flags
}

func (this *base) Log() gopi.Logger {
	if logger, ok := this.UnitInstance("logger").(gopi.Logger); ok {
		return logger
	} else {
		return nil
	}
}

func (this *base) Timer() gopi.Timer {
	if timer, ok := this.UnitInstance("timer").(gopi.Timer); ok {
		return timer
	} else {
		return nil
	}
}

func (this *base) Bus() gopi.Bus {
	if bus, ok := this.UnitInstance("bus").(gopi.Bus); ok {
		return bus
	} else {
		return nil
	}
}

func (this *base) UnitInstance(name string) gopi.Unit {
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
// STRINGIFY

func (this *base) String() string {
	return fmt.Sprintf("<gopi.App flags=%v>", this.flags)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *base) unitsWithDependencies(unit gopi.UnitConfig) []gopi.UnitConfig {
	return nil
}
