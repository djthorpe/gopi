/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Filepoll_000(t *testing.T) {
	t.Log("Test_Filepoll_000")
}

func Test_Filepoll_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Filepoll_001, []string{"-debug"}, "gopi/filepoll"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Filepoll_001(app gopi.App, t *testing.T) {
	filepoll := app.UnitInstance("gopi/filepoll").(gopi.FilePoll)
	t.Log(filepoll)
}
