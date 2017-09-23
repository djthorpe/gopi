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
	"os/signal"
	"path"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// AppConfig defines how an application should be created
type AppConfig struct {
	Modules  []*Module
	AppArgs  []string
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
	debug    bool
	verbose  bool
	sigchan  chan os.Signal
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
// modules which should be created, the arguments are either by type
// or by name
func NewAppConfig(modules ...string) AppConfig {
	var err error

	config := AppConfig{}

	// retrieve modules and dependencies, using appendModule
	if config.Modules, err = ModuleWithDependencies("logger"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return AppConfig{}
	}
	for _, module := range modules {
		if module_array, err := ModuleWithDependencies(module); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return AppConfig{}
		} else {
			config.Modules = append(config.Modules, module_array...)
		}
	}

	// Set the flags
	config.AppArgs = getTestlessArguments(os.Args[1:])
	config.AppFlags = NewFlags(path.Base(os.Args[0]))
	config.Debug = false
	config.Verbose = false

	// Set 'debug' and 'verbose' flags
	config.AppFlags.FlagBool("debug", false, "Set debugging mode")
	config.AppFlags.FlagBool("verbose", false, "Verbose logging")

	// Call module.Config for each module
	for _, module := range config.Modules {
		if module.Config != nil {
			module.Config(&config)
		}
	}

	// Return the configuration
	return config
}

// NewAppInstance method will create a new application object given an application
// configuration
func NewAppInstance(config AppConfig) (*AppInstance, error) {

	if config.AppFlags == nil {
		return nil, ErrAppError
	}

	// Parse flags. We want to ignore flags which start with "-test."
	// in the testing environment
	if config.AppFlags != nil && config.AppFlags.Parsed() == false {
		if err := config.AppFlags.Parse(config.AppArgs); err != nil {
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

	// Set up signalling
	this.sigchan = make(chan os.Signal, 1)
	signal.Notify(this.sigchan, syscall.SIGTERM, syscall.SIGINT)

	// Create module instances
	var once sync.Once
	for _, module := range config.Modules {
		fmt.Println(module)
		// Report open (once after logger module is created)
		if this.Logger != nil {
			once.Do(func() {
				this.Logger.Debug2("gopi.AppInstance.Open()")
			})
		}
		if this.Logger != nil {
			this.Logger.Debug2("module.New{ %v }", module)
		}
		if driver, err := module.New(this); err != nil {
			return nil, err
		} else if err := this.setModuleInstance(module, driver); err != nil {
			return nil, err
		}
	}

	// report Open() again if it's not been done yet
	once.Do(func() {
		this.Logger.Debug2("gopi.AppInstance.Open()")
	})

	// success
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

// WaitForSignal blocks until a signal is caught
func (this *AppInstance) WaitForSignal() {
	signal := <-this.sigchan
	this.Logger.Debug("gopi.AppInstance.WaitForSignal: %v", signal)
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

func (this *AppInstance) setModuleInstance(module *Module, driver Driver) error {
	var ok bool
	switch module.Type {
	case MODULE_TYPE_LOGGER:
		if this.Logger, ok = driver.(Logger); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Logger", module)
		}
	case MODULE_TYPE_HARDWARE:
		if this.Hardware, ok = driver.(HardwareDriver2); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Hardware", module)
		}
	case MODULE_TYPE_DISPLAY:
		if this.Display, ok = driver.(DisplayDriver2); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Display", module)
		}
	case MODULE_TYPE_LAYOUT:
		if this.Layout, ok = driver.(Layout); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Layout", module)
		}
	default:
		return fmt.Errorf("Not implemented: setModuleInstance: %v", module)
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
