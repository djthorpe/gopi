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
	"runtime"
	"strings"
	"sync"

	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// AppConfig defines how an application should be created
type AppConfig struct {
	LogLevel util.LogLevel
	Modules  []ModuleType
	AppFlags *Flags
	Debug    bool
	Verbose  bool
}

// AppInstance defines the running application instance with modules
type AppInstance struct {
	AppFlags *Flags
	Logger   Logger
	Hardware HardwareDriver2
	Display  DisplayDriver2
	Bitmap   Driver
	Vector   Driver
	VGFont   Driver
	OpenGL   Driver
	Layout   Layout
	GPIO     Driver
	I2C      Driver
	SPI      Driver
	Input    Driver

	debug   bool
	verbose bool
}

// Task defines a function which can run, and has a channel which
// indicates when the main thread has finished
type Task func(app *AppInstance, done chan struct{}) error

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// DONE is the message sent on the channel to indicate task is completed
	DONE = struct{}{}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// NewAppConfig method will create a new configuration file given the set of
// modules which should be created
func NewAppConfig(modules ...ModuleType) AppConfig {
	config := AppConfig{}

	// Only allow each module type once, and don't allow NONE or OTHER
	// plus add logger as a mandatory module later
	module_hash := make(map[ModuleType]bool, len(modules))
	for _, t := range modules {
		if t == MODULE_TYPE_NONE || t == MODULE_TYPE_OTHER || t == MODULE_TYPE_LOGGER {
			continue
		}
		module_hash[t] = true
	}

	// Now enumerate the modules, we always have the LOGGER type first
	config.Modules = make([]ModuleType, 1, len(module_hash)+1)
	config.Modules[0] = MODULE_TYPE_LOGGER
	for k := range module_hash {
		config.Modules = append(config.Modules, k)
	}

	// Set the flags
	config.AppFlags = flags
	config.Debug = false
	config.Verbose = false

	// Return the configuration
	return config
}

// NewAppInstance method will create a new application object given an application
// configuration
func NewAppInstance(config AppConfig) (*AppInstance, error) {

	// Parse flags. We want to ignore flags which start with "-test."
	// in the testing environment
	if config.AppFlags != nil && config.AppFlags.Parsed() == false {
		if err := config.AppFlags.Parse(getTestlessArguments(os.Args[1:])); err != nil {
			return nil, err
		}
	}

	// Set debug and verbose flags
	if debug, exists := config.AppFlags.GetBool("debug"); exists {
		config.Debug = debug
	}
	if verbose, exists := config.AppFlags.GetBool("verbose"); exists {
		config.Verbose = verbose
	}

	// Create instance
	this := new(AppInstance)
	this.debug = config.Debug
	this.verbose = config.Verbose
	this.AppFlags = config.AppFlags

	// Create subsystems
	var once sync.Once
	for _, t := range config.Modules {
		// Report open (once after logger module is created)
		if this.Logger != nil {
			once.Do(func() {
				this.Logger.Debug2("gopi.AppInstance.Open()")
			})
		}
		if module, err := ModuleByType(t); err != nil {
			return nil, err
		} else if driver, err := module.New(&config, this.Logger); err != nil {
			return nil, err
		} else if err := this.setModuleInstance(t, driver); err != nil {
			return nil, err
		}
	}

	return this, nil
}

// Run all tasks simultaneously, the first task in the list on the main thread and the
// remaining tasks elsewhere.
func (this *AppInstance) Run(tasks ...Task) error {
	// Lock this to run in the current operating system thread (ie, the main thread)
	runtime.LockOSThread()

	// if no tasks then return
	if len(tasks) == 0 {
		return ErrNoTasks
	}

	// create the channels we'll use to signal the goroutines
	channels := make([]chan struct{}, len(tasks))
	for i := range tasks {
		channels[i] = make(chan struct{})
	}

	// if more than one task, then give them a channel which is signalled
	// by the main thread for ending
	var wg sync.WaitGroup
	if len(tasks) > 1 {
		for i, task := range tasks[1:] {
			wg.Add(1)
			go func(i int, t Task) {
				defer wg.Done()
				if err := t(this, channels[i+1]); err != nil {
					this.Logger.Error("Error: %v [task %v]", err, i+1)
				}
			}(i, task)
		}
	}

	go func() {
		// Wait for mainDone
		_ = <-channels[0]
		if this.Logger != nil {
			this.Logger.Debug2("Main thread done")
		}
		// Signal other tasks to complete
		for i := 1; i < len(tasks); i++ {
			channels[i] <- DONE
		}
	}()

	// Now run main task
	err := tasks[0](this, channels[0])

	// Wait for other tasks to finish
	if len(tasks) > 1 {
		if this.Logger != nil {
			this.Logger.Debug2("Waiting for tasks to finish")
		}

	}
	wg.Wait()
	if this.Logger != nil {
		this.Logger.Debug2("All tasks finished")
	}

	return err
}

// Debug returns whether the application has the debug flag set
func (this *AppInstance) Debug() bool {
	return this.debug
}

// Verbose returns whether the application has the verbose flag set
func (this *AppInstance) Verbose() bool {
	return this.verbose
}

// Close method for app
func (this *AppInstance) Close() error {
	this.Logger.Debug2("gopi.AppInstance.Close()")
	if this.Layout != nil {
		if err := this.Layout.Close(); err != nil {
			return err
		}
		this.Layout = nil
	}
	if this.Display != nil {
		if err := this.Display.Close(); err != nil {
			return err
		}
		this.Display = nil
	}
	if this.Hardware != nil {
		if err := this.Hardware.Close(); err != nil {
			return err
		}
		this.Hardware = nil
	}
	if this.Logger != nil {
		if err := this.Logger.Close(); err != nil {
			return err
		}
		this.Logger = nil
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *AppInstance) setModuleInstance(t ModuleType, driver Driver) error {
	var ok bool
	switch t {
	case MODULE_TYPE_LOGGER:
		if this.Logger, ok = driver.(Logger); !ok {
			return fmt.Errorf("Module of type %v cannot be cast to gopi.Logger", t)
		}
	case MODULE_TYPE_HARDWARE:
		if this.Hardware, ok = driver.(HardwareDriver2); !ok {
			return fmt.Errorf("Module of type %v cannot be cast to gopi.Hardware", t)
		}
	case MODULE_TYPE_DISPLAY:
		if this.Display, ok = driver.(DisplayDriver2); !ok {
			return fmt.Errorf("Module of type %v cannot be cast to gopi.Display", t)
		}
	case MODULE_TYPE_LAYOUT:
		if this.Layout, ok = driver.(Layout); !ok {
			return fmt.Errorf("Module of type %v cannot be cast to gopi.Layout", t)
		}
	default:
		return fmt.Errorf("Not implemented: setModuleInstance: %v", t)
	}
	// success
	return nil
}

func getTestlessArguments(input []string) []string {
	output := make([]string, 0, len(input))
	for _, arg := range input {
		if strings.HasPrefix(arg, "-test.") {
			continue
		}
		output = append(output, arg)
	}
	return output
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *AppInstance) String() string {
	return fmt.Sprintf("gopi.App{ debug=%v verbose=%v }", this.debug, this.verbose)
}
