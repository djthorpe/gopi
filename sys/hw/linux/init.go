/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register GPIO
	gopi.RegisterModule(gopi.Module{
		Name:     "gpio/linux",
		Type:     gopi.MODULE_TYPE_GPIO,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(GPIO{}, app.Logger)
		},
	})
}
