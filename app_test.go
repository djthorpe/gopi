package gopi_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/djthorpe/gopi"

	_ "github.com/djthorpe/gopi/sys/mock"
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

func TestRunApp_002(t *testing.T) {
	// Create an application with an empty configuration, make sure logger
	// device is also created
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig()); err != nil {
		t.Error(err)
	} else if err := app.Run(CheckLogger); err != nil {
		t.Error(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// RUN BACKGROUND TASKS
func TestRunTasks_001(t *testing.T) {
	// Create an application with an empty configuration, make sure logger
	// device is also created
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(Task001, Task002); err != nil {
		t.Error(err)
	}
}

func TestRunTasks_002(t *testing.T) {
	// Create an application with two sub-tasks
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(Task001, Task002, Task003); err != nil {
		t.Error(err)
	}
}

func TestRunTasks_003(t *testing.T) {
	// Create an application with no tasks, which should return an error
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(); err != gopi.ErrNoTasks {
		t.Error("Unexpected response, expected gopi.ErrNoTasks")
	}
}

func TestRunTasks_004(t *testing.T) {
	// Have a subtask wait for 200ms before finishing
	timer := time.Now()
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(Task001, Task002); err != nil {
		t.Error(err)
	}
	if time.Since(timer) < time.Duration(200*time.Millisecond) {
		t.Error("Ended too early, expected to wait for 200ms")
	}
}

func TestRunTasks_005(t *testing.T) {
	// Have a subtask wait for 1 second before finishing
	timer := time.Now()
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(Task001, Task002, Task003, Task010); err != nil {
		t.Error(err)
	}
	if time.Since(timer) < time.Duration(1000*time.Millisecond) {
		t.Error("Ended too early, expected to wait for 1000ms")
	}
}

////////////////////////////////////////////////////////////////////////////////
// SIGNALLING TESTS

func TestRunSignal_001(t *testing.T) {
	// Have main thread signal it is done
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(MainThreadSignal); err != nil {
		t.Error(err)
	}
}

func TestRunSignal_002(t *testing.T) {
	// Have main thread signal it is done to one sub-task
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(MainThreadSignal, WaitTask001); err != nil {
		t.Error(err)
	}
}

func TestRunSignal_003(t *testing.T) {
	// Have main thread signal and two sub-tasks
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(MainThreadSignal, WaitTask001, WaitTask002); err != nil {
		t.Error(err)
	}
}

func TestRunSignal_004(t *testing.T) {
	// Have main thread signal and two sub-tasks, plus a task which does occassional work
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(MainThreadSignal, WaitTask001, WaitTask002, WaitTask003); err != nil {
		t.Error(err)
	}
}

func TestRunSignal_005(t *testing.T) {
	// Have main thread signal and two sub-tasks, plus two tasks which do occassional work
	config := gopi.NewAppConfig()
	config.Debug = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else if err := app.Run(MainThreadSignal, WaitTask001, WaitTask002, WaitTask003, WaitTask004); err != nil {
		t.Error(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TASKS

func HelloWorld(app *gopi.AppInstance, done chan struct{}) error {
	fmt.Println("Hello, World")
	return nil
}

func ReturnAnError(app *gopi.AppInstance, done chan struct{}) error {
	return gopi.ErrAppError
}

func CheckLogger(app *gopi.AppInstance, done chan struct{}) error {
	if app.Logger == nil {
		return errors.New("Expected a logger object")
	}
	return nil
}

func Task001(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Info("Running Task 001")
	time.Sleep(100 * time.Millisecond)
	return nil
}

func Task002(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Info("Running Task 002")
	time.Sleep(200 * time.Millisecond)
	return nil
}

func Task003(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Info("Running Task 003")
	time.Sleep(300 * time.Millisecond)
	return nil
}

func Task010(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Info("Running Task 010")
	time.Sleep(1000 * time.Millisecond)
	return nil
}

func MainThreadSignal(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Main thread: waiting for 1 second to complete")
	time.Sleep(1 * time.Second)
	app.Logger.Debug("Main thread: signalling we are done")
	done <- gopi.DONE
	app.Logger.Debug("Main thread: signaled we are done")
	return nil
}

func WaitTask001(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Wait Task 001 thread: waiting for done")
	_ = <-done
	app.Logger.Debug("Wait Task 001 thread: got done")
	return nil
}

func WaitTask002(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Wait Task 002 thread: waiting for done")
	_ = <-done
	app.Logger.Debug("Wait Task 002 thread: got done, waiting for another 1 second")
	time.Sleep(1 * time.Second)
	app.Logger.Debug("Wait Task 002 thread: now finished")
	return nil
}

func WaitTask003(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Debug("WaitTask003 thread: doing work occasionally")
	t := time.NewTicker(200 * time.Millisecond)
OUTER_LOOP:
	for {
		select {
		case <-t.C:
			app.Logger.Debug("WaitTask003 Tick")
		case <-done:
			app.Logger.Debug("WaitTask003 Done")
			t.Stop()
			break OUTER_LOOP
		}
	}
	app.Logger.Debug("WaitTask003: now finished")
	return nil
}

func WaitTask004(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Debug("WaitTask004 thread: doing work occasionally")
	t := time.NewTicker(200 * time.Millisecond)
OUTER_LOOP:
	for {
		select {
		case <-t.C:
			app.Logger.Debug("WaitTask004 Start work")
			time.Sleep(500 * time.Millisecond)
			app.Logger.Debug("WaitTask004 End work")
		case <-done:
			app.Logger.Debug("WaitTask004 Done")
			t.Stop()
			break OUTER_LOOP
		}
	}
	app.Logger.Debug("WaitTask004: now finished")
	return nil
}
