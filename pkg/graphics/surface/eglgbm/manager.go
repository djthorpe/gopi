// +build egl,gbm,drm

package surface

import (
	"context"
	"fmt"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/graphics/internal/drm"
	gbmegl "github.com/djthorpe/gopi/v3/pkg/graphics/internal/gbmegl"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger
	gopi.Metrics

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
	if this.Metrics == nil {
		return gopi.ErrInternalAppError.WithPrefix("Metrics")
	}

	// Create EGL, GBM and EGL interfaces
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

	// Record framerate
	if _, err := this.Metrics.NewMeasurement("vrefresh", "hertz float64"); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Dispose of all interfaces
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

func (this *Manager) Run(ctx context.Context) error {
	var modeset bool

	// Ensure connector and Crtc have been set
	if this.drm.Connector() == nil {
		return gopi.ErrBadParameter.WithPrefix("Connector")
	} else if this.drm.Crtc() == nil {
		return gopi.ErrBadParameter.WithPrefix("Crtc")
	} else {
		if exists := this.drm.Connector().SetProperty("CRTC_ID", uint64(this.drm.Crtc().Id())); exists == false {
			this.Debug("No CRTC_ID property for connector")
		}
		if exists := this.drm.Crtc().SetProperty("ACTIVE", 1); exists == false {
			this.Debug("No ACTIVE property for crtc")
		}

		// TODO: Create Mode Blob and set in Crtc

		this.Debug(this.drm.Crtc())
		this.Debug(this.drm.Connector())

		// Flag to set the modeset on the first run of the loop
		modeset = true
	}

FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		default:
			now := time.Now()

			// Swap EGL buffers for drawing
			if err := this.egl.SwapBuffers(); err != nil {
				this.Print("SwapBuffers: ", err)
			}

			// EGL DRAW HERE

			// DO STUFF IN HERE
			time.Sleep(500 * time.Millisecond)

			// Commit changes to flip page
			if err := this.drm.CommitChanges(modeset); err != nil {
				this.Print("CommitChanges: ", err)
			}

			// Don't set modeset a second time
			modeset = false

			// Send framerate metrics
			framerate := float64(1.0e9) / float64(time.Since(now).Nanoseconds())
			if err := this.Metrics.Emit("vrefresh", framerate); err != nil {
				this.Print(err)
			}
		}
	}
	return ctx.Err()
}

/*
func (this *Manager) RunOnce() {

		if (gbm->surface) {
			next_bo = gbm_surface_lock_front_buffer(gbm->surface);
		} else {
			next_bo = gbm->bos[frame % NUM_BUFFERS];
		}
		if (!next_bo) {
			printf("Failed to lock frontbuffer\n");
			return -1;
		}
		fb = drm_fb_get_from_bo(next_bo);
		if (!fb) {
			printf("Failed to get a new framebuffer BO\n");
			return -1;
		}

	// release last buffer to render on again
	if (bo && gbm->surface)
		gbm_surface_release_buffer(gbm->surface, bo);
	bo = next_bo;:
*/

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) CreateBackground(flags gopi.SurfaceFlags) (gopi.Surface, error) {
	api := ""
	version := uint(0)
	switch flags {
	case SURFACE_FLAG_BITMAP, SURFACE_FLAG_OPENGL_ES:
		api = "OpenGL_ES"
		version = 1
	case SURFACE_FLAG_OPENGL:
		api = "OpenGL"
	case SURFACE_FLAG_OPENGL_ES2:
		api = "OpenGL_ES"
		version = 2
	case SURFACE_FLAG_OPENGL_ES3:
		api = "OpenGL_ES"
		version = 3
	case SURFACE_FLAG_OPENVG:
		api = "OpenVG"
	default:
		return gopi.ErrBadParameter.WithPrefix("CreateBackground", flags)
	}
	if ctx, err := this.egl.CreateSurface(api, version, 100, 100, 0); err != nil {
		return nil, err
	} else {
		// TODO
	}
}

func (this *Manager) DisposeSurface(surface gopi.Surface) error {
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
