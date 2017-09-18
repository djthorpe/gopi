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
	Name     string
	Type     ModuleType
	Config   ModuleConfigFunc
	New      ModuleNewFunc
	Requires []interface{}
}

// ModuleNewFunc is the signature for creating a new module instance
type ModuleNewFunc func(*AppInstance) (Driver, error)

// ModuleConfigFunc is the signature for setting up the configuration
// for creating the app
type ModuleConfigFunc func(*AppConfig)

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
	MODULE_TYPE_LAYOUT
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

// RegisterModule registers the Config and New functions
// for creating a module, there is no return value
func RegisterModule(module Module) {
	// Check for module.New method
	if module.New == nil {
		fmt.Fprintln(os.Stderr, "Module missing New method:", &module)
		os.Exit(-1)
	}
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

// ModuleByType returns a module given the type. It will
// return nil if the module is not registered
func ModuleByType(t ModuleType) *Module {
	if module, exists := modules_by_type[t]; exists {
		return &module
	} else {
		return nil
	}
}

// ModuleByName returns a module given the name. It will
// return nil if the module is not registered
func ModuleByName(n string) *Module {
	if module, exists := modules_by_name[n]; exists {
		return &module
	} else {
		return nil
	}
}

// ModuleByValue returns a module given either a type
// or a name. It will return nil if the module is
// not registered
func ModuleByValue(v interface{}) *Module {
	switch v.(type) {
	case string:
		return ModuleByName(v.(string))
	case ModuleType:
		return ModuleByType(v.(ModuleType))
	default:
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// appendModule will append a module to the array, satisfying dependencies
// by appending them to the array as well. It won't currently detect any
// endless recursion, but should probably do that
func appendModule(modules []*Module, name interface{}) ([]*Module, error) {
	var err error

	// Create modules array
	if modules == nil {
		modules = make([]*Module, 0, 1)
	}
	// Find module, and satisfy dependencies
	if module := ModuleByValue(name); module == nil {
		return nil, fmt.Errorf("Module not found: %v", name)
	} else {
		// Satisfy dependencies
		if len(module.Requires) > 0 {
			for _, dependency := range module.Requires {
				if modules, err = appendModule(modules, dependency); err != nil {
					return nil, err
				}
			}
		}
		// Append module if it doesn't exist in the list of modules
		if existsModule(modules, module) == false {
			modules = append(modules, module)
		}
	}
	// Return modules
	return modules, nil
}

// existsModule returns true if the module already exists
// in the array of modules by type (or if MODULE_TYPE_OTHER)
// then by name
func existsModule(modules []*Module, module *Module) bool {
	for _, other := range modules {
		if module.Type == MODULE_TYPE_OTHER || module.Type == MODULE_TYPE_NONE {
			if module.Name == other.Name {
				return true
			}
		} else {
			if module.Type == other.Type {
				return true
			}
		}
	}
	// No module found
	return false
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
	case MODULE_TYPE_LAYOUT:
		return "MODULE_TYPE_LAYOUT"
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
	return fmt.Sprintf("gopi.Module{ name=\"%v\" type=%v }", this.Name, this.Type)
}
