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
	"sync"

	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// AppConfig defines how an application should be created
type AppConfig struct {
	LogLevel util.LogLevel
	Modules  []ModuleType
	Flags    *util.Flags
	Debug    bool
}

// AppInstance defines the running application instance with modules
type AppInstance struct {
	Logger   Logger
	Hardware Driver
	Display  Driver
	Bitmap   Driver
	Vector   Driver
	VGFont   Driver
	OpenGL   Driver
	GPIO     Driver
	I2C      Driver
	SPI      Driver
	Input    Driver

	debug bool
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
	config.Flags = flags
	config.Debug = false

	// Return the configuration
	return config
}

// NewAppInstance method will create a new application object given an application
// configuration
func NewAppInstance(config AppConfig) (*AppInstance, error) {

	// Parse flags
	if config.Flags != nil && config.Flags.Parsed() == false {
		if err := config.Flags.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	// Set debug flag
	if debug, exists := config.Flags.GetBool("debug"); exists {
		config.Debug = debug
	}

	// Create instance
	this := new(AppInstance)
	this.debug = config.Debug

	// Create subsystems
	for _, t := range config.Modules {
		if module, err := ModuleByType(t); err != nil {
			return nil, err
		} else if driver, err := module.New(&config); err != nil {
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *AppInstance) setModuleInstance(t ModuleType, driver Driver) error {
	var ok bool
	switch t {
	case MODULE_TYPE_LOGGER:
		if this.Logger, ok = driver.(Logger); ok != true {
			return fmt.Errorf("Module of type %v cannot be cast to gopi.Logger", t)
		}
	case MODULE_TYPE_HARDWARE:
		this.Hardware = driver
	default:
		return fmt.Errorf("Not implmenented: setModuleInstance: %v", t)
	}
	// success
	return nil
}
