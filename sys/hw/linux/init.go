// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register FilePoll
	gopi.RegisterModule(gopi.Module{
		Name: "filepoll/linux",
		Type: gopi.MODULE_TYPE_OTHER,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(FilePoll{}, app.Logger)
		},
	})

	// Register GPIO
	gopi.RegisterModule(gopi.Module{
		Name:     "gpio/linux",
		Requires: []string{"filepoll/linux"},
		Type:     gopi.MODULE_TYPE_GPIO,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagBool("gpio.unexport", true, "Unexport exported pins on exit")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			unexport, _ := app.AppFlags.GetBool("gpio.unexport")
			return gopi.Open(GPIO{
				UnexportOnClose: unexport,
				FilePoll:        app.ModuleInstance("filepoll/linux").(FilePollInterface),
			}, app.Logger)
		},
	})
}
