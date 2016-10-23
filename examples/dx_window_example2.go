/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"flag"
	"fmt"
	"os"
)

import (
	"../" /* import "github.com/djthorpe/gopi" */
	"../util" /* import "github.com/djthorpe/gopi/util" */
	"../khronos" /* import "github.com/djthorpe/gopi/khronos" */
	"../device/rpi" /* import "github.com/djthorpe/gopi/device/rpi" */
)

////////////////////////////////////////////////////////////////////////////////

// The application
type App struct {
	// The logger
	Logger *util.LoggerDevice

	// The hardware device
	Device gopi.HardwareDriver

	// The opened display
	Display gopi.DisplayDriver

	// The EGL driver
	EGL khronos.EGLDriver
}

// Application configuration
type AppConfig struct {
	// The display number to open
	Display uint16

	// The file to log information to
	LogFile string

	// The level of logging
	LogLevel util.LogLevel

	// Whether to append to the log file
	LogAppend bool
}

// Run callback
type AppCallback func(*App) error

////////////////////////////////////////////////////////////////////////////////

var (
	flagDisplay = flag.Uint("display", 0,"Display number")
	flagVerbose = flag.Bool("verbose",false,"Output verbose logging messages")
	flagLogFile = flag.String("log","","Logging file. If empty, logs to stderr")
)

////////////////////////////////////////////////////////////////////////////////

func NewApp(config AppConfig) (*App, error) {
	var err error

	// Create application
	this := new(App)

	// Create a logger
	if len(config.LogFile) != 0 {
		this.Logger, err = util.Logger(util.FileLogger{ Filename: config.LogFile, Append: config.LogAppend })
	} else {
		this.Logger, err = util.Logger(util.StderrLogger{ })
	}
	if err != nil {
		return nil, err
	}

	// Set logging level
	this.Logger.SetLevel(config.LogLevel)

	// Create the device
	device, err := gopi.Open(rpi.Device{ },this.Logger)
	if err != nil {
		this.Logger.Close()
		return nil, err
	}
	// Convert device into a HardwareDriver
	this.Device = device.(gopi.HardwareDriver)

	// Open the display
	display, err := gopi.Open(rpi.DXDisplayConfig{
		Device: this.Device,
		Display: config.Display,
	},this.Logger)
	if err != nil {
		this.Device.Close()
		this.Logger.Close()
		return nil, err
	}
	// Convert device into a DisplayDriver
	this.Display = display.(gopi.DisplayDriver)

	// Create the EGL interface
	egl, err := gopi.Open(rpi.EGL{ Display: this.Display },this.Logger)
	if err != nil {
		this.Display.Close()
		this.Device.Close()
		this.Logger.Close()
		return nil, err
	}
	// Convert device into a EGLDriver
	this.EGL = egl.(khronos.EGLDriver)

	this.Logger.Debug("<App>Open")

	// success
	return this, nil
}

// Close the application
func (this *App) Close() error {
	this.Logger.Debug("<App>Close")

	if this.EGL != nil {
		if err := this.EGL.Close(); err != nil {
			return err
		}
	}
	if this.Display != nil {
		if err := this.Display.Close(); err != nil {
			return err
		}
	}
	if this.Device != nil {
		if err := this.Device.Close(); err != nil {
			return err
		}
	}
	if this.Logger != nil {
		if err := this.Logger.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Run the application with callback
func (this *App) Run(callback AppCallback) error {
	this.Logger.Debug("<App>Run")
	return callback(this)
}

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *App) error {
	app.Logger.Debug("Device=%v",app.Device)
	app.Logger.Debug("Display=%v",app.Display)
	app.Logger.Debug("EGL=%v",app.EGL)

	// Create a background
	bg, err := app.EGL.CreateBackground("OpenVG")
	if err != nil {
		return app.Logger.Error("Error: %v",err)
	}
	defer app.EGL.CloseWindow(bg)

	app.Logger.Debug("Background=%v",bg)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Parse flags
	flag.Parse()

	// Determine level of logging
	var level util.LogLevel
	if(*flagVerbose) {
		level = util.LOG_ANY
	} else {
		level = util.LOG_INFO
	}

	// Create the application
	app,err := NewApp(AppConfig{ Display: uint16(*flagDisplay), LogFile: *flagLogFile, LogAppend: false, LogLevel: level })
	if err != nil {
		fmt.Fprintln(os.Stderr,"Error:",err)
		return
	}
	defer app.Close()

	// Run the application
	if err := app.Run(RunLoop); err != nil {
		fmt.Fprintln(os.Stderr,"Error:",err)
		return
	}
}
