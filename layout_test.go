package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/default/layout"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE LAYOUT MODULE

func TestLayout_000(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Error(err)
	}
	if app == nil {
		t.Error("Expecting app object")
	}
	if app.Logger == nil {
		t.Error("Expecting app.Logger object")
	}
	if app.Layout == nil {
		t.Error("Expecting app.Layout object")
	}
	app.Logger.Info("layout=%v", app.Layout)
}
