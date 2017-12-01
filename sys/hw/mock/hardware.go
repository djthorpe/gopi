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
	// Register hardware and display
	gopi.RegisterModule(gopi.Module{Name: "mock/hw", Type: gopi.MODULE_TYPE_HARDWARE, New: newHardware})
	registerDisplayFlags(gopi.RegisterModule(gopi.Module{
		Name: "mock/display",
		Type: gopi.MODULE_TYPE_DISPLAY,
		New:  newDisplay,
	}))
}

func registerDisplayFlags(flags *gopi.Flags) {
	flags.FlagUint("display", 0, "Display")
}

func newHardware(config *gopi.AppConfig, logger gopi.Logger) (gopi.Driver, error) {
	var err gopi.Error
	if driver, ok := gopi.Open2(Hardware{}, logger, &err).(gopi.HardwareDriver2); !ok {
		return nil, err
	} else {
		return driver, nil
	}
}

func newDisplay(config *gopi.AppConfig, logger gopi.Logger) (gopi.Driver, error) {
	var err gopi.Error
	var display Display

	// set display argument
	if display_number, exists := config.AppFlags.GetUint("display"); exists {
		display.Display = display_number
	}
	if driver, ok := gopi.Open2(display, logger, &err).(gopi.DisplayDriver2); !ok {
		return nil, err
	} else {
		return driver, nil
	}
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
