// +build gbm,drm

package gbmegl

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GBM struct {
	sync.RWMutex

	fd  uintptr
	dev *gbm.GBMDevice
}

type Format gbm.GBMBufferFormat

const (
	GBM_BO_FORMAT_NONE     Format = 0
	GBM_BO_FORMAT_XRGB8888 Format = Format(gbm.GBM_BO_FORMAT_XRGB8888)
	GBM_BO_FORMAT_ARGB8888 Format = Format(gbm.GBM_BO_FORMAT_ARGB8888)
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Return a GBM object from filehandle
func NewGBM(fd uintptr) (*GBM, error) {
	this := new(GBM)

	if fd == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewGBM")
	} else {
		this.fd = fd
	}

	if dev, err := gbm.GBMCreateDevice(this.fd); err != nil {
		return nil, err
	} else {
		this.dev = dev
	}

	// Success
	return this, nil
}

func (this *GBM) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	if this.dev != nil {
		this.dev.Free()
	}

	// Release resources
	this.fd = 0
	this.dev = nil

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *GBM) Device() *gbm.GBMDevice {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.dev
}

// Returns r,g,b,a or all zeros if format not supported
func (this *GBM) BitsForFormat(format Format) (uint, uint, uint, uint) {
	switch format {
	case GBM_BO_FORMAT_XRGB8888:
		return 8, 8, 8, 8
	case GBM_BO_FORMAT_ARGB8888:
		return 8, 8, 8, 8
	default:
		return 0, 0, 0, 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *GBM) NewSurface(w, h uint32, format Format) (*gbm.GBMSurface, error) {
	modifiers := []uint64{drm.DRM_FORMAT_MOD_LINEAR}
	if this.dev == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewSurface")
	} else {
		return this.dev.SurfaceCreateWithModifiers(w, h, gbm.GBMBufferFormat(format), modifiers)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRIMGIFY

func (this *GBM) String() string {
	str := "<gbm"
	if this.dev != nil {
		str += " dev=" + fmt.Sprint(this.dev)
	}
	return str + ">"
}
