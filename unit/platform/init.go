/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package logger

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/platform",
		Type: gopi.UNIT_PLATFORM,
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Platform{}, nil)
		},
	})
}
