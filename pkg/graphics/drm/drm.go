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

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// NewPlanes returns all planes. They neeed to be disposed of
func (this *DRM) NewPlanes() []*Plane {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.fh == nil {
		return nil
	}

	planes := drm.Planes(this.fh.Fd())
	if planes == nil {
		return nil
	}

	var result []*Plane
	for _, plane := range planes {
		if ctx, err := drm.GetPlane(this.fh.Fd(), plane); err != nil {
			continue
		} else if plane := NewPlane(this.fh.Fd(), ctx); plane == nil {
			continue
		} else {
			result = append(result, plane)
		}
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *DRM) String() string {
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
