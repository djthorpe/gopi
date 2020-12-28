// +build drm

package drm

import (
	"fmt"
	"sync"

	"github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Encoder struct {
	sync.RWMutex

	fd  uintptr
	ctx *drm.ModeEncoder
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEncoder(fd uintptr, ctx *drm.ModeEncoder) (*Encoder, error) {
	this := new(Encoder)
	if ctx == nil || fd == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewEncoder")
	} else {
		this.fd = fd
		this.ctx = ctx
	}
	return this, nil
}

func (this *Encoder) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Context
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
// PROPERTIES

func (this *Encoder) Crtc() uint32 {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return 0
	} else {
		return this.ctx.Crtc()
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Encoder) String() string {
	str := "<drm.encoder"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	return str + ">"
}
