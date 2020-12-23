// +build drm

package drm

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libdrm
#include <xf86drm.h>
#include <xf86drmMode.h>
*/
import "C"
import (
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	ModeInfo C.drmModeModeInfo
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this ModeInfo) Name() string {
	ctx := (C.drmModeModeInfo)(this)
	return C.GoString(&ctx.name[0])
}

func (this ModeInfo) Size() (uint32, uint32) {
	ctx := (C.drmModeModeInfo)(this)
	return uint32(ctx.hdisplay), uint32(ctx.vdisplay)
}

func (this ModeInfo) VRefresh() uint32 {
	ctx := (C.drmModeModeInfo)(this)
	return uint32(ctx.vrefresh)
}

func (this ModeInfo) Type() ModeInfoType {
	ctx := (C.drmModeModeInfo)(this)
	return ModeInfoType(ctx._type)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this ModeInfo) String() string {
	str := "<drm.info"
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if t := this.Type(); t != DRM_MODE_TYPE_NONE {
		str += fmt.Sprintf(" type=%v", t)
	}
	if w, h := this.Size(); w > 0 && h > 0 {
		str += fmt.Sprintf(" size={ %v,%v }", w, h)
	}
	if vrefresh := this.VRefresh(); vrefresh > 0 {
		str += fmt.Sprintf(" vrefresh=%v", vrefresh)
	}
	return str + ">"
}
