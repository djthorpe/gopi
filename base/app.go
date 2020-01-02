/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package base

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/config"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type App struct {
	sync.Mutex

	flags            gopi.Flags
	units            []*gopi.UnitConfig
	instanceByConfig map[*gopi.UnitConfig]gopi.Unit
	instancesByName  map[string][]gopi.Unit
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *App) Init(name string, units []string) error {
	// Make flags
	if flags := config.NewFlags(name); flags == nil {
		return nil
	} else {
		this.flags = flags
	}

	// Get units and dependendies
	units = append([]string{"logger"}, units...)
	if units_, err := gopi.UnitWithDependencies(units...); err != nil {
		return err
	} else {
		// Call configuration for units - don't visit a unit more
		// than once
		unitmap := make(map[*gopi.UnitConfig]bool)
		this.units = make([]*gopi.UnitConfig, 0, len(units_))
		for _, unit := range units_ {
			if _, exists := unitmap[unit]; exists {
				continue
			} else if unit.Config != nil {
				if err := unit.Config(this); err != nil {
					return fmt.Errorf("%s: %w", unit.Name, err)
				}
			}
			this.units = append(this.units, unit)
			unitmap[unit] = true
		}
		// Set units and instances map
		this.instanceByConfig = make(map[*gopi.UnitConfig]gopi.Unit, len(this.units))
		this.instancesByName = make(map[string][]gopi.Unit, len(this.units))
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.App

func (this *App) Run() int {
	return this.Start([]string{})
}

func (this *App) Start(args []string) int {
	if err := this.flags.Parse(args); errors.Is(err, gopi.ErrHelp) {
		this.flags.Usage(os.Stderr)
		return -1
	} else if err != nil {
		fmt.Fprintln(os.Stderr, this.flags.Name()+":", err)
		return -1
	} else if this.flags.HasFlag("version", gopi.FLAG_NS_DEFAULT) && this.flags.GetBool("version", gopi.FLAG_NS_DEFAULT) {
		this.flags.Version(os.Stderr)
		return -1
	}
	// Create unit instances
	for _, unit := range this.units {
		if unit.New == nil {
			continue
		}
		if instance, err := unit.New(this); err != nil {
			fmt.Fprintln(os.Stderr, unit.Name+":", err)
			return -1
		} else {
			if instance != nil {
				this.instanceByConfig[unit] = instance
			}
		}
	}

	// Success
	return 0
}

func (this *App) Close() error {
	// Close in reverse order
	errs := &gopi.CompoundError{}
	for i := range this.units {
		unit := this.units[len(this.units)-i-1]
		if instance, exists := this.instanceByConfig[unit]; exists {
			errs.Add(instance.Close())
		}
	}

	// Release resources
	this.flags = nil
	this.units = nil
	this.instanceByConfig = nil
	this.instancesByName = nil

	// Return success
	return errs.ErrorOrSelf()
}

func (this *App) WaitForSignal(ctx context.Context, signals ...os.Signal) error {
	sigchan := make(chan os.Signal, 1)
	defer close(sigchan)

	signal.Notify(sigchan, signals...)
	select {
	case s := <-sigchan:
		return gopi.ErrSignalCaught.WithPrefix(s.String())
	case <-ctx.Done():
		return ctx.Err()
	}
}

////////////////////////////////////////////////////////////////////////////////
// RETURN PROPERTIES

func (this *App) Flags() gopi.Flags {
	return this.flags
}

func (this *App) Log() gopi.Logger {
	if logger, ok := this.UnitInstance("logger").(gopi.Logger); ok {
		return logger
	} else {
		return nil
	}
}

func (this *App) Timer() gopi.Timer {
	if timer, ok := this.UnitInstance("timer").(gopi.Timer); ok {
		return timer
	} else {
		return nil
	}
}

func (this *App) Bus() gopi.Bus {
	if bus, ok := this.UnitInstance("bus").(gopi.Bus); ok {
		return bus
	} else {
		return nil
	}
}

func (this *App) Platform() gopi.Platform {
	if platform, ok := this.UnitInstance("platform").(gopi.Platform); ok {
		return platform
	} else {
		return nil
	}
}

func (this *App) UnitInstance(name string) gopi.Unit {
	if units := this.UnitInstancesByName(name); len(units) == 0 {
		return nil
	} else {
		return units[0]
	}
}

func (this *App) UnitInstancesByName(name string) []gopi.Unit {
	// Cached unit names
	if units, exists := this.instancesByName[name]; exists {
		return units
	}
	// Otherwise, get configurations by name and match with
	// configurations for this applicatiomn
	if configs := gopi.UnitsByName(name); len(configs) == 0 {
		return nil
	} else {
		units := make([]gopi.Unit, 0, len(configs))
		for _, config := range configs {
			if instance, exists := this.instanceByConfig[config]; exists {
				units = append(units, instance)
			}
		}
		// TODO: Sort units by some sort of priority field
		this.instancesByName[name] = units
		return units
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *App) String() string {
	return fmt.Sprintf("<gopi.App flags=%v instances=%v>", this.flags, this.instanceByConfig)
}
