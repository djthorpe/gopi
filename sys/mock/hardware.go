package mock /* import "github.com/djthorpe/gopi/sys/mock" */

import (
	"fmt"
	"strings"

	"github.com/djthorpe/gopi"
	"github.com/rs/xid"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type hardwareDriver struct{}
type displayDriver struct{}

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware and display
	gopi.RegisterModule(gopi.Module{Name: "mock/hardware", Type: gopi.MODULE_TYPE_HARDWARE, New: newHardware})
	gopi.RegisterModule(gopi.Module{Name: "mock/display", Type: gopi.MODULE_TYPE_DISPLAY, New: newDisplay})
}

func newHardware(config *gopi.AppConfig) (gopi.Driver, error) {
	return new(hardwareDriver), nil
}

func newDisplay(config *gopi.AppConfig) (gopi.Driver, error) {
	return new(displayDriver), nil
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	guid = xid.New() // unique id we can use as a fake serial number
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - HARDWARE

// GetName returns the name of the hardware (ie, mock, mac, linux, rpi, etc)
func (this *hardwareDriver) Name() string {
	return "mock/hardware"
}

// GetName returns the name of the hardware (ie, mock, mac, linux, rpi, etc)
func (this *hardwareDriver) SerialNumber() string {
	return strings.ToUpper(guid.String())
}

func (this *hardwareDriver) NumberOfDisplays() uint {
	return 0
}

func (this *hardwareDriver) Close() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - DISPLAY

func (this *displayDriver) Close() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *hardwareDriver) String() string {
	return fmt.Sprintf("sys.mock.Hardware{ name=%v serial=%v displays=%v }", this.Name(), this.SerialNumber(), this.NumberOfDisplays())
}

func (this *displayDriver) String() string {
	return fmt.Sprintf("sys.mock.Display{ }")
}
