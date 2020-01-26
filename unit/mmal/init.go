// +build mmal

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mmal

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: MMAL{}.Name(),
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(MMAL{}, app.Log().Clone(MMAL{}.Name()))
		},
	})
}
