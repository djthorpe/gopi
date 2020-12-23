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

	fd    uintptr
	ctx   *drm.Plane
	props map[string]uint64
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
	this.props = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Plane) Type() PlaneType {
	if v, exists := this.GetProperty("type"); exists {
		return PlaneType(v)
	} else {
		return DRM_PLANE_TYPE_NONE
	}
}

func (this *Plane) GetProperty(name string) (uint64, bool) {
	if this.props == nil {
		this.RWMutex.Lock()
		this.props = this.getProperties()
		this.RWMutex.Unlock()
	}
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.fd == 0 || this.ctx == nil {
		return 0, false
	}
	value, exists := this.props[name]
	return value, exists
}

// MatchesCrtc returns true when this plane can be rendered
// by the crtc in the argument
func (this *Plane) MatchesCrtc(crtc *Crtc) bool {
	if crtc == nil || this.ctx == nil {
		return false
	} else {
		return (uint32(1)<<crtc.index)&this.ctx.PossibleCrtcs() != 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Plane) String() string {
	str := "<drm.plane"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	if t := this.Type(); t != DRM_PLANE_TYPE_NONE {
		str += " type=" + fmt.Sprint(t)
	}
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Plane) getProperties() map[string]uint64 {
	if this.fd == 0 || this.ctx == nil {
		return nil
	}

	props := drm.GetPlaneProperties(this.fd, this.ctx.Id())
	if props == nil {
		return nil
	}
	defer props.Free()

	result := make(map[string]uint64)
	values := props.Values()
	for i, key := range props.Keys() {
		prop := drm.NewProperty(this.fd, key)
		if prop == nil {
			continue
		}
		defer prop.Free()
		name := prop.Name()
		result[name] = values[i]
	}
	return result
}
