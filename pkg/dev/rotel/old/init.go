/*
  Mutablehome Automation: Rotel
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: Rotel{}.Name(),
		Config: func(app gopi.App) error {
			app.Flags().FlagString("rotel.tty", "/dev/ttyUSB0", "RS232 device")
			FlagUint("rotel.baudrate", BAUD_RATE_DEFAULT, "RS232 speed")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Rotel{
				TTY:      app.Flags().GetString("rotel.tty", gopi.FLAG_NS_DEFAULT),
				BaudRate: app.Flags().GetUint("rotel.baudrate", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone(Rotel{}.Name()))
		},
	})
}
