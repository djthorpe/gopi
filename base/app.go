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
	"sort"
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
		// Call configuration for units - don't visit a unit more than once
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
	if err := this.Start([]string{}); err != nil {
		panic(err)
	}
	return 0
}

func (this *App) Start(args []string) error {
	var log gopi.Logger

	if err := this.flags.Parse(args); errors.Is(err, gopi.ErrHelp) {
		this.flags.Usage(os.Stderr)
		return gopi.ErrHelp
	} else if err != nil {
		return err
	} else if this.flags.HasFlag("version", gopi.FLAG_NS_DEFAULT) && this.flags.GetBool("version", gopi.FLAG_NS_DEFAULT) {
		this.flags.Version(os.Stderr)
		return gopi.ErrHelp
	}

	// Create unit instances
	for _, unit := range this.units {
		if unit.New == nil {
			continue
		}
		if log != nil {
			log.Debug("New:", unit)
		}
		if instance, err := unit.New(this); err != nil {
			return err
		} else {
			// `Set logging instance`
			if unit.Type == gopi.UNIT_LOGGER && log == nil {
				if log_, ok := instance.(gopi.Logger); ok {
					log = log_
				}
			}
			// Register instance
			if instance != nil {
				this.instanceByConfig[unit] = instance
			}
		}
	}

	// Call Unit run functions
	for _, unit := range this.units {
		if unit.Run == nil {
			continue
		}
		if log != nil {
			log.Debug("Run:", unit)
		}
		if instance, exists := this.instanceByConfig[unit]; exists {
			if err := unit.Run(this, instance); err != nil {
				return err
			}
		}
	}

	// Success
	return nil
}

func (this *App) Close() error {
	log := this.Log()

	// Close in reverse order
	errs := &gopi.CompoundError{}
	for i := range this.units {
		unit := this.units[len(this.units)-i-1]
		if instance, exists := this.instanceByConfig[unit]; exists {
			if log != nil {
				log.Debug("Close:", unit)
			}
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
		signal.Reset(signals...)		
		return gopi.ErrSignalCaught.WithPrefix(s.String())
	case <-ctx.Done():
		return ctx.Err()
	}
}

////////////////////////////////////////////////////////////////////////////////
// EMIT EVENTS

func (this *App) Emit(e gopi.Event) {
	if bus, exists := this.instancesByName["bus"]; exists && len(bus) > 0 {
		bus[0].(gopi.Bus).Emit(e)
	}
}

////////////////////////////////////////////////////////////////////////////////
// RETURN PROPERTIES

func (this *App) Flags() gopi.Flags {
	return this.flags
}

////////////////////////////////////////////////////////////////////////////////
// RETURN UNIT INSTANCES

func (this *App) Log() gopi.Logger {
	return this.UnitInstance("logger").(gopi.Logger)
}

func (this *App) Timer() gopi.Timer {
	return this.UnitInstance("timer").(gopi.Timer)
}

func (this *App) Bus() gopi.Bus {
	return this.UnitInstance("bus").(gopi.Bus)
}

func (this *App) Platform() gopi.Platform {
	return this.UnitInstance("platform").(gopi.Platform)
}

func (this *App) Display() gopi.Display {
	return this.UnitInstance("display").(gopi.Display)
}

func (this *App) Fonts() gopi.FontManager {
	return this.UnitInstance("fonts").(gopi.FontManager)
}

func (this *App) GPIO() gopi.GPIO {
	return this.UnitInstance("gpio").(gopi.GPIO)
}

func (this *App) I2C() gopi.I2C {
	return this.UnitInstance("i2c").(gopi.I2C)
}

func (this *App) SPI() gopi.SPI {
	return this.UnitInstance("spi").(gopi.SPI)
}

func (this *App) LIRC() gopi.LIRC {
	return this.UnitInstance("lirc").(gopi.LIRC)
}

func (this *App) Surfaces() gopi.SurfaceManager {
	return this.UnitInstance("surfaces").(gopi.SurfaceManager)
}

func (this *App) Input() gopi.InputManager {
	return this.UnitInstance("input").(gopi.InputManager)
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
		pri := make([]uint, 0, len(configs))
		for _, config := range configs {
			if instance, exists := this.instanceByConfig[config]; exists {
				units = append(units, instance)
				pri = append(pri, config.Pri)
			}
		}
		// Sort units by priority field
		sort.Slice(units, func(i, j int) bool {
			return pri[i] > pri[j]
		})
		// Cache unit names
		this.instancesByName[name] = units
		return units
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *App) String() string {
	return fmt.Sprintf("<gopi.App flags=%v instances=%v>", this.flags, this.instanceByConfig)
}
