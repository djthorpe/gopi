package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/default/layout"
	_ "github.com/djthorpe/gopi/sys/default/logger"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE LAYOUT MODULE

func TestLayout_000(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	if app == nil {
		t.Fatal("Expecting app object")
	}
	if app.Logger == nil {
		t.Fatal("Expecting app.Logger object")
	}
	if app.Layout == nil {
		t.Fatal("Expecting app.Layout object")
	}
	app.Logger.Info("layout=%v", app.Layout)
}

func TestLayout_001(t *testing.T) {
	// Create a root node with tag 0
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()

	// Check layout direction defaults to LEFTRIGHT
	layout := app.Layout
	default_direction := gopi.LAYOUT_DIRECTION_LEFTRIGHT
	if layout.Direction() != default_direction {
		t.Errorf("Layout direction is %v, expected %v", layout.Direction(), default_direction)
	}
}
