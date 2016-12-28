/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// APP
//
// This package provides the main application framework for developing
// applications for the GOPI library. In order to create an application,
// first create a configuration object, then create an application object
// with your configuration. For example,
//
//    config := app.Config(flags)
//    /* ... add additional configuration code here ... */
//    app, err := app.NewApp(config)
//    if err != nil { /* handle application error */ }
//    defer app.Close()
//
// The application will create all the necessary subsystems you need for your
// program, including hardware, display, GPU, GPIO, I2C and so forth. You
// specify what features you need through the 'flags' when calling the app.Config
// method.
//
package app /* import "github.com/djthorpe/gopi/app" */

// import
import (
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

// import abstract drivers
import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	khronos "github.com/djthorpe/gopi/khronos"
	util "github.com/djthorpe/gopi/util"
)

// import rpi drivers
import (
	rpi "github.com/djthorpe/gopi/device/rpi"
)

// import linux drivers
import (
	linux "github.com/djthorpe/gopi/device/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// The application
type App struct {
	// The logger
	Logger *util.LoggerDevice

	// Command-line flags
	FlagSet *Flags

	// The hardware device
	Device gopi.HardwareDriver

	// The opened display
	Display gopi.DisplayDriver

	// The EGL driver
	EGL khronos.EGLDriver

	// The OpenVG driver
	OpenVG khronos.VGDriver

	// The Font driver
	Fonts khronos.VGFontDriver

	// The GPIO driver
	GPIO hw.GPIODriver

	// The I2C driver
	I2C hw.I2CDriver

	// The Input driver
	Input hw.InputDriver

	// Signal channel on catching signals
	signal_channel chan os.Signal

	// Signal to place a finish bool on to indicate application should end
	finish_channel chan bool

	// debug and verbose flags
	debug, verbose bool
}

// Application configuration
type AppConfig struct {
	// Command-line flags
	FlagSet *Flags

	// Application features
	Features AppFlags

	// The display number to open
	Display uint16

	// The I2C bus
	I2CBus uint

	// The file to log information to
	LogFile string

	// The level of logging
	LogLevel util.LogLevel

	// Whether to append to the log file
	LogAppend bool

	// An array of font paths from which to load fonts
	FontPaths []string
}

// Run callback
type AppCallback func(*App) error

// Application flags
type AppFlags uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Constants used to determine what subsystems are needed
	APP_DEVICE    AppFlags = 0x0001
	APP_DISPLAY   AppFlags = 0x0002
	APP_EGL       AppFlags = 0x0004
	APP_OPENVG    AppFlags = 0x0008
	APP_GPIO      AppFlags = 0x0010
	APP_I2C       AppFlags = 0x0020
	APP_OPENGL_ES AppFlags = 0x0040
	APP_VGFONT    AppFlags = 0x0080
	APP_INPUT     AppFlags = 0x0100
)

const (
	// Constants used to determine default flags
	APP_DEFAULT_DEBUG   bool   = false
	APP_DEFAULT_VERBOSE bool   = false
	APP_DEFAULT_DISPLAY uint16 = 0
)

////////////////////////////////////////////////////////////////////////////////
// Public Functions

// Create a configuration based on the features you'd like in your application.
// This will return a pre-filled AppConfig object, with appropriate command-line
// flags also in place.
func Config(flags AppFlags) AppConfig {
	config := AppConfig{}

	// create flagset and set appflags
	config.FlagSet = NewFlags(path.Base(os.Args[0]))
	config.Features = flags

	// create nil value for FontPaths
	config.FontPaths = nil

	// Add on -log flag for path to logfile
	config.FlagSet.FlagString("log", "", "File for logging")
	config.FlagSet.FlagBool("verbose", APP_DEFAULT_VERBOSE, "Log verbosely")
	config.FlagSet.FlagBool("debug", APP_DEFAULT_DEBUG, "Trigger debugging support")

	// Add -display
	if config.Features&(APP_DISPLAY|APP_EGL|APP_OPENVG|APP_OPENGL_ES) != 0 {
		config.Display = APP_DEFAULT_DISPLAY
		config.FlagSet.FlagUint("display", uint(config.Display), "Display to use")
	}

	// Add -ppi
	if config.Features&(APP_DISPLAY|APP_EGL|APP_OPENVG|APP_OPENGL_ES|APP_VGFONT) != 0 {
		config.FlagSet.FlagString("displaysize", "", "Display physical size (ie, 99in 99mm 99cm 99x99in 99x99mm 99x99cm)")
	}

	// Add -i2cbus
	if config.Features&(APP_I2C) != 0 {
		config.FlagSet.FlagUint("i2cbus", 1, "I2C Bus")
	}

	// Add -fontpath
	if config.Features&(APP_VGFONT) != 0 {
		config.FlagSet.FlagString("fontpath", "", "Path to font file or folder, will recurse into subfolders (file extensions allowed: .TTF, .TTC, .OTF, .OTC, .DFONT)")
	}

	return config
}

