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
		Name:     "gopi/lirc",
		Type:     gopi.UNIT_LIRC,
		Requires: []string{"bus", "gopi/filepoll"},
		Config: func(app gopi.App) error {
			app.Flags().FlagString("lirc.dev", "0,1", "Comma-separated list of LIRC devices")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(LIRC{
				Dev:      app.Flags().GetString("lirc.dev", gopi.FLAG_NS_DEFAULT),
				Bus:      app.Bus(),
				Filepoll: app.UnitInstance("gopi/filepoll").(gopi.FilePoll),
			}, app.Log().Clone("gopi/lirc"))
		},
	})
}
