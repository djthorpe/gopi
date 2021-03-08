// +build egl,dispmanx

package surface

import (
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger
	gopi.Platform

	display *uint
	handle  dx.Display
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) Define(cfg gopi.Config) error {
	this.display = cfg.FlagUint("display", 0, "Graphics Display Number")
	return nil
}

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger, this.Platform)

	// Open display
	if handle, err := dx.DisplayOpen(uint32(*this.display)); err != nil {
		return err
	} else {
		this.handle = handle
	}

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	if err := dx.DisplayClose(this.handle); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Manager) CreateBackground(gopi.SurfaceFlags) (gopi.Surface, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Manager) DisposeSurface(gopi.Surface) error {
	return gopi.ErrNotImplemented
}

func (this *Manager) CreateBitmap(gopi.SurfaceFormat, gopi.Size) (gopi.Bitmap, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Manager) DisposeBitmap(gopi.Bitmap) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<surfacemanager"
	return str + ">"
}
