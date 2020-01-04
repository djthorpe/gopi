/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package platform_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Platform_000(t *testing.T) {
	t.Log("Test_Platform_000")
}

func Test_Platform_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Platform_001, []string{"-debug"}, "platform"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Platform_001(app gopi.App, t *testing.T) {
	platform := app.Platform()
	if platform == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("Platform() failed"))
	}
	app.Log().Debug("Platform", platform)
	if type_ := platform.Type(); type_ == 0 {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("Type() failed"))
	} else {
		t.Log("Type", type_)
	}
	if platform.SerialNumber() == "" {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("SerialNumber() failed"))
	} else {
		t.Log("SerialNumber", platform.SerialNumber())
	}
	if platform.Uptime() == 0 {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("Uptime() failed"))
	} else {
		t.Log("Uptime", platform.Uptime())
	}
	if l1, l5, l15 := platform.LoadAverages(); l1 == 0 || l5 == 0 || l15 == 0 {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("LoadAverages() failed"))
	} else {
		t.Log("Load Average", l1, l5, l15)
	}

	t.Log("Number of displays", platform.NumberOfDisplays())
}
