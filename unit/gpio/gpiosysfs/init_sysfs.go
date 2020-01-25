// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiosysfs

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     GPIO{}.Name(),
		Type:     gopi.UNIT_GPIO,
		Requires: []string{"gopi/filepoll"},
		Config: func(app gopi.App) error {
			app.Flags().FlagBool("gpio.unexport", false, "Unexport exported pins on exit")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(GPIO{
				FilePoll:        app.UnitInstance("gopi/filepoll").(gopi.FilePoll),
				UnexportOnClose: app.Flags().GetBool("gpio.unexport", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone(GPIO{}.Name()))
		},
	})
}
