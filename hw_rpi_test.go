package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/rpi"
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
