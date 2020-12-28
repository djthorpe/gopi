// +build drm

package drm

import (
	"fmt"
	"sync"

	"github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Crtc struct {
	sync.RWMutex
	Properties

	fd    uintptr
	ctx   *drm.ModeCRTC
	index int
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCrtc(fd uintptr, ctx *drm.ModeCRTC, index int) (*Crtc, error) {
	this := new(Crtc)
	if ctx == nil || fd == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewCrtc")
	} else {
		this.fd = fd
		this.ctx = ctx
		this.index = index
	}

	if err := this.Properties.New(fd, ctx.Id()); err != nil {
		return nil, err
	} else {
		return this, nil
	}
}

func (this *Crtc) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Release properties
	if err := this.Properties.Dispose(); err != nil {
		result = multierror.Append(result, err)
	}

	// Context
	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil
	this.fd = 0

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Crtc) P() *Properties {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return &this.Properties
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Crtc) String() string {
	str := "<drm.crtc"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	str += " index=" + fmt.Sprint(this.index)
	str += " props=" + fmt.Sprint(this.P())
	return str + ">"
}
