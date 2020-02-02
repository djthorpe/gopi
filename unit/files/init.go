/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: FilePoll{}.Name(),
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(FilePoll{}, app.Log().Clone(FilePoll{}.Name()))
		},
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     FSEvents{}.Name(),
		Requires: []string{FilePoll{}.Name()},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(FSEvents{
				FilePoll: app.UnitInstance(FilePoll{}.Name()).(gopi.FilePoll),
			}, app.Log().Clone(FSEvents{}.Name()))
		},
	})
}

