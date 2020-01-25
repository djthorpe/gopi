/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpio

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: GPIO{}.Name(),
		Type: gopi.UNIT_GPIO,
		Pri:  1,
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(GPIO{}, app.Log().Clone(GPIO{}.Name()))
		},
		Run: func(app gopi.App, unit gopi.Unit) error {
			this := unit.(*gpio)
			if rpi := app.UnitInstance("gopi/gpio/rpi"); rpi != nil {
				this.rpi = rpi.(gopi.GPIO)
			}
			if sysfs := app.UnitInstance("gopi/gpio/sysfs"); sysfs != nil {
				this.sysfs = sysfs.(gopi.GPIO)
			}
			return nil
		},
	})
}
