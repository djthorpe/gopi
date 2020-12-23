// +build drm

package drm

import (
	"fmt"
	"sync"

	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Crtc struct {
	sync.RWMutex

	fd  uintptr
	ctx *drm.ModeCRTC
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCrtc(fd uintptr, ctx *drm.ModeCRTC) *Crtc {
	this := new(Crtc)
	if ctx == nil || fd == 0 {
		return nil
	}
	this.fd = fd
	this.ctx = ctx
	return this
}

func (this *Crtc) Dispose() error {
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

func (this *Crtc) String() string {
	str := "<drm.crtc"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	return str + ">"
}
