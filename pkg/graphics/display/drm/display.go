// +build drm

package display

import (
	"fmt"
	"strconv"
	"sync"

	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Display struct {
	sync.RWMutex

	ctx  *drm.ModeConnector
	mode drm.ModeInfo
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDisplay(ctx *drm.ModeConnector) *Display {
	this := new(Display)

	if ctx == nil {
		return nil
	} else {
		this.ctx = ctx
	}

	// We always choose the zero-indexed mode
	if modes := this.ctx.Modes(); len(modes) == 0 {
		return nil
	} else {
		this.mode = modes[0]
	}

	return this
}

func (this *Display) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Release DRM resources
	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Display) Id() uint32 {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return 0
	} else {
		return this.ctx.Id()
	}
}

func (this *Display) Name() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return ""
	} else {
		return this.mode.Name()
	}
}

func (this *Display) Size() (uint32, uint32) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return 0, 0
	} else {
		return this.mode.Size()
	}
}

func (this *Display) PixelsPerInch() uint32 {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// TODO: Not yet implemented
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Display) String() string {
	str := "<display.drm"
	if id := this.Id(); id != 0 {
		str += " id=" + fmt.Sprint(id)
	}
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if w, h := this.Size(); w != 0 && h != 0 {
		str += fmt.Sprint(" size={", w, ",", h, "}")
	}
	if ppi := this.PixelsPerInch(); ppi != 0 {
		str += " pixels_per_inch=" + fmt.Sprint(ppi)
	}
	if this.ctx != nil {
		str += " ctx=" + fmt.Sprint(this.ctx)
	}
	return str + ">"
}
