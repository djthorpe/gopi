/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "gopi/surfaces",
		Type:     gopi.UNIT_SURFACE_MANAGER,
		Requires: []string{"display"},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(SurfaceManager{
				Display: app.Display(),
			}, app.Log().Clone("gopi/surfaces"))
		},
	})
}
