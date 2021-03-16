package input

import (
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger
	gopi.FilePoll

	devices map[uintptr]gopi.InputDevice
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger, this.FilePoll)

	// Create devices
	this.devices = make(map[uintptr]gopi.InputDevice)

	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Free devices
	var result error
	this.devices = nil

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) Devices() []gopi.InputDevice {
	return nil
}

func (this *Manager) RegisterDevice(gopi.InputDevice) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS
