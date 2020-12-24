// +build drm

package drm

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Plane struct {
	sync.RWMutex
	Properties

	fd  uintptr
	ctx *drm.Plane
}

type PlaneType uint64

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DRM_PLANE_TYPE_OVERLAY PlaneType = drm.DRM_PLANE_TYPE_OVERLAY
	DRM_PLANE_TYPE_PRIMARY PlaneType = drm.DRM_PLANE_TYPE_PRIMARY
	DRM_PLANE_TYPE_CURSOR  PlaneType = drm.DRM_PLANE_TYPE_CURSOR
	DRM_PLANE_TYPE_NONE    PlaneType = DRM_PLANE_TYPE_CURSOR + 1
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewPlane(fd uintptr, ctx *drm.Plane) (*Plane, error) {
	this := new(Plane)
	if ctx == nil || fd == 0 {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewPlane")
	}
	this.fd = fd
	this.ctx = ctx

	if err := this.Properties.New(fd, ctx.Id()); err != nil {
		return nil, err
	} else {
		return this, nil
	}
}

func (this *Plane) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.ctx != nil {
		this.ctx.Free()
	}

	var result error
	if err := this.Properties.Dispose(); err != nil {
		result = multierror.Append(result, err)
	}

	// Release resources
	this.ctx = nil
	this.fd = 0

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Plane) Id() uint32 {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return 0
	} else {
		return this.ctx.Id()
	}
}

func (this *Plane) Type() PlaneType {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if v, exists := this.GetProperty("type"); exists {
		return PlaneType(v)
	} else {
		return DRM_PLANE_TYPE_NONE
	}
}

// MatchesCrtc returns true when this plane can be rendered
// by the crtc in the argument
func (this *Plane) MatchesCrtc(crtc *Crtc) bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if crtc == nil || this.ctx == nil {
		return false
	} else {
		return (uint32(1)<<crtc.index)&this.ctx.PossibleCrtcs() != 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Plane) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<drm.plane"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	if t := this.Type(); t != DRM_PLANE_TYPE_NONE {
		str += " type=" + fmt.Sprint(t)
	}
	str += " props=" + fmt.Sprint(&this.Properties)
	return str + ">"
}

func (t PlaneType) String() string {
	switch t {
	case DRM_PLANE_TYPE_NONE:
		return "DRM_PLANE_TYPE_NONE"
	case DRM_PLANE_TYPE_OVERLAY:
		return "DRM_PLANE_TYPE_OVERLAY"
	case DRM_PLANE_TYPE_PRIMARY:
		return "DRM_PLANE_TYPE_PRIMARY"
	case DRM_PLANE_TYPE_CURSOR:
		return "DRM_PLANE_TYPE_CURSOR"
	default:
		return "[?? Invalid PlaneType value]"
	}
}
