/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package i2c

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/i2c",
		Type: gopi.UNIT_I2C,
		Config: func(app gopi.App) error {
			app.Flags().FlagUint("i2c.bus", 1, "I2C Bus")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(I2C{
				Bus: app.Flags().GetUint("i2c.bus", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone("gopi/i2c"))
		},
	})
}
