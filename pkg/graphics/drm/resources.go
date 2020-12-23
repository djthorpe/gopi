// +build drm

package drm

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
	"github.com/hashicorp/go-multierror"
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
	} else {
		this.res = res
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

// NewActiveConnectors returns active connectors. Connectors returned
// need to be disposed
func (this *Resources) NewActiveConnectors() []*Connector {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.res == nil {
		return nil
	}

	result := []*Connector{}
	for _, id := range this.res.Connectors() {
		if conn, err := drm.GetConnector(this.fd, id); err != nil {
			continue
		} else if conn.Status() != drm.ModeConnectionConnected {
			conn.Free()
			continue
		} else {
			result = append(result, NewConnector(this.fd, conn))
		}
	}

	// Return connectors
	return result
}

// NewActiveConnectorsForMode returns active connectors for a named mode
// which falls back to NewActiveConnectors when the name is not included
func (this *Resources) NewActiveConnectorsForMode(name string, vrefresh uint32) ([]*Connector, error) {
	connectors := this.NewActiveConnectors()
	if name == "" || len(connectors) == 0 {
		return connectors, nil
	}

	// Find connectors with correct mode
	var result []*Connector
	var errs error
	for _, connector := range connectors {
		if modes := connector.Modes(name, vrefresh, false); len(modes) > 0 {
			result = append(result, connector)
		} else if err := connector.Dispose(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	// Return result
	return result, errs
}

// NewEncoderForConnector returns the encoder object
func (this *Resources) NewEncoderForConnector(connector *Connector) (*Encoder, error) {
	if this.res == nil || this.fd == 0 {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewEncoderForConnector")
	} else if connector == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewEncoderForConnector")
	}

	if ctx, err := drm.GetEncoder(this.fd, connector.Encoder()); err != nil {
		return nil, err
	} else if encoder := NewEncoder(this.fd, ctx); encoder == nil {
		ctx.Free()
		return nil, gopi.ErrInternalAppError.WithPrefix("NewEncoderForConnector")
	} else {
		return encoder, nil
	}
}

// NewCrtcForEncoder returns the ctrc object
func (this *Resources) NewCrtcForEncoder(encoder *Encoder) (*Crtc, error) {
	if this.res == nil || this.fd == 0 {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewCrtcForEncoder")
	} else if encoder == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewCrtcForEncoder")
	}

	// Go through all Crtcs
	for index, id := range this.res.CRTCs() {
		if id != encoder.Crtc() {
			continue
		}
		if ctx, err := drm.GetCRTC(this.fd, id); err != nil {
			return nil, err
		} else if crtc := NewCrtc(this.fd, ctx, index); crtc == nil {
			ctx.Free()
			return nil, gopi.ErrInternalAppError.WithPrefix("NewCrtcForEncoder")
		} else {
			return crtc, nil
		}
	}

	// Not found
	return nil, gopi.ErrNotFound.WithPrefix("NewCrtcForEncoder")
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Resources) String() string {
	str := "<drm.resources"
	if this.res != nil {
		str += " res=" + fmt.Sprint(this.res)
	}
	return str + ">"
}
