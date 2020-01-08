// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/fonts/freetype",
		Type: gopi.UNIT_FONT_MANAGER,
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(FontManager{}, app.Log().Clone("gopi/fonts/freetype"))
		},
	})
}
