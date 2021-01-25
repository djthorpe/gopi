// +build dvb

package dvb

import (
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) Tuners() []gopi.DVBTuner {
	tuners := []gopi.DVBTuner{}
	for _, device := range dvb.Devices() {
		tuners = append(tuners, NewTuner(device))
	}
	return tuners
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<dvb.mamager"
	return str + ">"
}
