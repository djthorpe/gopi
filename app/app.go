/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package app /* import "github.com/djthorpe/gopi/app" */

import (
	gopi "../"           /* import "github.com/djthorpe/gopi" */
	rpi "../device/rpi"  /* import "github.com/djthorpe/gopi/util" */
	khronos "../khronos" /* import "github.com/djthorpe/gopi/khronos" */
	util "../util"       /* import "github.com/djthorpe/gopi/device/rpi" */
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

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

	// The OpenVG driver
	OpenVG khronos.VGDriver

	// The GPIO driver
	GPIO gopi.GPIODriver
}

// Application configuration
type AppConfig struct {
	// Application features
	Features AppFlags

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

// Application flags
type AppFlags uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Constants used to determine what features are needed
	APP_DEVICE     AppFlags = 0x0001
	APP_DISPLAY    AppFlags = 0x0002
	APP_EGL        AppFlags = 0x0004
	APP_OPENVG     AppFlags = 0x0008
	APP_GPIO       AppFlags = 0x0010
	APP_I2C        AppFlags = 0x0020
	APP_OPENGL_ES  AppFlags = 0x0040
)

////////////////////////////////////////////////////////////////////////////////
// Public Functions

func NewApp(config AppConfig) (*App, error) {
	var err error

	// Create application
	this := new(App)

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
	this.Logger.Debug("<App>Open device=%v display=%v egl=%v openvg=%v opengl_es=%v gpio=%v i2c=%v",
		config.Features&(APP_DEVICE|APP_DISPLAY|APP_EGL|APP_OPENVG|APP_GPIO|APP_I2C) != 0,
		config.Features&(APP_DISPLAY|APP_EGL|APP_OPENVG) != 0,
		config.Features&(APP_EGL|APP_OPENVG|APP_OPENGL_ES) != 0,
		config.Features&(APP_OPENVG) != 0,
		config.Features&(APP_OPENGL_ES) != 0,
		config.Features&(APP_GPIO) != 0,
		config.Features&(APP_I2C) != 0,
	)

	// Create the device
	if config.Features&(APP_DEVICE|APP_DISPLAY|APP_EGL|APP_OPENVG|APP_GPIO|APP_I2C) != 0 {
		device, err := gopi.Open(rpi.Device{}, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a HardwareDriver
		this.Device = device.(gopi.HardwareDriver)
	}

	// Create the display
	if config.Features&(APP_DISPLAY|APP_EGL|APP_OPENVG) != 0 {
		display, err := gopi.Open(rpi.DXDisplayConfig{
			Device:  this.Device,
			Display: config.Display,
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
	if config.Features & (APP_OPENVG) != 0 {
		openvg, err := gopi.Open(rpi.OpenVG{ EGL: this.EGL }, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a VGDriver
		this.OpenVG = openvg.(khronos.VGDriver)
	}

	// Create the GPIO interface
	if config.Features & (APP_GPIO) != 0 {
		openvg, err := gopi.Open(rpi.GPIO{ Device: this.Device }, this.Logger)
		if err != nil {
			this.Close()
			return nil, err
		}
		// Convert device into a GPIODriver
		this.GPIO = openvg.(gopi.GPIODriver)
	}

	// success
	return this, nil
}

// Close the application
func (this *App) Close() error {
	this.Logger.Debug("<App>Close")

	if this.GPIO != nil {
		if err := this.GPIO.Close(); err != nil {
			return err
		}
		this.GPIO = nil
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
	this.Logger.Debug("<App>Run")
	if err := callback(this); err != nil {
		return this.Logger.Error("%v",err)
	}
	return nil
}
