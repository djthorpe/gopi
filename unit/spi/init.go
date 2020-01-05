/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package spi

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/spi",
		Type: gopi.UNIT_SPI,
		Config: func(app gopi.App) error {
			app.Flags().FlagUint("spi.bus", 0, "SPI Bus")
			app.Flags().FlagUint("spi.slave", 0, "SPI Slave")
			app.Flags().FlagUint("spi.delay", 0, "SPI Transfer delay in microseconds")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(SPI{
				Bus:   app.Flags().GetUint("spi.bus", gopi.FLAG_NS_DEFAULT),
				Slave: app.Flags().GetUint("spi.slave", gopi.FLAG_NS_DEFAULT),
				Delay: app.Flags().GetUint("spi.delay", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone("gopi/spi"))
		},
	})
}
