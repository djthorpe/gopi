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

	// Create an application with a layout module
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
	// Check direction default
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a layout module
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

func TestLayout_002(t *testing.T) {
	// Create a root node with tag 1
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()

	// Create a view
	layout := app.Layout
	view1 := layout.NewRootView(1, "root")
	if view1 == nil {
		t.Error("NewRootView failed")
	}
	if view1.Tag() != 1 {
		t.Errorf("view1.Tag() expected 1, received %v", view1.Tag())
	}
	if view1.Class() != "root" {
		t.Errorf("view1.Tag() expected root, received %v", view1.Class())
	}
}

func TestLayout_003(t *testing.T) {
	// Check class names
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

	class_name_tests := map[string]bool{
		"":       false,
		"a":      true,
		"-":      false,
		"0":      false,
		"test":   true,
		"t0":     true,
		"t-":     true,
		"t-test": true,
		"t!test": false,
	}

	// Create root view with particular class names
	tag := uint(1)
	for k, v := range class_name_tests {
		view := app.Layout.NewRootView(tag, k)
		failed := (view == nil)
		if failed == v {
			t.Errorf("class %v => %v, expected %v", k, !failed, v)
		}
		if view != nil {
			if view.Tag() != tag {
				t.Errorf("view.Tag() expected %v, received %v", tag, view.Tag())
			}
			if view.Class() != k {
				t.Errorf("view.Class() expected %v, received %v", k, view.Class())
			}
		}
		tag = tag + 1
	}
}
