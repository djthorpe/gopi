/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// ModuleType defines the type of module
type ModuleType uint

// Module is a structure which determines details about a module
type Module struct {
	Name   string
	Type   ModuleType
	Config ModuleConfigFunc
	New    ModuleNewFunc
}

// ModuleConfigFunc is the signature for altering the configuration before
// the application instance is made
type ModuleConfigFunc func(*AppConfig)

// ModuleNewFunc is the signature for creating a new module instance
type ModuleNewFunc func(*AppConfig) (Driver, error)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MODULE_TYPE_NONE ModuleType = iota
	MODULE_TYPE_OTHER
	MODULE_TYPE_LOGGER
	MODULE_TYPE_HARDWARE
	MODULE_TYPE_DISPLAY
	MODULE_TYPE_BITMAP
	MODULE_TYPE_VECTOR
	MODULE_TYPE_VGFONT
	MODULE_TYPE_OPENGL
	MODULE_TYPE_GPIO
	MODULE_TYPE_I2C
	MODULE_TYPE_SPI
	MODULE_TYPE_INPUT
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	modules_by_name = make(map[string]Module)
	modules_by_type = make(map[ModuleType]Module)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func RegisterModule(module Module) {
	// Register by name
	if _, exists := modules_by_name[module.Name]; exists {
		fmt.Fprintln(os.Stderr, "Duplicate Module registered:", &module)
		os.Exit(-1)
	} else {
		modules_by_name[module.Name] = module
	}
	// Register by type if module type is not None or Other
	if module.Type != MODULE_TYPE_OTHER && module.Type != MODULE_TYPE_NONE {
		if _, exists := modules_by_type[module.Type]; exists {
			fmt.Fprintln(os.Stderr, "Duplicate Module registered:", &module)
			os.Exit(-1)
		} else {
			modules_by_type[module.Type] = module
		}
	}
}

func ModuleByType(t ModuleType) (*Module, error) {
	if module, exists := modules_by_type[t]; exists {
		return &module, nil
	} else {
		return nil, ErrModuleNotFound
	}
}

func ModuleByName(n string) (*Module, error) {
	if module, exists := modules_by_name[n]; exists {
		return &module, nil
	} else {
		return nil, ErrModuleNotFound
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t ModuleType) String() string {
	switch t {
	case MODULE_TYPE_NONE:
		return "MODULE_TYPE_NONE"
	case MODULE_TYPE_OTHER:
		return "MODULE_TYPE_OTHER"
	case MODULE_TYPE_LOGGER:
		return "MODULE_TYPE_LOGGER"
	case MODULE_TYPE_HARDWARE:
		return "MODULE_TYPE_HARDWARE"
	case MODULE_TYPE_DISPLAY:
		return "MODULE_TYPE_DISPLAY"
	case MODULE_TYPE_BITMAP:
		return "MODULE_TYPE_BITMAP"
	case MODULE_TYPE_VECTOR:
		return "MODULE_TYPE_VECTOR"
	case MODULE_TYPE_VGFONT:
		return "MODULE_TYPE_VGFONT"
	case MODULE_TYPE_OPENGL:
		return "MODULE_TYPE_OPENGL"
	case MODULE_TYPE_GPIO:
		return "MODULE_TYPE_GPIO"
	case MODULE_TYPE_I2C:
		return "MODULE_TYPE_I2C"
	case MODULE_TYPE_SPI:
		return "MODULE_TYPE_SPI"
	case MODULE_TYPE_INPUT:
		return "MODULE_TYPE_INPUT"
	default:
		return "[Invalid ModuleType value]"
	}
}

func (this *Module) String() string {
	return fmt.Sprintf("gopi.Module{ name=%v type=%v }", this.Name, this.Type)
}