// Create a new application object. This will return the application object
// or an error. If the error is flag.ErrHelp then the usage information for
// the application is printed on stderr, and you should simply quit the
// application. Other errors might occur depending what features have been
// requests for the application.
func NewApp(config AppConfig) (*App, error) {
	var err error

	// Parse command-line flags
	if config.FlagSet != nil {
		// Parse command-line flags
		if err := config.FlagSet.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	// Create application
	this := new(App)

	// Set FlagSet
	this.FlagSet = config.FlagSet

	// Set Debug and Verbose
	this.debug = this.getDebug()
	this.verbose = this.getVerbose()

	// Set display and diagnol size in inches
	display, screensize, err := this.getDisplay(config.Display)
	if err != nil {
		return nil, err
	}

	// Set log level
	if config.LogLevel == util.LOG_ANY {
		if this.debug && this.verbose {
			config.LogLevel = util.LOG_DEBUG2
		} else if this.debug {
			config.LogLevel = util.LOG_DEBUG
		} else if this.verbose {
			config.LogLevel = util.LOG_INFO
		} else {
			config.LogLevel = util.LOG_WARN
		}
	}

	// Create a logger - either log to file or to stderr
	if len(config.LogFile) != 0 {
		this.Logger, err = util.Logger(util.FileLogger{Filename: config.LogFile, Append: config.LogAppend})
	} else {
		this.Logger, err = util.Logger(util.StderrLogger{})
	}
	if err != nil {
		return nil, err
	}

	// Set logging level
	this.Logger.SetLevel(config.LogLevel)

	// Debugging
	this.Logger.Debug("<App>Open device=%v display=%v egl=%v openvg=%v opengl_es=%v vgfont=%v input=%v gpio=%v i2c=%v",
		config.Features&(APP_DEVICE|APP_DISPLAY|APP_EGL|APP_OPENVG|APP_VGFONT|APP_GPIO|APP_I2C) != 0,
		config.Features&(APP_DISPLAY|APP_EGL|APP_OPENVG|APP_VGFONT) != 0,
		config.Features&(APP_EGL|APP_OPENVG|APP_OPENGL_ES) != 0,
		config.Features&(APP_OPENVG) != 0,
		config.Features&(APP_OPENGL_ES) != 0,
		config.Features&(APP_VGFONT) != 0,
		config.Features&(APP_INPUT) != 0,
		config.Features&(APP_GPIO) != 0,
		config.Features&(APP_I2C) != 0,
	)

	// Signal handlers
	this.signal_channel = make(chan os.Signal, 1)
	this.finish_channel = make(chan bool, 1)
	signal.Notify(this.signal_channel, syscall.SIGTERM, syscall.SIGINT)

	// Create the device
	if config.Features&(APP_DEVICE|APP_DISPLAY|APP_EGL|APP_OPENVG|APP_VGFONT|APP_GPIO|APP_I2C) != 0 {
		device, err := gopi.Open(rpi.Hardware{}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a HardwareDriver
		this.Device = device.(gopi.HardwareDriver)
	}

	// Create the display
	if config.Features&(APP_DISPLAY|APP_EGL|APP_OPENVG|APP_VGFONT) != 0 {
		display, err := gopi.Open(rpi.DXDisplayConfig{
			Device:         this.Device,
			Display:        display,
			PhysicalInches: screensize,
		}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a DisplayDriver
		this.Display = display.(gopi.DisplayDriver)
	}

	// Create the EGL interface
	if config.Features&(APP_EGL|APP_OPENVG) != 0 {
		egl, err := gopi.Open(rpi.EGL{Display: this.Display}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a EGLDriver
		this.EGL = egl.(khronos.EGLDriver)
	}

	// Create the OpenVG interface
	if config.Features&(APP_OPENVG) != 0 {
		openvg, err := gopi.Open(rpi.OpenVG{EGL: this.EGL}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a VGDriver
		this.OpenVG = openvg.(khronos.VGDriver)
	}

	// Create the Font driver subsystem
	if config.Features&(APP_VGFONT) != 0 {
		ppi := uint(0)
		if this.Display != nil {
			ppi = uint(this.Display.GetPixelsPerInch())
		}
		fontdriver, err := gopi.Open(rpi.VGFont{PPI: ppi}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert fontdriver into a VGFontDriver
		this.Fonts = fontdriver.(khronos.VGFontDriver)

		// Load fonts
		if config.FontPaths != nil {
			err = this.loadFonts(config.FontPaths)
		} else {
			if fontpaths, err := this.getFontPaths(); err == nil {
				err = this.loadFonts(fontpaths)
			}
		}
		if err != nil {
			this.Close()
			return nil, err
		}
	}

	// Create the Input subsystem
	if config.Features&(APP_INPUT) != 0 {
		input, err := gopi.Open(linux.Input{}, this.Logger)
		if err != nil {
			this.Logger.Info("Error %v",err)
			this.Close()
			return nil, err
		}
		// Convert device into an InputDriver
		this.Input = input.(hw.InputDriver)
	}

	// Create the GPIO interface
	if config.Features&(APP_GPIO) != 0 {
		gpio, err := gopi.Open(rpi.GPIO{Device: this.Device}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a GPIODriver
		this.GPIO = gpio.(hw.GPIODriver)
	}

	// Create the I2C interface
	if config.Features&(APP_I2C) != 0 {
		bus, exists := this.FlagSet.GetUint("i2cbus")
		if exists {
			config.I2CBus = bus
		}
		i2c, err := gopi.Open(linux.I2C{Bus: config.I2CBus}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a GPIODriver
		this.I2C = i2c.(hw.I2CDriver)
	}

	// success
	return this, nil
}

// Close the application. This will free any resources opened. It will return
// an error on unsuccessful, or nil otherwise.
func (this *App) Close() error {
	this.Logger.Debug2("<App>Close")

	if this.I2C != nil {
		if err := this.I2C.Close(); err != nil {
			return err
		}
		this.I2C = nil
	}
	if this.Input != nil {
		if err := this.Input.Close(); err != nil {
			return err
		}
		this.Input = nil
	}
	if this.GPIO != nil {
		if err := this.GPIO.Close(); err != nil {
			return err
		}
		this.GPIO = nil
	}
	if this.Fonts != nil {
		if err := this.Fonts.Close(); err != nil {
			return err
		}
		this.Fonts = nil
	}
	if this.OpenVG != nil {
		if err := this.OpenVG.Close(); err != nil {
			return err
		}
		this.OpenVG = nil
	}
	if this.EGL != nil {
		if err := this.EGL.Close(); err != nil {
			return err
		}
		this.EGL = nil
	}
	if this.Display != nil {
		if err := this.Display.Close(); err != nil {
			return err
		}
		this.Display = nil
	}
	if this.Device != nil {
		if err := this.Device.Close(); err != nil {
			return err
		}
		this.Device = nil
	}
	if this.Logger != nil {
		if err := this.Logger.Close(); err != nil {
			return err
		}
		this.Logger = nil
	}
	return nil
}

// Run the application with callback
func (this *App) Run(callback AppCallback) error {
	this.Logger.Debug2("<App>Run pid=%v", os.Getpid())

	// Load the fonts
	//if err := this.loadfonts(); err != nil {
	//	return this.Logger.Error("Error loading fonts: %v",err)
	//}

	// Go routine to wait for signal, and send finish signal in that case
	go func() {
		signal := <-this.signal_channel
		this.Logger.Debug("<App>Run: caught signal: %v", signal)
		this.finish_channel <- true
	}()

	if err := callback(this); err != nil {
		return this.Logger.Error("%v", err)
	}

	// TODO: Remove signal handlers

	return nil
}

// Wait until the finish channel has an event on it
func (this *App) WaitUntilDone() {
	this.Logger.Debug2("<App>WaitUntilDone")

	// Runloop accepting events, until done
	done := false
	for done == false {
		select {
		case done = <-this.finish_channel:
			done = true
			break
		}
	}
}

// Set the Done signal which terminates the WaitUntilDone loop
func (this *App) Done() {
	this.Logger.Debug2("<App>Done")
	this.finish_channel <- true
}

// Return the debug flag
func (this *App) GetDebug() bool {
	return this.debug
}

// Return the verbose flag
func (this *App) GetVerbose() bool {
	return this.verbose
}

////////////////////////////////////////////////////////////////////////////////
// Private Methods

// Return the debug flag
func (this *App) getDebug() bool {
	if this.FlagSet == nil {
		return APP_DEFAULT_DEBUG
	}
	debug, exists := this.FlagSet.GetBool("debug")
	if exists == false {
		return APP_DEFAULT_DEBUG
	}
	return debug
}

// Return the verbose flag
func (this *App) getVerbose() bool {
	if this.FlagSet == nil {
		return APP_DEFAULT_VERBOSE
	}
	verbose, exists := this.FlagSet.GetBool("verbose")
	if exists == false {
		return APP_DEFAULT_VERBOSE
	}
	return verbose
}

// Return screen size
func (this *App) getDisplay(default_display uint16) (uint16, float64, error) {
	if this.FlagSet == nil {
		return default_display, 0.0, nil
	}
	display, exists := this.FlagSet.GetUint16("display")
	if exists == false {
		display = default_display
	}
	screensize, exists := this.FlagSet.GetString("displaysize")
	if exists == false {
		return display, 0.0, nil
	}
	inches, err := util.ParseLengthString(screensize)
	if err != nil {
		return display, inches, this.Logger.Error("Invalid -displaysize: %v", err)
	}
	return display, inches, nil
}

// Get font path from command line
func (this *App) getFontPaths() ([]string, error) {
	if this.FlagSet == nil {
		return nil, nil
	}
	path, exists := this.FlagSet.GetString("fontpath")
	if exists == false {
		return nil, nil
	}
	return []string{path}, nil
}

// Load fonts
func (this *App) loadFonts(paths []string) error {
	// Ignore if fonts are not yet loaded
	if this.Fonts == nil {
		return nil
	}
	// Ignore if there are no font paths
	if paths == nil || len(paths) == 0 {
		return nil
	}
	// Now iterate through font paths
	for _, filepath := range paths {
		err := this.Fonts.OpenFacesAtPath(filepath, func(filename string, info os.FileInfo) bool {
			if strings.HasPrefix(info.Name(), ".") {
				// ignore hidden files and folders
				return false
			}
			if info.IsDir() {
				// recurse into folders
				return true
			}
			if path.Ext(filename) == ".ttf" || path.Ext(filename) == ".TTF" {
				// support TTF loading
				return true
			}
			if path.Ext(filename) == ".ttc" || path.Ext(filename) == ".TTC" {
				// support TTC loading
				return true
			}
			if path.Ext(filename) == ".otf" || path.Ext(filename) == ".OTF" {
				// support OTF loading
				return true
			}
			if path.Ext(filename) == ".otc" || path.Ext(filename) == ".OTC" {
				// support OTC loading
				return true
			}
			if path.Ext(filename) == ".dfont" {
				// support DFONT loading
				return true
			}
			if path.Ext(filename) == ".txt" || path.Ext(filename) == ".TXT" || path.Ext(filename) == "" {
				// silently ignore txt files and files without extension
				return false
			}
			this.Logger.Warn("Whilst loading fonts at path %v: Ignoring file %v", filepath, filename)
			return false
		})
		if err != nil {
			return this.Logger.Error("Error loading fonts at path: %v: %v", filepath, err)
		}
	}

	return nil
}
