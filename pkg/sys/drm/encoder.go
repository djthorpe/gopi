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
	ModeEncoder C.drmModeEncoder
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func GetEncoder(fd uintptr, id uint32) (*ModeEncoder, error) {
	if enc := C.drmModeGetEncoder(C.int(fd), C.uint32_t(id)); enc == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("GetEncoder")
	} else {
		return (*ModeEncoder)(unsafe.Pointer(enc)), nil
	}
}

func (this *ModeEncoder) Free() {
	ctx := (*C.drmModeEncoder)(this)
	C.drmModeFreeEncoder(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *ModeEncoder) Id() uint32 {
	ctx := (*C.drmModeEncoder)(this)
	return uint32(ctx.encoder_id)
}

func (this *ModeEncoder) Crtc() uint32 {
	ctx := (*C.drmModeEncoder)(this)
	return uint32(ctx.crtc_id)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ModeEncoder) String() string {
	str := "<drm.encoder"
	str += " id=" + fmt.Sprint(this.Id())
	if crtc := this.Crtc(); crtc != 0 {
		str += " crtc=" + fmt.Sprint(crtc)
	}
	return str + ">"
}
