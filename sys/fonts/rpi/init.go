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
	// Register font manager
	gopi.RegisterModule(gopi.Module{
		Name: "rpi/fonts",
		Type: gopi.MODULE_TYPE_FONTS,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("fonts.path", "", "Path for font files")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			font_path, _ := app.AppFlags.GetString("fonts.path")
			return gopi.Open(FontManager{
				RootPath: font_path,
			}, app.Logger)
		},
	})
}
