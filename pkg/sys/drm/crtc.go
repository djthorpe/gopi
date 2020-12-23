// +build drm

package drm

import (
	"fmt"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libdrm
#include <xf86drm.h>
#include <xf86drmMode.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	ModeCRTC C.drmModeCrtc
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func GetCRTC(fd uintptr, id uint32) (*ModeCRTC, error) {
	if crtc := C.drmModeGetCrtc(C.int(fd), C.uint32_t(id)); crtc == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("GetCRTC")
	} else {
		return (*ModeCRTC)(unsafe.Pointer(crtc)), nil
	}
}

func (this *ModeCRTC) Free() {
	ctx := (*C.drmModeCrtc)(this)
	C.drmModeFreeCrtc(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *ModeCRTC) Id() uint32 {
	ctx := (*C.drmModeCrtc)(this)
	return uint32(ctx.crtc_id)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ModeCRTC) String() string {
	str := "<drm.crtc"
	str += " id=" + fmt.Sprint(this.Id())
	return str + ">"
}
