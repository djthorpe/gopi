package gopi_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE LOGGER

func TestLogging_000(t *testing.T) {
	// Create config file with logger
	gopi.NewAppConfig("logger")
}
func TestLogging_001(t *testing.T) {
	// Get logging module
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig("logger")); err != nil {
		t.Error(err)
	} else if app == nil {
		t.Error("app==nil")
	} else if app.Logger == nil {
		t.Error("app.Logger==nil")
	}
}

func TestLogging_002(t *testing.T) {
	// Check IsDebug flag
	config := gopi.NewAppConfig("logger")
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if app.Logger.IsDebug() == false {
		t.Error("app.Logger.IsDebug() == false, expected true")
	}
}

func TestLogging_003(t *testing.T) {
	// Check IsDebug flag
	config := gopi.NewAppConfig("logger")
	config.Debug = false
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if app.Logger.IsDebug() == true {
		t.Error("app.Logger.IsDebug() == true, expected false")
	}
}
