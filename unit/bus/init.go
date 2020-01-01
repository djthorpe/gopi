/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bus

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/bus",
		Type: gopi.UNIT_BUS,
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Bus{}, app.Log())
		},
	})
}
