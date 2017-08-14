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

	// Attempt to create a view with tag TagNone
	view2 := layout.NewRootView(gopi.TagNone, "")
	if view2 != nil {
		t.Error("NewRootView succeeded but should have failed")
	}
}

func TestLayout_003(t *testing.T) {
	// Check class names
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

func TestLayout_004(t *testing.T) {
	// Check layout starts as absolute with auto edges
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
	view := layout.NewRootView(1, "root")
	if view == nil {
		t.Error("NewRootView failed")
	}
	if view.Positioning() != gopi.VIEW_POSITIONING_ABSOLUTE {
		t.Error("Expected positioning on root element to be absolute")
	}
	app.Logger.Info("view=%v", view)
}

func TestLayout_005(t *testing.T) {
	m := map[gopi.ViewDirection]string{
		gopi.VIEW_DIRECTION_COLUMN:         "VIEW_DIRECTION_COLUMN",
		gopi.VIEW_DIRECTION_COLUMN_REVERSE: "VIEW_DIRECTION_COLUMN_REVERSE",
		gopi.VIEW_DIRECTION_ROW:            "VIEW_DIRECTION_ROW",
		gopi.VIEW_DIRECTION_ROW_REVERSE:    "VIEW_DIRECTION_ROW_REVERSE",
	}
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
	// Create a view and test direction property
	if view := app.Layout.NewRootView(1, "root"); view == nil {
		t.Error("View could not be created")
	} else {
		for k, v := range m {
			view.SetDirection(k)
			if k.String() != v {
				t.Errorf("Expected string to return %v but it returned %v", v, k.String())
			}
			if view.Direction() != k {
				t.Errorf("Expected Direction() to return %v but it returned %v", k, view.Direction())
			}
		}
	}
}

func TestLayout_010(t *testing.T) {
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
	// Create a view and test direction property
	view := app.Layout.NewRootView(1, "root")
	if view == nil {
		t.Fatal("View could not be created")
	}

	// Output view
	app.Logger.Info("view=%v", view)
}
