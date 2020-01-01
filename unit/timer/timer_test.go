/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package timer_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/timer"
)

func Test_Timer_000(t *testing.T) {
	t.Log("Test_Timer_000")
}

func Test_Timer_001(t *testing.T) {
	if app, err := app.NewCommandLineTool(func(app gopi.App, _ []string) error {
		if timer := app.Timer(); timer == nil {
			t.Error("nil timer unit")
		}
		return nil
	}, "timer", "bus"); err != nil {
		t.Error(err)
	} else if returnValue := app.Run(); returnValue != 0 {
		t.Error("Unexpected return value")
	}
}
