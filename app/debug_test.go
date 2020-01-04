/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Debug_000(t *testing.T) {
	t.Log("Test_Debug_000")
}

func Test_Debug_001(t *testing.T) {
	args := []string{}
	units := []string{}
	if app, err := app.NewTestTool(t, func(app gopi.App, t *testing.T) {
		t.Log("In", app.Flags().Name())
	}, args, units...); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}
