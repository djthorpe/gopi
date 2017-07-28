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
	"path"
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

	// Create the flags
	config.Flags = util.NewFlags(path.Base(os.Args[0]))

	// For each module, now call the configuration if the module
	// is registered
	for _, t := range config.Modules {
		if module, err := ModuleByType(t); err == nil {
			module.Config(&config)
		}
	}

	// Return the configuration
	return config
}

// NewAppInstance method will create a new application object given an application
// configuration
func NewAppInstance(config AppConfig) (*AppInstance, error) {
	this := new(AppInstance)

	// Create subsystems
	for _, t := range config.Modules {
		if module, err := ModuleByType(t); err == nil {
			module.Config(&config)
		}
		fmt.Println("TODO: Creating subsystem", t)
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
