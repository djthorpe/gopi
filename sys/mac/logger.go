package mac /* import "github.com/djthorpe/gopi/sys/mac" */

import (
	"errors"
)

import (
	gopi "github.com/djthorpe/gopi"
)

func init() {
	// Register logger
	gopi.RegisterModule(gopi.Module{Name: "mac/logger", Type: gopi.MODULE_TYPE_LOGGER, Config: registerFlags, New: newLogger})
}

func registerFlags(config *gopi.AppConfig) {
	config.Flags.FlagString("log", "", "File for logging (default: log to stderr)")
	config.Flags.FlagBool("verbose", false, "Log verbosely")
	config.Flags.FlagBool("debug", false, "Trigger debugging support")
}

func newLogger(config *gopi.AppConfig) (gopi.Driver, error) {
	return nil, errors.New("Not implemented")
}
