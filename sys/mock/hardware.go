package mock /* import "github.com/djthorpe/gopi/sys/mock" */

import (
	"fmt"

	"github.com/djthorpe/gopi"
	"github.com/rs/xid"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type hardwareDriver struct{}

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware
	gopi.RegisterModule(gopi.Module{Name: "mock/hardware", Type: gopi.MODULE_TYPE_HARDWARE, New: newHardware})
}

func newHardware(config *gopi.AppConfig) (gopi.Driver, error) {
	return new(hardwareDriver), nil
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	guid = xid.New()
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetName returns the name of the hardware (ie, mock, mac, linux, rpi, etc)
func (this *hardwareDriver) GetName() string {
	return "mock/hardware"
}

func (this *hardwareDriver) GetSerialNumber() string {
	return guid.String()
}

func (this *hardwareDriver) NumberOfDisplays() uint {
	return 0
}

func (this *hardwareDriver) Close() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *hardwareDriver) String() string {
	return fmt.Sprintf("sys.mock.Hardware{ name=%v serial=%v displays=%v }", this.GetName(), this.GetSerialNumber(), this.NumberOfDisplays())
}
