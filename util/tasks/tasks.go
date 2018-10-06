package tasks

import (
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/util/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tasks struct {
	sync.WaitGroup
	tasks []*task
	err   errors.CompoundError
}

type TaskFunc func(start chan<- struct{}, stop <-chan struct{}) error

type task struct {
	start, stop chan struct{}
}

var DONE = struct{}{}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Start tasks in the background and waits for all "start" signals to be
// returned before unblocking
func (this *Tasks) Start(funcs ...TaskFunc) {
	// Create any instance variables
	this.new()

	var wg sync.WaitGroup // wg blocks until all start signals are in
	wg.Add(len(funcs))
	this.Add(len(funcs)) // this blocks until all stop signals are in
	for _, fn := range funcs {
		t := &task{make(chan struct{}), make(chan struct{})}
		// Append onto the list of tasks
		this.tasks = append(this.tasks, t)
		// Run the task, then mark as done
		go func(fn TaskFunc, t *task) {
			if err := this.run(fn, t); err != nil {
				this.err.Add(err)
			}
			this.Done()
		}(fn, t)
		// In the background, wait for start signal, or nil if the task
		// ends without sending a start signal
		go func(t *task) {
			select {
			case <-t.start:
				wg.Done() // Indicate the task has started

			}
		}(t)
	}
	// Wait for all started
	wg.Wait()
}

// Close sends done signals to each go routine and will return any error
// from the tasks. Each go routine may end before the 'stop' signal is returned...
func (this *Tasks) Close() error {
	if len(this.tasks) > 0 {
		// Signal all functions to stop
		for _, t := range this.tasks {
			t.stop <- DONE
		}
		// Close all stop channels
		for _, t := range this.tasks {
			close(t.stop)
		}
		// Wait for all tasks to complete
		this.Wait()
	}
	// clear tasks & errors
	this.tasks = nil
	if this.err.Success() {
		return nil
	} else if err := this.err.One(); err != nil {
		return err
	} else {
		return &this.err
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Tasks) new() {
	if this.tasks == nil {
		this.tasks = make([]*task, 0)
	}
}

func (this *Tasks) run(fn TaskFunc, t *task) error {

	// Start the function and wait for error return
	err := fn(t.start, t.stop)

	// Close the start channel
	close(t.start)

	// Receive stop signal in background - this
	// gets the 'nil' on close of the channel
	go func() { <-t.stop }()

	// return any errors
	return err
}
