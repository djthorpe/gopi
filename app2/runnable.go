package app2 /* import "github.com/djthorpe/gopi/app2" */

import (
	"errors"
	"runtime"
	"sync"

	gopi "github.com/djthorpe/gopi"

	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES AND INTERFACES

// AppConfig defines how an application should be created
type AppConfig struct {
	LogLevel util.LogLevel
	Modules  []ModuleType
}

// AppInstance defines the running application instance with modules
type AppInstance struct {
	Logger  gopi.Logger
	Device  gopi.Driver
	Display gopi.Driver
	EGL     gopi.Driver
	OpenVG  gopi.Driver
	VGFont  gopi.Driver
	OpenGL  gopi.Driver
	GPIO    gopi.Driver
	I2C     gopi.Driver
	SPI     gopi.Driver
	Input   gopi.Driver
}

// Task defines a function which can run, and has a channel which
// indicates when the main thread has finished
type Task func(app *AppInstance, done chan struct{}) error

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// ErrNoTasks is an error returned when no tasks are to be run
	ErrNoTasks = errors.New("No tasks to run")
	// ErrAppError is a general application error
	ErrAppError = errors.New("General application error")
	// ErrModuleNotFound is an error when module cannot be found by name or type
	ErrModuleNotFound = errors.New("Module not found")
)

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
	config.Modules = modules
	return config
}

// NewAppInstance method will create a new application object given an application
// configuration
func NewAppInstance(config AppConfig) (*AppInstance, error) {
	this := new(AppInstance)

	// Create Logger object
	if logger, err := util.Logger(util.StderrLogger{}); err != nil {
		return nil, err
	} else {
		logger.SetLevel(config.LogLevel)
		this.Logger = logger
	}
	if this.Logger == nil {
		return nil, ErrAppError
	}

	// Create subsystems
	for _, name := range config.Modules {
		this.Logger.Debug("Creating subsystem %v", name)

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
		this.Logger.Debug2("Main thread done")
		// Signal other tasks to complete
		for i := 1; i < len(tasks); i++ {
			channels[i] <- DONE
		}
	}()

	// Now run main task
	err := tasks[0](this, channels[0])

	// Wait for other tasks to finish
	if len(tasks) > 1 {
		this.Logger.Debug2("Waiting for tasks to finish")
	}
	wg.Wait()
	this.Logger.Debug2("All tasks finished")

	return err
}
