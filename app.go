/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/util/errors"
	"github.com/djthorpe/gopi/util/tasks"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// AppParam is a list of application parameters
type AppParam uint

// AppConfig defines how an application should be created
type AppConfig struct {
	Modules  []*Module
	AppArgs  []string
	AppFlags *Flags
	Params   map[AppParam]interface{}
	Debug    bool
	Verbose  bool
}

// AppInstance defines the running application instance with modules
type AppInstance struct {
	AppFlags   *Flags
	Logger     Logger
	Hardware   Hardware
	Display    Display
	Graphics   SurfaceManager
	Sprites    SpriteManager
	Input      InputManager
	Fonts      FontManager
	Layout     Layout
	Timer      Timer
	GPIO       GPIO
	I2C        I2C
	SPI        SPI
	PWM        PWM
	LIRC       LIRC
	ClientPool RPCClientPool
	debug      bool
	verbose    bool
	sigchan    chan os.Signal
	modules    []*Module
	byname     map[string]Driver
	bytype     map[ModuleType]Driver
	byorder    []Driver

	// background tasks implementation
	tasks.Tasks
}

// MainTask defines a function which can run as a main task
// and has a channel which can be written to when the task
// has completed
type MainTask func(app *AppInstance, done chan<- struct{}) error

// BackgroundTask defines a function which can run as a
// background task and has a channel which receives a value of gopi.DONE
// then the background task should complete
type BackgroundTask func(app *AppInstance, done <-chan struct{}) error
type BackgroundTask2 func(app *AppInstance, start chan<- struct{}, stop <-chan struct{}) error

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	// PARAM_SERVICENAME_DEFAULT is the default service type
	PARAM_SERVICE_TYPE_DEFAULT = "gopi"
)

var (
	// DONE is the message sent on the channel to indicate task is completed
	DONE = struct{}{}
)

const (
	// Application paramaters are used to store global constant data
	// from startup, usually used for storing versions and other
	// information from time of building
	PARAM_NONE AppParam = iota
	PARAM_TIMESTAMP
	PARAM_EXECNAME
	PARAM_SERVICE_NAME
	PARAM_SERVICE_TYPE
	PARAM_SERVICE_SUBTYPE
	PARAM_GOVERSION
	PARAM_GOBUILDTIME
	PARAM_GITTAG
	PARAM_GITBRANCH
	PARAM_GITHASH
	PARAM_MAX = PARAM_GITHASH
	PARAM_MIN = PARAM_TIMESTAMP
)

var (
	// Build and version flags
	GitTag      string
	GitBranch   string
	GitHash     string
	GoBuildTime string
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

	// append other modules
	if config.Modules, err = AppendModulesByName(config.Modules, modules...); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return AppConfig{}
	}

	// Set the flags
	config.AppArgs = getTestlessArguments(os.Args[1:])
	config.AppFlags = NewFlags(path.Base(os.Args[0]))
	config.Debug = false
	config.Verbose = false

	// Set the parameters
	config.AppFlags.params[PARAM_SERVICE_TYPE] = PARAM_SERVICE_TYPE_DEFAULT
	config.AppFlags.params[PARAM_EXECNAME] = config.AppFlags.Name()
	config.AppFlags.params[PARAM_TIMESTAMP] = time.Now()
	config.AppFlags.params[PARAM_GOVERSION] = runtime.Version()
	config.AppFlags.params[PARAM_GITTAG] = GitTag
	config.AppFlags.params[PARAM_GITBRANCH] = GitBranch
	config.AppFlags.params[PARAM_GITHASH] = GitHash
	config.AppFlags.params[PARAM_GOBUILDTIME] = GoBuildTime

	// Set 'debug', 'verbose' and 'version' flags
	config.AppFlags.FlagBool("debug", false, "Set debugging mode")
	config.AppFlags.FlagBool("verbose", false, "Verbose logging")
	config.AppFlags.FlagBool("version", false, "Print version information and exit")

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
		// Check for version flag
		if version, _ := config.AppFlags.GetBool("version"); version {
			config.AppFlags.PrintVersion()
			return nil, ErrHelp
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

	// Set module maps
	this.modules = config.Modules
	this.byname = make(map[string]Driver, len(config.Modules))
	this.bytype = make(map[ModuleType]Driver, len(config.Modules))
	this.byorder = make([]Driver, 0, len(config.Modules))

	// Create module instances
	var once sync.Once
	for _, module := range config.Modules {
		// Report open (once after logger module is created)
		if this.Logger != nil {
			once.Do(func() {
				this.Logger.Debug("gopi.AppInstance.Open(){ modules=%v }", config.Modules)
			})
		}
		if module.New != nil {
			if this.Logger != nil {
				this.Logger.Debug2("module.New{ %v }", module)
			}
			if driver, err := module.New(this); err != nil {
				return nil, err
			} else if driver == nil {
				return nil, fmt.Errorf("%v: New: return nil", module.Name)
			} else if err := this.setModuleInstance(module, driver); err != nil {
				if err := driver.Close(); err != nil {
					this.Logger.Error("module.Close(): %v", err)
				}
				return nil, err
			}
		}
	}

	// report Open() again if it's not been done yet
	once.Do(func() {
		this.Logger.Debug("gopi.AppInstance.Open()")
	})

	// success
	return this, nil
}

