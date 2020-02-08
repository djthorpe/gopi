/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package display

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Display{}.Name(),
		Type:     gopi.UNIT_DISPLAY,
		Requires: []string{"platform"},
		Config: func(app gopi.App) error {
			app.Flags().FlagUint("display", 0, "Display")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Display{
				Id:       app.Flags().GetUint("display", gopi.FLAG_NS_DEFAULT),
				Platform: app.Platform(),
			}, app.Log().Clone(Display{}.Name()))
		},
	})
}
