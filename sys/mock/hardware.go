package mock /* import "github.com/djthorpe/gopi/sys/mock" */

import (
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware
	gopi.RegisterModule(gopi.Module{Name: "mock/hardware", Type: gopi.MODULE_TYPE_HARDWARE, New: newHardware})
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func newHardware(config *gopi.AppConfig) (gopi.Driver, error) {
	return nil, fmt.Errorf("Not implemented")
}
