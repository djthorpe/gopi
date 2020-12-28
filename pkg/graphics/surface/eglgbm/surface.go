// +build egl,gbm,drm

package surface

import (
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	gbmegl "github.com/djthorpe/gopi/v3/pkg/graphics/internal/gbmegl"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surface struct {
	sync.RWMutex
	gopi.Logger

	ctx  *gbmegl.Surface
	x, y uint32
	w, h uint32
}

func NewSurface(ctx *gbmegl.Surface, w, h uint32) *Surface {
	this := new(Surface)
	if ctx == nil {
		return nil
	} else {
		this.ctx = ctx
		this.w, this.h = w, h
	}
	return this
}
