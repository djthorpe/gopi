package event_test

import (
	"errors"
	"testing"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ErrNumber3 = errors.New("Error Number 3")
	ErrNumber6 = errors.New("Error Number 6")
)

////////////////////////////////////////////////////////////////////////////////
// CREATE TASKS OBJECT

func TestTasks_000(t *testing.T) {
	tasks := &event.Tasks{}
	defer tasks.Close()

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_000(t, start, stop)
	})
}

func TestTasks_001(t *testing.T) {
	tasks := &event.Tasks{}
	defer tasks.Close()

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_000(t, start, stop)
	}, func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_001(t, start, stop)
	})
}

func TestTasks_002(t *testing.T) {
	tasks := &event.Tasks{}

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_001(t, start, stop)
	}, func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_002(t, start, stop)
	})

	tasks.Close()
}

func TestTasks_003(t *testing.T) {
	tasks := &event.Tasks{}

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_003(t, start, stop)
	})

	err := tasks.Close()
	if err != ErrNumber3 {
		t.Error("Expected ErrNumber003, got", err)
	}
}
func TestTasks_004(t *testing.T) {
	tasks := &event.Tasks{}

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_003(t, start, stop)
	}, func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_004(t, start, stop)
	})

	err := tasks.Close()
	if err != ErrNumber3 {
		t.Error("Expected ErrNumber003, got", err)
	}
}

func TestTasks_005(t *testing.T) {
	tasks := &event.Tasks{}

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_003(t, start, stop)
	}, func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_004(t, start, stop)
	}, func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_005(t, start, stop)
	})

	err := tasks.Close()
	if err != ErrNumber3 {
		t.Error("Expected ErrNumber003, got", err)
	}
}

func TestTasks_006(t *testing.T) {
	tasks := &event.Tasks{}

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_004(t, start, stop)
	})
	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_005(t, start, stop)
	})
	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_006(t, start, stop)
	})

	err := tasks.Close()
	if err != ErrNumber6 {
		t.Error("Expected error number 6, got", err)
	}
}

func TestTasks_007(t *testing.T) {
	tasks := &event.Tasks{}

	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_003(t, start, stop)
	}, func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_004(t, start, stop)
	})
	tasks.Start(func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_005(t, start, stop)
	}, func(start chan<- event.Signal, stop <-chan event.Signal) error {
		return task_006(t, start, stop)
	})

	err := tasks.Close()
	if err == nil {
		t.Error("Expected compound error")
	}
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

func task_000(t *testing.T, start chan<- event.Signal, stop <-chan event.Signal) error {
	// This task just returns nil
	t.Log("Entered task_000")
	return nil
}

func task_001(t *testing.T, start chan<- event.Signal, stop <-chan event.Signal) error {
	// This task returns 'start' after 1 second
	t.Log("Entered task_001, sleeping for 1 second")
	time.Sleep(time.Second)
	start <- event.DONE
	t.Log("Woke up task_001, returning nil")
	return nil
}

func task_002(t *testing.T, start chan<- event.Signal, stop <-chan event.Signal) error {
	// This task returns once 'stop' signal is sent
	t.Log("Entered task_002, sleeping for 1 second")
	time.Sleep(time.Second)
	start <- event.DONE
	t.Log("Woke up task_002, waiting for stop signal")
	<-stop
	t.Log("Stop signal received for task_002")
	return nil
}

func task_003(t *testing.T, start chan<- event.Signal, stop <-chan event.Signal) error {
	// This task returns an error straight away
	t.Log("Entered task_003, returning error 3")
	return ErrNumber3
}

func task_004(t *testing.T, start chan<- event.Signal, stop <-chan event.Signal) error {
	// This task returns success immediately after sending start
	t.Log("Entered task_004, sending immediate start")
	start <- event.DONE
	return nil
}

func task_005(t *testing.T, start chan<- event.Signal, stop <-chan event.Signal) error {
	// This task returns success immediately after sending start and receiving stop
	t.Log("Entered task_005, sending immediate start")
	start <- event.DONE
	t.Log("Waiting for stop in task_005")
	<-stop
	t.Log("task_005 got stop, returning nil")
	return nil
}

func task_006(t *testing.T, start chan<- event.Signal, stop <-chan event.Signal) error {
	// This task waits for stop and sends back an error
	t.Log("Entered task_006, sending start and waiting for stop")
	start <- event.DONE
	<-stop
	t.Log("task_006 got stop, returning ErrNumber6")
	return ErrNumber6
}
