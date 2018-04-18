/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package timer

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register Timer
	gopi.RegisterModule(gopi.Module{
		Name: "sys/timer",
		Type: gopi.MODULE_TYPE_TIMER,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Timer{}, app.Logger)
		},
	})
}
