// +build drm

package drm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type DRM struct {
	sync.RWMutex

	fh        *os.File
	connector *Connector
	mode      *Mode
	encoder   *Encoder
	crtc      *Crtc
	primary   *Plane
	cursor    *Plane
	overlay   []*Plane
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Return a DRM object with primary node
func NewDRM(name string, vrefresh uint32) (*DRM, error) {
	this := new(DRM)
	name = strings.TrimSpace(name)

	// Open primary device
	if fh, err := OpenPrimaryDevice(); err != nil {
		return nil, err
	} else {
		this.fh = fh
	}

	// Get Resources
	res, err := NewResources(this.fh.Fd())
	if err != nil {
		this.fh.Close()
		return nil, err
	}

	// Get preferred connector using name and vrefresh
	connectors, err := res.NewActiveConnectorsForMode(name, vrefresh)
	if err != nil {
		res.Dispose()
		this.fh.Close()
		return nil, err
	}
	for _, connector := range connectors {
		if mode := connector.PreferredMode(name, vrefresh); mode != nil {
			this.connector = connector
			this.mode = mode
		} else {
			connector.Dispose()
		}
	}

	// If no connector and/or mode here. bail
	if this.connector == nil || this.mode == nil {
		res.Dispose()
		this.fh.Close()
		return nil, gopi.ErrNotFound.WithPrefix("NewDRM: ", name)
	}

	// Set encoder and Crtc
	if encoder, err := res.NewEncoderForConnector(this.connector); err != nil {
		return nil, err
	} else {
		this.encoder = encoder
	}
	if crtc, err := res.NewCrtcForEncoder(this.encoder); err != nil {
		return nil, err
	} else {
		this.crtc = crtc
	}

	// Dispose of resources
	if err := res.Dispose(); err != nil {
		this.fh.Close()
		return nil, err
	}

	// Set atomic capability
	if err := drm.SetClientCap(this.fh.Fd(), drm.DRM_CLIENT_CAP_ATOMIC, 1); err != nil {
		this.fh.Close()
		return nil, fmt.Errorf("%w: DRM_CLIENT_CAP_ATOMIC", err)
	}

	// Success
	return this, nil
}

func (this *DRM) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	for _, plane := range this.overlay {
		if err := plane.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.cursor != nil {
		if err := this.cursor.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.primary != nil {
		if err := this.primary.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.crtc != nil {
		if err := this.crtc.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.encoder != nil {
		if err := this.encoder.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.mode != nil {
		if err := this.mode.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.connector != nil {
		if err := this.connector.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.fh != nil {
		if err := this.fh.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.fh = nil
	this.connector = nil
	this.mode = nil
	this.encoder = nil
	this.crtc = nil
	this.primary = nil
	this.cursor = nil
	this.overlay = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *DRM) Fd() uintptr {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.fh == nil {
		return 0
	} else {
		return this.fh.Fd()
	}
}

func (this *DRM) Connector() *Connector {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.connector
}

func (this *DRM) Mode() *Mode {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.mode
}

func (this *DRM) Crtc() *Crtc {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.crtc
}

func (this *DRM) Encoder() *Encoder {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.encoder
}

func (this *DRM) Primary() *Plane {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.primary == nil {
		if this.crtc == nil {
			return nil
		} else if planes := this.NewPlanesForCrtc(DRM_PLANE_TYPE_PRIMARY, this.crtc, 1); len(planes) == 0 {
			return nil
		} else {
			this.primary = planes[0]
		}
	}
	return this.primary
}

func (this *DRM) Cursor() *Plane {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.cursor == nil {
		if this.crtc == nil {
			return nil
		} else if planes := this.NewPlanesForCrtc(DRM_PLANE_TYPE_CURSOR, this.crtc, 1); len(planes) == 0 {
			return nil
		} else {
			this.cursor = planes[0]
		}
	}
	return this.cursor
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// NewPlanesForCrtc returns all planes which can be rendered by Crtc with
// appropriate type. Can be limited to first "count" planes.
func (this *DRM) NewPlanesForCrtc(t PlaneType, crtc *Crtc, count int) []*Plane {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.fh == nil {
		return nil
	}

	planes := drm.Planes(this.fh.Fd())
	if len(planes) == 0 {
		return nil
	}

	var result []*Plane
	for index, plane := range planes {
		if ctx, err := drm.GetPlane(this.fh.Fd(), plane); err != nil {
			continue
		} else if plane, err := NewPlane(this.fh.Fd(), ctx, index); err != nil {
			ctx.Free()
			continue
		} else if t != DRM_PLANE_TYPE_NONE && t != plane.Type() {
			plane.Dispose()
			continue
		} else if crtc != nil && plane.MatchesCrtc(crtc) == false {
			plane.Dispose()
			continue
		} else {
			result = append(result, plane)
		}
		if count != 0 && len(result) >= count {
			break
		}
	}

	return result
}

func (this *DRM) CommitChanges(modeset bool) error {
	var result error

	flags := uint32(drm.DRM_MODE_ATOMIC_NONBLOCK)
	if modeset {
		flags |= drm.DRM_MODE_ATOMIC_ALLOW_MODESET
	}

	// Create atomic request
	req := drm.NewAtomic()

	// Get properties for connector, crtc and planes
	if err := setPropertiesFor(req, this.Connector().P()); err != nil {
		result = multierror.Append(result, err)
	}
	if err := setPropertiesFor(req, this.Crtc().P()); err != nil {
		result = multierror.Append(result, err)
	}
	if err := setPropertiesFor(req, this.Primary().P()); err != nil {
		result = multierror.Append(result, err)
	}
	if err := setPropertiesFor(req, this.Cursor().P()); err != nil {
		result = multierror.Append(result, err)
	}
	// TODO: Overlay planes

	// Commit changes atomically
	if err := req.Commit(this.Fd(), flags, 0); err != nil {
		result = multierror.Append(result, err)
	}

	// Free atomic request
	req.Free()

	// Return any errors
	return result
}

func setPropertiesFor(req *drm.Atomic, props *Properties) error {
	var result error
	for k, v := range props.GetDirtyProperties() {
		if err := req.SetObjectProperty(props.Id(), k, v); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *DRM) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<drm"
	if this.fh != nil {
		str += " node=" + strconv.Quote(this.fh.Name())
	}
	if this.connector != nil {
		str += " connector=" + fmt.Sprint(this.connector)
	}
	if this.mode != nil {
		str += " mode=" + fmt.Sprint(this.mode)
	}
	if this.encoder != nil {
		str += " encoder=" + fmt.Sprint(this.encoder)
	}
	if this.crtc != nil {
		str += " crtc=" + fmt.Sprint(this.crtc)
	}
	return str + ">"
}
