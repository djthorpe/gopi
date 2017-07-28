package app2_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	app "github.com/djthorpe/gopi/app2"
)

func TestCreateApp_000(t *testing.T) {
	// Create an application with an empty configuration
	if _, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	}
}

func TestRunApp_000(t *testing.T) {
	// Create an application with an empty configuration
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(HelloWorld); err != nil {
		t.Error(err)
	}
}

func TestRunApp_001(t *testing.T) {
	// Create an application with an empty configuration, one main
	// task and one sub task
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(HelloWorld, Task001); err != nil {
		t.Error(err)
	}
}

func TestRunApp_002(t *testing.T) {
	// Create an application with an empty configuration and
	// one main task plus three sub-tasks
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(HelloWorld, Task001, Task002, Task003); err != nil {
		t.Error(err)
	}
}
func TestRunApp_003(t *testing.T) {
	// Create an application with an empty configuration
	// call run but with no tasks to run...
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(); err != app.ErrNoTasks {
		t.Error(errors.New("Expected app.ErrNoTasks error code"))
	}
}

func TestRunApp_004(t *testing.T) {
	// Have a subtask wait for 500ms before finishing
	timer := time.Now()
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(HelloWorld, Task004); err != nil {
		t.Error(err)
	}
	if time.Since(timer) < time.Duration(500*time.Millisecond) {
		t.Error("Ended too early")
	}
}

func TestRunApp_005(t *testing.T) {
	// Have a subtask wait for 1 second before finishing
	timer := time.Now()
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(HelloWorld, Task004, Task005); err != nil {
		t.Error(err)
	}
	if time.Since(timer) < time.Duration(1000*time.Millisecond) {
		t.Error("Ended too early")
	}
}

func TestAppDone_001(t *testing.T) {
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(MainThreadSignal); err != nil {
		t.Error(err)
	}
}

func TestAppDone_002(t *testing.T) {
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(MainThreadSignal, Task001Signal); err != nil {
		t.Error(err)
	}
}

func TestAppDone_003(t *testing.T) {
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(MainThreadSignal, Task001Signal, Task002Signal); err != nil {
		t.Error(err)
	}
}

func TestAppDone_004(t *testing.T) {
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(MainThreadSignal, Task001Signal, Task002Signal, Task003Signal); err != nil {
		t.Error(err)
	}
}

func TestAppDone_005(t *testing.T) {
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(MainThreadSignal, Task004Signal); err != nil {
		t.Error(err)
	}
}

func TestAppDone_006(t *testing.T) {
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(MainThreadSignal, Task004Signal, Task005Signal); err != nil {
		t.Error(err)
	}
}

func TestAppDone_007(t *testing.T) {
	if myapp, err := app.NewAppInstance(app.AppConfig{}); err != nil {
		t.Error(err)
	} else if err := myapp.Run(MainThreadSignal, Task001Signal, Task002Signal, Task003Signal, Task004Signal, Task005Signal); err != nil {
		t.Error(err)
	}
}

//////////////////////////////////////////////////////////////////////////////////////

func HelloWorld(app *app.AppInstance, done chan struct{}) error {
	if app.Logger == nil {
		return fmt.Errorf("Invalid or empty logger objet")
	}
	app.Logger.Debug("Hello, World")
	return nil
}

func Task001(app *app.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Task 001")
	return nil
}

func Task002(app *app.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Task 002")
	return nil
}

func Task003(app *app.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Task 003")
	return nil
}

func Task004(app *app.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Task 004 - wait for 0.5 second")
	time.Sleep(500 * time.Millisecond)
	return nil
}

func Task005(app *app.AppInstance, done chan struct{}) error {
	app.Logger.Debug("Task 005 - wait for 1.0 second")
	time.Sleep(1 * time.Second)
	return nil
}

func MainThreadSignal(myapp *app.AppInstance, done chan struct{}) error {
	myapp.Logger.Debug("Main thread: waiting for 5 seconds to complete")
	time.Sleep(5 * time.Second)
	myapp.Logger.Debug("Main thread: signalling we are done")
	done <- app.DONE
	myapp.Logger.Debug("Main thread: signaled we are done")
	return nil
}

func Task001Signal(myapp *app.AppInstance, done chan struct{}) error {
	myapp.Logger.Debug("Task 001 thread: waiting for done")
	_ = <-done
	myapp.Logger.Debug("Task 001 thread: got done")
	return nil
}

func Task002Signal(myapp *app.AppInstance, done chan struct{}) error {
	myapp.Logger.Debug("Task 002 thread: waiting for done")
	_ = <-done
	myapp.Logger.Debug("Task 002 thread: got done")
	return nil
}

func Task003Signal(myapp *app.AppInstance, done chan struct{}) error {
	myapp.Logger.Debug("Task 003 thread: waiting for done")
	_ = <-done
	myapp.Logger.Debug("Task 003 thread: got done, waiting for another 3 seconds")
	time.Sleep(3 * time.Second)
	myapp.Logger.Debug("Task 003 thread: now finished")
	return nil
}

func Task004Signal(myapp *app.AppInstance, done chan struct{}) error {
	myapp.Logger.Debug("Task 004 thread: doing work occasionally")
	t := time.NewTicker(200 * time.Millisecond)
OUTER_LOOP:
	for {
		select {
		case <-t.C:
			myapp.Logger.Debug("Tick 004")
		case <-done:
			myapp.Logger.Debug("Done 004")
			t.Stop()
			break OUTER_LOOP
		}
	}
	myapp.Logger.Debug("Task 004 thread: now finished")
	return nil
}

func Task005Signal(myapp *app.AppInstance, done chan struct{}) error {
	myapp.Logger.Debug("Task 005 thread: doing work occasionally")
	t := time.NewTicker(1000 * time.Millisecond)
OUTER_LOOP:
	for {
		select {
		case <-t.C:
			myapp.Logger.Debug("Tick 005")
		case <-done:
			myapp.Logger.Debug("Done 005")
			t.Stop()
			break OUTER_LOOP
		}
	}
	myapp.Logger.Debug("Task 005 thread: now finished")
	return nil
}
