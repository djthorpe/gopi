// +build drm

package drm

import (
	"fmt"
	"sync"

	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Plane struct {
	sync.RWMutex

	fd  uintptr
	ctx *drm.Plane
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewPlane(fd uintptr, ctx *drm.Plane) *Plane {
	this := new(Plane)
	if ctx == nil || fd == 0 {
		return nil
	}
	this.fd = fd
	this.ctx = ctx
	return this
}

func (this *Plane) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil
	this.fd = 0

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Plane) String() string {
	str := "<drm.plane"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	return str + ">"
}
