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
package gopi

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
	Requires []interface{}
	edges    []*Module
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
	MODULE_TYPE_LOGGER   // Logging module
	MODULE_TYPE_HARDWARE // Hardware capabilities and monitoring
	MODULE_TYPE_DISPLAY  // Displays
	MODULE_TYPE_BITMAP   // Bitmap graphics
	MODULE_TYPE_VECTOR   // 2D Vector graphics
	MODULE_TYPE_VGFONT   // Font rendering
	MODULE_TYPE_OPENGL   // 3D Graphics
	MODULE_TYPE_LAYOUT   // Flex 2D Rectangular Layout
	MODULE_TYPE_GPIO     // GPIO Hardware interface
	MODULE_TYPE_I2C      // I2C Hardware interface
	MODULE_TYPE_SPI      // SPI Hardware interface
	MODULE_TYPE_INPUT    // User Input Device interface
	MODULE_TYPE_MDNS     // DNS Service Discovery
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	modules_by_name = make(map[string]*Module)
	modules_by_type = make(map[ModuleType]*Module)
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
			modules_by_name[module.Name] = &module
		}
	}
	// Register by type if module type is not None or Other
	if module.Type != MODULE_TYPE_OTHER && module.Type != MODULE_TYPE_NONE {
		if _, exists := modules_by_type[module.Type]; exists {
			panic(fmt.Errorf("Duplicate Module registered: %v", &module))
		} else {
			modules_by_type[module.Type] = &module
		}
	}
}

// ModuleByType returns a module given the type. It will
// return nil if the module is not registered
func ModuleByType(t ModuleType) *Module {
	if module, exists := modules_by_type[t]; exists {
		return module
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
		return module
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
// creation is to occur, and vice-versa on application exit
func ModuleWithDependencies(names ...interface{}) ([]*Module, error) {
	// Create modules array
	modules := make([]*Module, 0, len(names))

	// Iterate through the modules adding the edges to each module
	for _, name := range names {
		// Find module and generate array of dependencies
		if module := ModuleByValue(name); module == nil {
			return nil, fmt.Errorf("Module not registered with name: %v", name)
		} else if err := addModuleEdges(module); err != nil {
			return nil, err
		} else {
			modules = append(modules, module)
		}
	}

	// Iterate through the modules again, resolving dependencies
	dependencies := make([]*Module, 0, len(modules))
	for _, module := range modules {
		fmt.Printf("resolving dependencies for %v\n", module.Identifier())
		fmt.Printf(" ...requires=%v\n", module.Requires)
		if resolved, _, err := resolveModuleDependencies(module, nil, nil); err != nil {
			fmt.Printf(" ...returns error %v\n", err)
			return nil, err
		} else {
			fmt.Printf(" ...satisfies dependencies with %v\n", resolved)
			dependencies = append(dependencies, resolved...)
		}
	}

	// Return modules in reverse order to ensure initialization
	// it done in the correct order
	return dependencies, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// addModuleEdges simply resolves the 'Requires' array and initalizes the
// edges member value
func addModuleEdges(module *Module) error {
	// Check for already initialized
	if module.edges != nil {
		return nil
	}
	// make an empty slice with capacity for edges
	module.edges = make([]*Module, 0, len(module.Requires))
	for _, name := range module.Requires {
		// Find module and generate array of dependencies
		if dependency := ModuleByValue(name); dependency == nil {
			return fmt.Errorf("Module not registered with name: %v (required by %v)", name, module.Identifier())
		} else if dependency.In(module.edges) == false {
			module.edges = append(module.edges, dependency)
		}
	}

	return nil
}

func resolveModuleDependencies(module *Module, resolved []*Module, seen []*Module) ([]*Module, []*Module, error) {
	var err error
	if resolved == nil {
		resolved = make([]*Module, 0)
	}
	if seen == nil {
		seen = make([]*Module, 0)
	}
	fmt.Printf("   ....adding %v to seen\n", module.Identifier())
	seen = append(seen, module)
	for _, edge := range module.edges {
		if edge.In(resolved) == false {
			if edge.In(seen) {
				return nil, nil, fmt.Errorf("Circular reference for %v required by %v", edge.Identifier(), module.Identifier())
			}
			fmt.Printf("   ....resolving dependencies for %v\n", edge.Identifier())
			if resolved, seen, err = resolveModuleDependencies(edge, resolved, seen); err != nil {
				return nil, nil, err
			}
		}
	}
	resolved = append(resolved, module)
	fmt.Printf("   ....resolved for %v is %v\n", module.Identifier(), resolved)
	return resolved, seen, nil
}

// Equals returns true if one module is equal to another one
func (this *Module) Equals(other *Module) bool {
	return (this == other)
}

// In returns true if one module is in array of other modules
func (this *Module) In(modules []*Module) bool {
	for _, other := range modules {
		if this.Equals(other) {
			return true
		}
	}
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
	case MODULE_TYPE_MDNS:
		return "MODULE_TYPE_MDNS"
	default:
		return "[Invalid ModuleType value]"
	}
}

func (this *Module) String() string {
	return fmt.Sprintf("%v{ name=\"%v\" type=%v requires=%v }", this.Identifier(), this.Name, this.Type, this.Requires)
}

func (this *Module) Identifier() string {
	if this.Type == MODULE_TYPE_NONE || this.Type == MODULE_TYPE_OTHER {
		return fmt.Sprintf("gopi.Module<%v>", this.Name)
	} else {
		return fmt.Sprintf("gopi.Module.%v<%v>", this.Type, this.Name)
	}
}
