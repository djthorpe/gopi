package mock

/* sys/hw/mock */

import (
	"fmt"
	"strings"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Hardware struct{}

type Display struct {
	Display uint
}

type hardwareDriver struct{}

type displayDriver struct {
	id uint
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware
	gopi.RegisterModule(gopi.Module{
		Name: "hardware/mock",
		Type: gopi.MODULE_TYPE_HARDWARE,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Hardware{}, app.Logger)
		},
	})
	// Register display
	gopi.RegisterModule(gopi.Module{
		Name: "display/mock",
		Type: gopi.MODULE_TYPE_DISPLAY,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("display", 0, "Display")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			display := Display{}
			if display_number, exists := app.AppFlags.GetUint("display"); exists {
				display.Display = display_number
			}
			return gopi.Open(display, app.Logger)
		},
	})
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - HARDWARE

// Open
func (config Hardware) Open(logger gopi.Logger) (gopi.Driver, error) {
	return new(hardwareDriver), nil
}

// GetName returns the name of the hardware (ie, mock, mac, linux, rpi, etc)
func (this *hardwareDriver) Name() string {
	return "mock/hardware"
}

// GetName returns the name of the hardware (ie, mock, mac, linux, rpi, etc)
func (this *hardwareDriver) SerialNumber() string {
	return strings.ToUpper("SERIAL_NUMBER")
}

func (this *hardwareDriver) NumberOfDisplays() uint {
	return 0
}

func (this *hardwareDriver) Close() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - DISPLAY

// Open
func (config Display) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.mock.Display.Open{ }")

	this := new(displayDriver)
	this.id = config.Display

	// Success
	return this, nil
}

// Close
func (this *displayDriver) Close() error {
	return nil
}

// Return display number
func (this *displayDriver) Display() uint {
	return 0
}

// Return size
func (this *displayDriver) Size() (uint32, uint32) {
	return 0, 0
}

// Return pixels-per-inch
func (this *displayDriver) PixelsPerInch() uint32 {
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *hardwareDriver) String() string {
	return fmt.Sprintf("sys.mock.Hardware{ name=%v serial=%v displays=%v }", this.Name(), this.SerialNumber(), this.NumberOfDisplays())
}

func (this *displayDriver) String() string {
	return fmt.Sprintf("sys.mock.Display{ id=%v }", this.id)
}
