// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/fonts/freetype"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Freetype_000(t *testing.T) {
	t.Log("Test_Freetype_000")
}

func Test_Freetype_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Freetype_000, []string{"-debug"}, "fonts"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Freetype_000(app gopi.App, t *testing.T) {
	fonts := app.Fonts()
	t.Log(fonts)
}
