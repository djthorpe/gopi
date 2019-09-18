/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
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
	Run      ModuleRunFunc
	Requires []string
	edges    *module_array
}

// ModuleNewFunc is the signature for creating a new module instance
type ModuleNewFunc func(*AppInstance) (Driver, error)

// ModuleRunFunc is the signature for a finalizer just before running
// the tool
type ModuleRunFunc func(*AppInstance, Driver) error

// ModuleConfigFunc is the signature for setting up the configuration
// for creating the app
type ModuleConfigFunc func(*AppConfig)

// module_array is an internal structure which efficiently allows
// adding and removing of elements
type module_array struct {
	modules    []*Module
	module_map map[*Module]bool
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MODULE_TYPE_NONE ModuleType = iota
	MODULE_TYPE_OTHER
	MODULE_TYPE_LOGGER     // Logging module
	MODULE_TYPE_HARDWARE   // Hardware capabilities and monitoring
	MODULE_TYPE_DISPLAY    // Displays
	MODULE_TYPE_GRAPHICS   // Graphics (Graphics Manager, Surfaces, Bitmaps)
	MODULE_TYPE_SPRITES    // Sprites Manager
	MODULE_TYPE_FONTS      // Font Manager & Faces
	MODULE_TYPE_VECTOR     // 2D Vector graphics
	MODULE_TYPE_OPENGL     // 3D Graphics
	MODULE_TYPE_LAYOUT     // Flex 2D Rectangular Layout
	MODULE_TYPE_GPIO       // GPIO Hardware interface
	MODULE_TYPE_I2C        // I2C Hardware interface
	MODULE_TYPE_SPI        // SPI Hardware interface
	MODULE_TYPE_PWM        // PWM Hardware interface
	MODULE_TYPE_INPUT      // Input manager & devices
	MODULE_TYPE_TIMER      // Timer module
	MODULE_TYPE_LIRC       // LIRC module
	MODULE_TYPE_SERVICE    // RPC Service
	MODULE_TYPE_CLIENT     // RPC Client
	MODULE_TYPE_CLIENTPOOL // RPC Client Pool
	MODULE_TYPE_DISCOVERY  // RPC Service Discovery
	MODULE_TYPE_KEYMAP     // Key Mapper
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	modules_by_name = make(map[string]*Module)
	modules_by_type = make(map[ModuleType]*Module)
	module_name_map = map[string]ModuleType{
		"logger":     MODULE_TYPE_LOGGER,     // Logging
		"hw":         MODULE_TYPE_HARDWARE,   // Platform
		"display":    MODULE_TYPE_DISPLAY,    // Displays
		"graphics":   MODULE_TYPE_GRAPHICS,   // Surface management
		"sprites":    MODULE_TYPE_SPRITES,    // Sprites management
		"fonts":      MODULE_TYPE_FONTS,      // Font management
		"vector":     MODULE_TYPE_VECTOR,     // Vector (2D) graphics
		"opengl":     MODULE_TYPE_OPENGL,     // OpenGL graphics
		"layout":     MODULE_TYPE_LAYOUT,     // Graphics layout
		"gpio":       MODULE_TYPE_GPIO,       // General Purpose input/output
		"i2c":        MODULE_TYPE_I2C,        // I2C Hardware interface
		"spi":        MODULE_TYPE_SPI,        // SPI Hardware interface
		"pwm":        MODULE_TYPE_PWM,        // PWM Hardware interface
		"input":      MODULE_TYPE_INPUT,      // Input devices
		"timer":      MODULE_TYPE_TIMER,      // Timer and backoff
		"lirc":       MODULE_TYPE_LIRC,       // IR Sending and Receiving
		"service":    MODULE_TYPE_SERVICE,    // Microservice service
		"client":     MODULE_TYPE_CLIENT,     // Microservice client
		"clientpool": MODULE_TYPE_CLIENTPOOL, // Microservice client pool
		"discovery":  MODULE_TYPE_DISCOVERY,  // Microservice discovery
		"keymap":     MODULE_TYPE_KEYMAP,     // Keycode/Scancode mapping
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
	// Register by type if module type is not None, Other, Service or Client
	if module.Type != MODULE_TYPE_OTHER && module.Type != MODULE_TYPE_NONE && module.Type != MODULE_TYPE_SERVICE && module.Type != MODULE_TYPE_CLIENT {
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

// ModulesByType returns all modules of a given type
// or an empty array if there are no modules of that
// type
func ModulesByType(t ModuleType) []*Module {
	modules := make([]*Module, 0, 1)
	for _, module := range modules_by_name {
		if module.Type == t {
			modules = append(modules, module)
		}
	}
	return modules
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

// ModuleWithDependencies returns an array of pointers to modules
// which satisfy both the module itself and the dependencies. Will
// return an error with the array as nil if the module was not
// found or any dependencies are not met, or there are circular
// dependencies. The ordering of the modules returned is
// important: dependencies are first, and the module requested is
// last, so that they can be initialized in the right order when
// creation is to occur, and vice-versa on application exit
func ModuleWithDependencies(names ...string) ([]*Module, error) {
	unresolved := newModuleArray()
	resolved := newModuleArray()

	// Iterate through the modules adding the edges to each module
	for _, name := range names {
		if module := ModuleByName(name); module != nil {
			if err := resolveModuleDependencies(module, unresolved, resolved); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Module not registered with name: %v", name)
		}
	}
	return resolved.Array(), nil
}

////////////////////////////////////////////////////////////////////////////////
// module_array implementation

func newModuleArray() *module_array {
	this := new(module_array)
	this.modules = make([]*Module, 0)
	this.module_map = make(map[*Module]bool)
	return this
}

func (this *module_array) Append(module *Module) {
	if _, exists := this.module_map[module]; exists == true {
		return
	}
	this.modules = append(this.modules, module)
	this.module_map[module] = true
}

func (this *module_array) Remove(module *Module) {
	if _, exists := this.module_map[module]; exists == false {
		return
	}
	for i, m := range this.modules {
		if m == module {
			this.modules = append(this.modules[:i], this.modules[i+1:]...)
			break
		}
	}
	delete(this.module_map, module)
}

func (this *module_array) Contains(module *Module) bool {
	_, exists := this.module_map[module]
	return exists
}

func (this *module_array) Array() []*Module {
	return this.modules
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
	// add edges
	module.edges = newModuleArray()
	for _, name := range module.Requires {
		// Find module and generate array of dependencies
		if dependency := ModuleByName(name); dependency == nil {
			return fmt.Errorf("Module not registered with name: %v (required by %v)", name, module.Identifier())
		} else {
			module.edges.Append(dependency)
		}
	}
	// success
	return nil
}

func resolveModuleDependencies(module *Module, unresolved, resolved *module_array) error {
	// Resolve edges if this module hasn't been seen yet
	if module.edges == nil {
		if err := addModuleEdges(module); err != nil {
			return err
		}
	}
	// Mark as unresolved
	unresolved.Append(module)
	// Now resolve each edge as necesary
	for _, edge := range module.edges.Array() {
		if resolved.Contains(edge) == false {
			if unresolved.Contains(edge) {
				return fmt.Errorf("Circular module reference detected: %v => %v", module.Name, edge.Name)
			}
			if err := resolveModuleDependencies(edge, unresolved, resolved); err != nil {
				return err
			}
		}
	}
	// Module has been seen and can be removed from unresolved
	resolved.Append(module)
	unresolved.Remove(module)
	// Return success
	return nil
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
	case MODULE_TYPE_GRAPHICS:
		return "MODULE_TYPE_GRAPHICS"
	case MODULE_TYPE_SPRITES:
		return "MODULE_TYPE_SPRITES"
	case MODULE_TYPE_FONTS:
		return "MODULE_TYPE_FONTS"
	case MODULE_TYPE_VECTOR:
		return "MODULE_TYPE_VECTOR"
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
	case MODULE_TYPE_PWM:
		return "MODULE_TYPE_PWM"
	case MODULE_TYPE_INPUT:
		return "MODULE_TYPE_INPUT"
	case MODULE_TYPE_DISCOVERY:
		return "MODULE_TYPE_DISCOVERY"
	case MODULE_TYPE_TIMER:
		return "MODULE_TYPE_TIMER"
	case MODULE_TYPE_LIRC:
		return "MODULE_TYPE_LIRC"
	case MODULE_TYPE_SERVICE:
		return "MODULE_TYPE_SERVICE"
	case MODULE_TYPE_CLIENT:
		return "MODULE_TYPE_CLIENT"
	case MODULE_TYPE_CLIENTPOOL:
		return "MODULE_TYPE_CLIENTPOOL"
	case MODULE_TYPE_KEYMAP:
		return "MODULE_TYPE_KEYMAP"
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