// Run all tasks simultaneously, the first task in the list on the main thread and the
// remaining tasks background tasks.
func (this *AppInstance) Run(main_task MainTask, background_tasks ...BackgroundTask) error {
	// Lock this to run in the current operating system thread (ie, the main thread)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Call the Run method for each module. If any report an error, then don't run
	// the application. Note that some modules don't have a 'New' method in which
	// case the driver argument is set to nil
	for _, module := range this.modules {
		if module.Run == nil {
			continue
		}
		driver, _ := this.byname[module.Name]
		if err := module.Run(this, driver); err != nil {
			return err
		}
	}

	// create the channels we'll use to signal the goroutines
	channels := make([]chan struct{}, len(background_tasks)+1)
	channels[0] = make(chan struct{})
	for i := range background_tasks {
		channels[i+1] = make(chan struct{})
	}

	// if more than one task, then give them a channel which is signalled
	// by the main thread for ending
	var wg sync.WaitGroup
	if len(background_tasks) > 0 {
		for i, task := range background_tasks {
			wg.Add(1)
			go func(i int, t BackgroundTask) {
				defer wg.Done()
				if err := t(this, channels[i+1]); err != nil {
					if this.Logger != nil {
						this.Logger.Error("Error: %v [background_task %v]", err, i+1)
					}
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
		for i := 0; i < len(background_tasks); i++ {
			this.Logger.Debug2("Sending DONE to background task %v of %v", i+1, len(background_tasks))
			channels[i+1] <- DONE
		}
	}()

	// Now run main task
	err := main_task(this, channels[0])

	// Wait for other tasks to finish
	if len(background_tasks) > 0 {
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

// Run all tasks simultaneously, the first task in the list on the main thread and the
// remaining tasks background tasks. The main task doesn't start running until received
// start signals from the background tasks
func (this *AppInstance) Run2(main_task MainTask, background_tasks ...BackgroundTask2) error {
	// Lock this to run in the current operating system thread (ie, the main thread)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Call the Run method for each module. If any report an error, then don't run
	// the application. Note that some modules don't have a 'New' method in which
	// case the driver argument is set to nil
	for _, module := range this.modules {
		if module.Run == nil {
			continue
		}
		driver, _ := this.byname[module.Name]
		if err := module.Run(this, driver); err != nil {
			return err
		}
	}

	// Start the background tasks - wait for start signals
	if len(background_tasks) > 0 {
		t := make([]tasks.TaskFunc, len(background_tasks))
		for i := range background_tasks {
			f := background_tasks[i]
			t[i] = func(start chan<- struct{}, stop <-chan struct{}) error {
				return f(this, start, stop)
			}
		}
		this.Tasks.Start(t...)
	}

	// Run a background thread waiting for a done signal from main
	errs := make(chan error)
	done := make(chan struct{})
	go func() {
		// Wait for main done signal
		<-done
		if this.Logger != nil {
			this.Logger.Debug2("Main thread done, stopping background tasks")
		}
		// Signal other tasks to complete
		errs <- this.Tasks.Close()
		if this.Logger != nil {
			this.Logger.Debug2("Main thread done, Background tasks stopped")
		}
	}()

	// Now run main task
	return_error := new(errors.CompoundError)
	if err := main_task(this, done); err != nil {
		return_error.Add(err)
		if this.Logger != nil {
			this.Logger.Debug2("Error from main thread: %v", err)
		}
	}

	// Close the 'done' channel and wait for closing gorouting to finish
	close(done)

	// Receive any addition errors from tasks
	if tasks_err := <-errs; tasks_err != nil {
		return_error.Add(tasks_err)
	}

	// Indicate all tasks finished
	if this.Logger != nil {
		this.Logger.Debug2("All tasks finished")
	}

	// Return error from main
	return return_error.ErrorOrSelf()
}

// Debug returns whether the application has the debug flag set
func (this *AppInstance) Debug() bool {
	return this.debug
}

// Verbose returns whether the application has the verbose flag set
func (this *AppInstance) Verbose() bool {
	return this.verbose
}

// Service returns the current service name set from configuration
func (this *AppInstance) Service() (string, string, string, error) {
	var service, subtype, name string
	if service_, exists := this.AppFlags.params[PARAM_SERVICE_TYPE]; exists {
		service = fmt.Sprint(service_)
	} else {
		service = PARAM_SERVICE_TYPE_DEFAULT
	}
	if subtype_, exists := this.AppFlags.params[PARAM_SERVICE_SUBTYPE]; exists {
		subtype = fmt.Sprint(subtype_)
	}
	if name_, exists := this.AppFlags.params[PARAM_SERVICE_NAME]; exists {
		name = fmt.Sprint(name_)
	} else if hostname, err := os.Hostname(); err != nil {
		return "", "", "", err
	} else {
		name = hostname
	}
	return service, subtype, name, nil
}

// WaitForSignal blocks until a signal is caught
func (this *AppInstance) WaitForSignal() {
	s := <-this.sigchan
	this.Logger.Debug2("gopi.AppInstance.WaitForSignal: %v", s)
}

// WaitForSignalOrTimeout blocks until a signal is caught or
// timeout occurs and return true if the signal is caught
func (this *AppInstance) WaitForSignalOrTimeout(timeout time.Duration) bool {
	select {
	case s := <-this.sigchan:
		this.Logger.Debug2("gopi.AppInstance.WaitForSignalOrTimeout: %v", s)
		return true
	case <-time.After(timeout):
		return false
	}
}

// SendSignal will send the terminate signal, breaking the WaitForSignal
// block
func (this *AppInstance) SendSignal() error {
	if process, err := os.FindProcess(os.Getpid()); err != nil {
		return err
	} else if err := process.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	return nil
}

// Close method for app
func (this *AppInstance) Close() error {
	this.Logger.Debug("gopi.AppInstance.Close()")

	// In reverse order, call the Close method on each
	// driver
	for i := len(this.byorder); i > 0; i-- {
		driver := this.byorder[i-1]
		this.Logger.Debug2("gopi.AppInstance.Close() %v", driver)
		if err := driver.Close(); err != nil {
			this.Logger.Error("gopi.AppInstance.Close() error: %v", err)
		}
	}

	// Quit tasks if not already quit
	this.Tasks.Close()

	// Clear out the references
	this.bytype = nil
	this.byname = nil
	this.byorder = nil

	this.Layout = nil
	this.Display = nil
	this.Graphics = nil
	this.Sprites = nil
	this.Fonts = nil
	this.Hardware = nil
	this.Logger = nil
	this.Timer = nil
	this.Input = nil
	this.I2C = nil
	this.GPIO = nil
	this.SPI = nil
	this.PWM = nil
	this.LIRC = nil
	this.ClientPool = nil

	// Return success
	return nil
}

// ModuleInstance returns module instance by name, or returns nil if the module
// cannot be found. You can use reserved words (ie, logger, layout, etc)
// for common module types
func (this *AppInstance) ModuleInstance(name string) Driver {
	var instance Driver
	// Check for reserved words
	if module_type, exists := module_name_map[name]; exists {
		instance, _ = this.bytype[module_type]
	} else {
		instance, _ = this.byname[name]
	}
	return instance
}

// Append Modules by name onto the configuration
func AppendModulesByName(modules []*Module, names ...string) ([]*Module, error) {
	if modules == nil {
		modules = make([]*Module, 0, len(names))
	}
	for _, name := range names {
		if module_array, err := ModuleWithDependencies(name); err != nil {
			return nil, err
		} else {
			modules = appendModules(modules, module_array)
		}
	}
	return modules, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// appendModules adds modules from 'others' onto 'modules' without
// creating duplicate modules
func appendModules(modules []*Module, others []*Module) []*Module {
	if len(others) == 0 {
		return modules
	}
	for _, other := range others {
		if inModules(modules, other) {
			continue
		}
		modules = append(modules, other)
	}
	return modules
}

// inModules returns true if a module is in the array
func inModules(modules []*Module, other *Module) bool {
	for _, module := range modules {
		if module == other {
			return true
		}
	}
	return false
}

func (this *AppInstance) setModuleInstance(module *Module, driver Driver) error {
	var ok bool

	// Set by name. Currently returns an error if there is more than one module with the same name
	if _, exists := this.byname[module.Name]; exists {
		return fmt.Errorf("setModuleInstance: Duplicate module with name '%v'", module.Name)
	} else {
		this.byname[module.Name] = driver
	}

	// Set by type. Currently returns an error if there is more than one module with the same type
	// Allows multiple modules accessed by name if other, service or client
	if module.Type != MODULE_TYPE_NONE && module.Type != MODULE_TYPE_OTHER && module.Type != MODULE_TYPE_SERVICE && module.Type != MODULE_TYPE_CLIENT {
		if _, exists := this.bytype[module.Type]; exists {
			return fmt.Errorf("setModuleInstance: Duplicate module with type '%v'", module.Type)
		} else {
			this.bytype[module.Type] = driver
		}
	}

	// Append to list of modules in the order (so we can close in the right order
	// later)
	this.byorder = append(this.byorder, driver)

	// Now some convenience methods for already-cast drivers
	switch module.Type {
	case MODULE_TYPE_LOGGER:
		if this.Logger, ok = driver.(Logger); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Logger", module)
		}
	case MODULE_TYPE_HARDWARE:
		if this.Hardware, ok = driver.(Hardware); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Hardware", module)
		}
	case MODULE_TYPE_DISPLAY:
		if this.Display, ok = driver.(Display); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Display", module)
		}
	case MODULE_TYPE_GRAPHICS:
		if this.Graphics, ok = driver.(SurfaceManager); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.SurfaceManager", module)
		}
	case MODULE_TYPE_SPRITES:
		if this.Sprites, ok = driver.(SpriteManager); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.SpriteManager", module)
		}
	case MODULE_TYPE_FONTS:
		if this.Fonts, ok = driver.(FontManager); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.FontManager", module)
		}
	case MODULE_TYPE_LAYOUT:
		if this.Layout, ok = driver.(Layout); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Layout", module)
		}
	case MODULE_TYPE_GPIO:
		if this.GPIO, ok = driver.(GPIO); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.GPIO", module)
		}
	case MODULE_TYPE_I2C:
		if this.I2C, ok = driver.(I2C); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.I2C", module)
		}
	case MODULE_TYPE_SPI:
		if this.SPI, ok = driver.(SPI); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.SPI", module)
		}
	case MODULE_TYPE_PWM:
		if this.PWM, ok = driver.(PWM); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.PWM", module)
		}
	case MODULE_TYPE_TIMER:
		if this.Timer, ok = driver.(Timer); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.Timer", module)
		}
	case MODULE_TYPE_LIRC:
		if this.LIRC, ok = driver.(LIRC); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.LIRC", module)
		}
	case MODULE_TYPE_INPUT:
		if this.Input, ok = driver.(InputManager); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.InputManager", module)
		}
	case MODULE_TYPE_CLIENTPOOL:
		if this.ClientPool, ok = driver.(RPCClientPool); !ok {
			return fmt.Errorf("Module %v cannot be cast to gopi.RPCClientPool", module)
		}
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
	modules := make([]string, 0, len(this.byname))
	for k := range this.byname {
		modules = append(modules, k)
	}
	return fmt.Sprintf("gopi.App{ debug=%v verbose=%v modules=%v instances=%v }", this.debug, this.verbose, modules, this.byorder)
}

func (p AppParam) String() string {
	switch p {
	case PARAM_NONE:
		return "PARAM_NONE"
	case PARAM_TIMESTAMP:
		return "PARAM_TIMESTAMP"
	case PARAM_EXECNAME:
		return "PARAM_EXECNAME"
	case PARAM_SERVICE_NAME:
		return "PARAM_SERVICE_NAME"
	case PARAM_SERVICE_TYPE:
		return "PARAM_SERVICE_TYPE"
	case PARAM_SERVICE_SUBTYPE:
		return "PARAM_SERVICE_SUBTYPE"
	case PARAM_GOVERSION:
		return "PARAM_GOVERSION"
	case PARAM_GOBUILDTIME:
		return "PARAM_GOBUILDTIME"
	case PARAM_GITTAG:
		return "PARAM_GITTAG"
	case PARAM_GITBRANCH:
		return "PARAM_GITBRANCH"
	case PARAM_GITHASH:
		return "PARAM_GITHASH"
	default:
		return "[?? Invalid AppParam value]"
	}
}
