/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package timer

import (
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "gopi/timer",
		Type:     gopi.UNIT_TIMER,
		Requires: []string{"bus"},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Timer{
				Bus: app.UnitInstance("bus").(gopi.Bus),
			}, app.Log())
		},
	})
}
