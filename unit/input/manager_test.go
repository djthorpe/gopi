// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input_test

import (
	"os"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/gopi/v2/unit/input"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Input_000(t *testing.T) {
	t.Log("Test_Input_000")
}

func Test_Input_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Input_001, []string{"-debug"}, "input"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Input_001(app gopi.App, t *testing.T) {
	input := app.Input()
	if input == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("UnitInstance() failed"))
	} else {
		for bus := uint(0); bus < 10; bus++ {
			if device, err := input.OpenDevice(bus, false); err != nil && os.IsNotExist(err) == false {
				t.Error(err)
			} else if device != nil {
				t.Log(device)
			}
		}
	}
}

func Test_Input_002(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Input_002, []string{"-debug"}, "input"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Input_002(app gopi.App, t *testing.T) {
	input := app.Input()
	if input == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("UnitInstance() failed"))
	} else if devices, err := input.OpenDevicesByNameType("", gopi.INPUT_TYPE_ANY, false); err != nil {
		t.Error(err)
	} else {
		t.Log(devices)
	}
}
