package gopi_test

import (
	"fmt"
	"testing"

	"github.com/djthorpe/gopi"
)

import (
	_ "github.com/djthorpe/gopi/sys/mac"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE APPLICATION TESTS

func TestCreateApp_000(t *testing.T) {
	// Create an application with an empty configuration
	if _, err := gopi.NewAppInstance(gopi.NewAppConfig()); err != nil {
		t.Error(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// RUN COMMAND-LINE APP

func TestRunApp_000(t *testing.T) {
	// Create an application with an empty configuration
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig()); err != nil {
		t.Error(err)
	} else if err := app.Run(HelloWorld); err != nil {
		t.Error(err)
	}
}

func TestRunApp_001(t *testing.T) {
	// Create an application with an empty configuration
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig()); err != nil {
		t.Error(err)
	} else if err := app.Run(ReturnAnError); err != gopi.ErrAppError {
		t.Error("Unexpected error returned")
	}
}

////////////////////////////////////////////////////////////////////////////////
// TASKS

func HelloWorld(app *gopi.AppInstance, done chan struct{}) error {
	fmt.Println("Hello, World")
	done <- gopi.DONE
	return nil
}

func ReturnAnError(app *gopi.AppInstance, done chan struct{}) error {
	done <- gopi.DONE
	return gopi.ErrAppError
}
