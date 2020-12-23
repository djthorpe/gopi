// +build drm

package drm

import (
	"fmt"
	"sync"

	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Connector struct {
	sync.RWMutex

	ctx *drm.ModeConnector
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewConnector(ctx *drm.ModeConnector) *Connector {
	this := new(Connector)
	if ctx == nil {
		return nil
	}
	this.ctx = ctx
	return this
}

func (this *Connector) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Connector) String() string {
	str := "<drm.connector"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprintf(this.ctx)
	}
	return str + ">"
}
