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

type Connector struct {
	sync.RWMutex
	Properties

	fd  uintptr
	ctx *drm.ModeConnector
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewConnector(fd uintptr, ctx *drm.ModeConnector) (*Connector, error) {
	this := new(Connector)
	if ctx == nil || fd == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewConnector")
	} else {
		this.fd = fd
		this.ctx = ctx
	}

	if err := this.Properties.New(fd, ctx.Id()); err != nil {
		return nil, err
	} else {
		return this, nil
	}
}

func (this *Connector) Dispose() error {
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

func (this *Connector) P() *Properties {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return &this.Properties
}

func (this *Connector) Encoder() uint32 {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return 0
	} else {
		return this.ctx.Encoder()
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Modes returns all modes for the connector if the name and vrefresh are empty
// or returns all modes matching a specific name. If preferred argument is true
// will only return modes which have the that flag set
func (this *Connector) Modes(name string, vrefresh uint32, preferred bool) []*Mode {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return nil
	}

	result := []*Mode{}
	for _, mode := range this.ctx.Modes() {
		if name != "" && name != mode.Name() {
			continue
		} else if vrefresh != 0 && vrefresh != mode.VRefresh() {
			continue
		} else if preferred && mode.Type()&drm.DRM_MODE_TYPE_PREFERRED == 0 {
			continue
		} else {
			result = append(result, NewMode(mode))
		}
	}

	return result
}

// PreferredMode returns a single mode that matches the name and v refresh
// values (or any if they are zero-values or nil if not available
func (this *Connector) PreferredMode(name string, vrefresh uint32) *Mode {
	modes := this.Modes(name, vrefresh, false)
	// Attempt for preferred mode first
	for _, mode := range modes {
		if mode.Type()&drm.DRM_MODE_TYPE_PREFERRED != 0 {
			return mode
		}
	}
	// Use first mode otherwise
	if len(modes) > 0 {
		return modes[0]
	}
	// Return not found
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Connector) String() string {
	str := "<drm.connector"
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	str += " props=" + fmt.Sprint(this.P())
	return str + ">"
}
