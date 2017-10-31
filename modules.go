/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

/*
	This file defines a way to register modules, which is a self-contained
	driver or other piece of code which is created in two phases: a
	configuration phase, giving the module a way to hook in the configuration
	into the application configuration, and a creation phase ("new") which
	reads the configuration and returns the driver (interface gopi.Driver)

	Modules are referenced by type (there are several pre-defined types)
	or by name (where there is no pre-defined type), and can have dependencies
	on other modules.
*/
package gopi // import "github.com/djthorpe/gopi"

import (
	"fmt"
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
	Requires []string
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
	MODULE_TYPE_MDNS
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	modules_by_name = make(map[string]Module)
	modules_by_type = make(map[ModuleType]Module)
	module_name_map = map[string]ModuleType{
		"logger":  MODULE_TYPE_LOGGER,
		"hw":      MODULE_TYPE_HARDWARE,
		"display": MODULE_TYPE_DISPLAY,
		"bitmap":  MODULE_TYPE_BITMAP,
		"vector":  MODULE_TYPE_VECTOR,
		"font":    MODULE_TYPE_VGFONT,
		"opengl":  MODULE_TYPE_OPENGL,
		"layout":  MODULE_TYPE_LAYOUT,
		"gpio":    MODULE_TYPE_GPIO,
		"i2c":     MODULE_TYPE_I2C,
		"spi":     MODULE_TYPE_SPI,
		"input":   MODULE_TYPE_INPUT,
		"mdns":    MODULE_TYPE_MDNS,
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// RegisterModule registers the Config and New functions
// for creating a module, there is no return value
func RegisterModule(module Module) {
	// Check for module.New method
	if module.New == nil {
		panic(fmt.Errorf("Module missing New method: %v", &module))
	}
	// Satisfy module.Type or module.Name
	if module.Type == MODULE_TYPE_OTHER || module.Type == MODULE_TYPE_NONE {
		if module.Name == "" {
			panic(fmt.Errorf("Module name cannot be empty when type is OTHER: %v", &module))
		}
	}
	// Module name cannot be a reserved name
	if _, exists := module_name_map[module.Name]; exists {
		panic(fmt.Errorf("Module name uses reserved word: %v", &module))
	}
	// Register by name
	if module.Name != "" {
		if _, exists := modules_by_name[module.Name]; exists {
			panic(fmt.Errorf("Duplicate Module registered: %v", &module))
		} else {
			modules_by_name[module.Name] = module
		}
	}
	// Register by type if module type is not None or Other
	if module.Type != MODULE_TYPE_OTHER && module.Type != MODULE_TYPE_NONE {
		if _, exists := modules_by_type[module.Type]; exists {
			panic(fmt.Errorf("Duplicate Module registered: %v", &module))
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

// ModuleByName returns a module given the name, or by type
// if it is using the reserved word. It will return nil if
// the module is not registered
func ModuleByName(n string) *Module {
	if t, exists := module_name_map[n]; exists {
		return ModuleByType(t)
	}
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

// ModuleWithDependencies returns an array of pointers to modules
// which satisfy both the module itself and the dependencies. Will
// return an error with the array as nil if the module was not
// found or any dependencies are not met, or there are circular
// dependencies. The ordering of the modules returned is
// important: dependencies are first, and the module requested is
// last, so that they can be initialized in the right order when
// creation is to occur
func ModuleWithDependencies(names ...string) ([]*Module, error) {
	var err error

	// Create modules array
	modules := make([]*Module, 0, len(names))

	// Iterate through the modules
	for _, name := range names {
		// Find module and generate array of dependencies
		if module := ModuleByName(name); module == nil {
			return nil, fmt.Errorf("Module not registered with name: %v", name)
		} else if modules, err = appendModuleAndDependencies(modules, module); err != nil {
			return nil, fmt.Errorf("%v (in module %v)", err, name)
		}
	}

	// Return modules in reverse order to ensure initialization
	// it done in the correct order
	return reverseArray(modules), nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// appendModuleAndDependencies returns array of modules which are dependencies of
// the module provided, with the module appended to the end. Will return nil and
// set err if there is some sort of module not found issue or dependency on self
func appendModuleAndDependencies(arr []*Module, module *Module) ([]*Module, error) {
	var err error

	// Make modules slice with stab at the capacity
	if arr == nil {
		arr = make([]*Module, 0, len(module.Requires)+1)
	}

	// Append module
	arr = append(arr, module)

	// Obtain array of module dependencies
	for _, name := range module.Requires {
		dependency := ModuleByName(name)
		if dependency == nil {
			return nil, fmt.Errorf("Module not registered with name: %v", name)
		}
		if existsModule([]*Module{module}, dependency) {
			return nil, fmt.Errorf("Module cannot depend on self (when satisfying dependencies of %v)", name)
		}
		if existsModule(arr, dependency) {
			return nil, fmt.Errorf("Circular module dependencies (when satisfying dependencies of %v)", name)
		}
		if arr, err = appendModuleAndDependencies(arr, dependency); err != nil {
			return nil, err
		}
	}

	return arr, nil
}

// appendModules will append other modules onto the end of the slice,
// ensuring the module is not already in the slice
func appendModules(modules []*Module, others ...*Module) []*Module {
	for _, other := range others {
		if existsModule(modules, other) == false {
			modules = append(modules, other)
		}
	}
	return modules
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

// In place reversal of the array
func reverseArray(arr []*Module) []*Module {
	last := len(arr) - 1
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[last-i] = arr[last-i], arr[i]
	}
	return arr
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
	case MODULE_TYPE_MDNS:
		return "MODULE_TYPE_MDNS"
	default:
		return "[Invalid ModuleType value]"
	}
}

func (this *Module) String() string {
	return fmt.Sprintf("gopi.Module{ name=\"%v\" type=%v requires=%v }", this.Name, this.Type, this.Requires)
}
