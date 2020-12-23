// +build gbm,egl

package gbmegl

import (
	"fmt"
	"sync"

	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surface struct {
	sync.RWMutex

	gbm *gbm.GBMSurface
	egl egl.EGLSurface
	ctx egl.EGLContext
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSurface(gbm *gbm.GBMSurface, egl egl.EGLSurface, ctx egl.EGLContext) *Surface {
	this := new(Surface)
	this.gbm = gbm
	this.egl = egl
	this.ctx = ctx
	return this
}

func (this *Surface) Dispose(display egl.EGLDisplay) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	if this.egl != nil {
		if err := egl.EGLDestroySurface(display, this.egl); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.ctx != nil {
		if err := egl.EGLDestroyContext(display, this.ctx); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.gbm != nil {
		this.gbm.Free()
	}

	// Release resources
	this.gbm = nil
	this.egl = nil
	this.ctx = nil

	// Return errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Surface) String() string {
	str := "<surface"
	if this.egl != nil {
		str += " egl=" + fmt.Sprint(this.egl)
	}
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	if this.gbm != nil {
		str += " gbm=" + fmt.Sprint(this.gbm)
	}
	return str + ">"
}
