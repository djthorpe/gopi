/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     InputManager{}.Name(),
		Type:     gopi.UNIT_INPUT_MANAGER,
		Requires: []string{"gopi/filepoll", "bus"},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(InputManager{
				FilePoll: app.UnitInstance("gopi/filepoll").(gopi.FilePoll),
				Bus:      app.UnitInstance("bus").(gopi.Bus),
			}, app.Log().Clone(InputManager{}.Name()))
		},
	})
}
