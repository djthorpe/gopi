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
// PUBLIC METHODS

func (this *Manager) OpenDefaultSink() (gopi.AudioContext, error) {

}

func (this *Manager) Close(ctx gopi.AudioContext) error {

}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<pulseaudio.manager"
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION
