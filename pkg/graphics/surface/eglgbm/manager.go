// +build egl,gbm,drm

package surface

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/graphics/drm"
	gbmegl "github.com/djthorpe/gopi/v3/pkg/graphics/gbmegl"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	sync.RWMutex

	drm *drm.DRM
	gbm *gbmegl.GBM
	egl *gbmegl.EGL

	mode     *string
	vrefresh *uint
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) Define(cfg gopi.Config) error {
	this.mode = cfg.FlagString("graphics.mode", "", "Set specific graphics mode")
	this.vrefresh = cfg.FlagUint("graphics.vrefresh", 0, "Set specific refresh rate")
	return nil
}

func (this *Manager) New(gopi.Config) error {
	if drm, err := drm.NewDRM(*this.mode, uint32(*this.vrefresh)); err != nil {
		return err
	} else {
		this.drm = drm
	}
	if gbm, err := gbmegl.NewGBM(this.drm.Fd()); err != nil {
		return err
	} else {
		this.gbm = gbm
	}
	if egl, err := gbmegl.NewEGL(this.gbm); err != nil {
		return err
	} else {
		this.egl = egl
	}

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	if this.egl != nil {
		if err := this.egl.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.gbm != nil {
		if err := this.gbm.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.drm != nil {
		if err := this.drm.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.gbm = nil
	this.egl = nil
	this.drm = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) CreateBackground(display gopi.Display, flags gopi.SurfaceFlags) (gopi.Surface, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Manager) DisposeSurface(surface gopi.Surface) error {
	return gopi.ErrNotImplemented
}

func (this *Manager) SwapBuffers() error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<surfacemanager"
	if this.egl != nil {
		str += " egl=" + fmt.Sprint(this.egl)
	}
	if this.gbm != nil {
		str += " gbm=" + fmt.Sprint(this.gbm)
	}
	if this.drm != nil {
		str += " drm=" + fmt.Sprint(this.drm)
	}
	return str + ">"
}
