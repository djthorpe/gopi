/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware
	gopi.RegisterModule(gopi.Module{
		Name: "hw/rpi",
		Type: gopi.MODULE_TYPE_HARDWARE,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Hardware{}, app.Logger)
		},
	})

	// Register GPIO
	gopi.RegisterModule(gopi.Module{
		Name: "gpio/rpi",
		Type: gopi.MODULE_TYPE_GPIO,
		Requires: []string{ "hw/rpi" }
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(GPIO{}, app.Logger)
		},
	})	
}
