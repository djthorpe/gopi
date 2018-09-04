// +build rpi

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
		Name: "rpi/hw",
		Type: gopi.MODULE_TYPE_HARDWARE,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Hardware{}, app.Logger)
		},
	})

	// Register Display
	/*
		gopi.RegisterModule(gopi.Module{
			Name:     "rpi/display",
			Type:     gopi.MODULE_TYPE_DISPLAY,
			Requires: []string{"rpi/hw"},
			Config: func(config *gopi.AppConfig) {
				config.AppFlags.FlagUint("display", 0, "Display")
				config.AppFlags.FlagString("display.ppi", "", "Pixels per inch or diagonal size in mm, cm or in")
			},
			New: func(app *gopi.AppInstance) (gopi.Driver, error) {
				display := Display{}
				if display_number, exists := app.AppFlags.GetUint("display"); exists {
					display.Display = display_number
				}
				if pixels_per_inch, exists := app.AppFlags.GetString("display.ppi"); exists {
					display.PixelsPerInch = pixels_per_inch
				}
				return gopi.Open(display, app.Logger)
			},
		})
	*/

	// Register GPIO
	gopi.RegisterModule(gopi.Module{
		Name:     "rpi/gpio",
		Type:     gopi.MODULE_TYPE_GPIO,
		Requires: []string{"rpi/hw"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(GPIO{Hardware: app.Hardware}, app.Logger)
		},
	})
}
