// +build pulse

package pulseaudio

import (
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex

	// TODO out []*outputctx
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Manager) New(gopi.Config) error {

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<mediamanager.pulseaudio"
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION
