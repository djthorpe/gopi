// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiorpi

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: GPIO{}.Name(),
		Type: gopi.UNIT_GPIO,
		Requires: []string{ "platform" },
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(GPIO{
				Platform: app.Platform(),
			}, app.Log().Clone(GPIO{}.Name()))
		},
	})
}
