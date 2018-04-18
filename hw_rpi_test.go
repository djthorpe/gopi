// +build rpi

package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	rpi "github.com/djthorpe/gopi/sys/hw/rpi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE HARDWARE MODULE
func TestHardware_000(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("hw")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer app.Close()
	if app == nil {
		t.Fatal("Expecting app object")
		return
	}
	if app.Logger == nil {
		t.Fatal("Expecting app.Logger object")
		return
	}
	if app.Hardware == nil {
		t.Fatal("Expecting app.Hardware object")
		return
	}
	app.Logger.Info("hardware=%v", app.Hardware)
}

////////////////////////////////////////////////////////////////////////////////
// DISPLAY TEST FUNCTIONS

func TestDisplay_000(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("display")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer app.Close()
	if app == nil {
		t.Fatal("Expecting app object")
		return
	}
	if app.Logger == nil {
		t.Fatal("Expecting app.Logger object")
		return
	}
	if app.Hardware == nil {
		t.Fatal("Expecting app.Hardware object")
		return
	}
	if app.Display == nil {
		t.Fatal("Expecting app.Display object")
		return
	}
	app.Logger.Info("display=%v", app.Display)
}
func TestDisplay_001(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("hw")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer app.Close()

	// Number of displays should be greater than zero
	if app.Hardware.NumberOfDisplays() == 0 {
		t.Error("Expected non-zero return fromNumberOfDisplays, actual return value is ", app.Hardware.NumberOfDisplays())
	}
}

func TestDisplay_002(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("display")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer app.Close()

	// Open up displays and check the dimensions of each
	num_displays := app.Hardware.NumberOfDisplays()
	for display := uint(0); display < num_displays; display++ {
		if d, err := gopi.Open(rpi.Display{Display: uint(display)}, app.Logger); err != nil {
			t.Error(err)
		} else {
			defer d.Close()

			if display != d.(gopi.Display).Display() {
				t.Errorf("Expected Display() to return %v but it returned %v", display, d.(gopi.Display).Display())
			}

			w, h := d.(gopi.Display).Size()
			t.Logf("Display=%v Size={%v,%v}", display, w, h)
			if w == 0 || h == 0 {
				t.Errorf("Expected non-zero size, got size={ %v,%v }", w, h)
			}
		}
	}
}
