// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package display_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/display"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/platform"
)

func Test_Display_000(t *testing.T) {
	t.Log("Test_Display_000")
}

func Test_Display_001(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewDebugTool(Main_Test_Display_001, flags, []string{"display"}); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Display_001(app gopi.App, _ []string) error {
	display := app.Display()
	if display == nil {
		return gopi.ErrInternalAppError.WithPrefix("Display() failed")
	}
	app.Log().Debug(display)
	// Success
	return nil
}
