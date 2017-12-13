/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package mock

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hw
	gopi.RegisterModule(gopi.Module{
		Name: "hw/mock",
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

	// Register gpio
	gopi.RegisterModule(gopi.Module{
		Name: "gpio/mock",
		Type: gopi.MODULE_TYPE_GPIO,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(GPIO{}, app.Logger)
		},
	})
}
