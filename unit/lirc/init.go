/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package lirc

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/lirc",
		Type: gopi.UNIT_LIRC,
		Config: func(app gopi.App) error {
			app.Flags().FlagString("lirc.in", "", "LIRC input device")
			app.Flags().FlagString("lirc.out", "", "LIRC output device")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(LIRC{
				DevIn:  app.Flags().GetString("lirc.in", gopi.FLAG_NS_DEFAULT),
				DevOut: app.Flags().GetString("lirc.out", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone("gopi/lirc"))
		},
	})
}
