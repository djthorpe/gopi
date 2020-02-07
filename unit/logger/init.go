/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package logger

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: Log{}.Name(),
		Type: gopi.UNIT_LOGGER,
		Pri:  1,
		Config: func(app gopi.App) error {
			app.Flags().FlagBool("verbose", true, "Verbose output")
			app.Flags().FlagBool("debug", false, "Debugging output")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Log{
				Unit:    app.Flags().Name(),
				Debug:   app.Flags().GetBool("debug", gopi.FLAG_NS_DEFAULT),
				Verbose: app.Flags().GetBool("verbose", gopi.FLAG_NS_DEFAULT),
				Writer:  os.Stderr,
			}, nil)
		},
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name: TestLogger{}.Name(),
		Type: gopi.UNIT_LOGGER,
		New: func(app gopi.App) (gopi.Unit, error) {
			_ = app.(gopi.DebugApp)
			if testapp, ok := app.(gopi.DebugApp); ok {
				return gopi.New(TestLogger{
					Unit: app.Flags().Name(),
					T:    testapp.T(),
				}, nil)
			} else {
				return nil, fmt.Errorf("Can't use %v without gopi.DebugApp", TestLogger{}.Name())
			}
		},
	})
}
