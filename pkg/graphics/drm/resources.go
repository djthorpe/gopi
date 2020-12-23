// +build drm

package drm

import (
	"sync"

	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Resources struct {
	sync.RWMutex

	fd  uintptr
	res *drm.ModeResources
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewResources(fd uintptr) (*Resources, error) {
	this := new(Resources)
	this.fd = fd
	if res, err := drm.GetResources(this.fd); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *Resources) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.res != nil {
		this.res.Free()
	}

	// Release resources
	this.res = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return active connectors with specific name or all active
// connectors if no name is provided
func (this *Resources) ActiveConnectors(name string) []*Connector {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.res == nil {
		return nil
	}

	for _, id := range this.res.Connectors() {
		if conn, err := this.res.GetConnector(this.fd, id); err != nil {
			continue
		} else if conn.Status() != ModeConnectionConnected {
			conn.Free()
			continue
		} else {
			result = append(result, NewConnector(conn))
		}
	}

	// Return connectors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY
