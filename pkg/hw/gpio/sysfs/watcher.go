package sysfs

import (
	"sync"
	"syscall"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Watcher struct {
	gopi.Unit
	gopi.FilePoll
	gopi.Logger
	sync.Mutex

	pins map[gopi.GPIOPin]state
}

type state struct {
	fd   uintptr
	edge gopi.GPIOEdge
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Watcher) New(gopi.Config) error {
	// Check for filepoll
	if this.FilePoll == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing gopi.FilePoll")
	}

	// Set up state
	this.pins = make(map[gopi.GPIOPin]state)

	// Return success
	return nil
}

func (this *Watcher) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Stop watching edges and close file descriptors
	var result error
	for pin, state := range this.pins {
		if err := writeEdge(pin, "none"); err != nil {
			result = multierror.Append(result, err)
		}
		if err := this.FilePoll.Unwatch(state.fd); err != nil {
			result = multierror.Append(result, err)
		}
		if err := syscall.Close(int(state.fd)); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Watcher) Watch(fd uintptr, pin gopi.GPIOPin, edge gopi.GPIOEdge) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.FilePoll.Watch(fd, gopi.FILEPOLL_FLAG_READ, func(uintptr, gopi.FilePollFlags) {
		if value, err := readPin(pin); err != nil {
			this.Print("Watch: ", pin, ": ", err)
		} else {
			this.Debug("Watch: ", pin, ": ", value)
		}
	})
}

func (this *Watcher) Unwatch(pin gopi.GPIOPin) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.Print("TODO: Unwatch", pin)

	// Return success
	return nil
}
