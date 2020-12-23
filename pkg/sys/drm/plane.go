// +build drm

package drm

import (
	"fmt"
	"reflect"
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
	Plane C.drmModePlane
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func GetPlane(fd uintptr, id uint32) (*Plane, error) {
	if plane := C.drmModeGetPlane(C.int(fd), C.uint32_t(id)); plane == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("GetPlane")
	} else {
		return (*Plane)(unsafe.Pointer(plane)), nil
	}
}

func (this *Plane) Free() {
	ctx := (*C.drmModePlane)(this)
	C.drmModeFreePlane(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Planes(fd uintptr) []uint32 {
	var result []uint32

	planes := C.drmModeGetPlaneResources(C.int(fd))
	if planes == nil {
		return nil
	}

	// Make fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(planes.count_planes)
	sliceHeader.Len = int(planes.count_planes)
	sliceHeader.Data = uintptr(unsafe.Pointer(planes.planes))

	// Return result
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Plane) Id() uint32 {
	ctx := (*C.drmModePlane)(this)
	return uint32(ctx.plane_id)
}

func (this *Plane) Fb() uint32 {
	ctx := (*C.drmModePlane)(this)
	return uint32(ctx.fb_id)
}

func (this *Plane) XY() (uint32, uint32) {
	ctx := (*C.drmModePlane)(this)
	return uint32(ctx.x), uint32(ctx.y)
}

func (this *Plane) CrtcXY() (uint32, uint32) {
	ctx := (*C.drmModePlane)(this)
	return uint32(ctx.crtc_x), uint32(ctx.crtc_y)
}

/* TODO
func (this *Plane) Formats() []uint32 {

}
*/

func (this *Plane) Crtc() uint32 {
	ctx := (*C.drmModePlane)(this)
	return uint32(ctx.crtc_id)
}

func (this *Plane) PossibleCrtcs() uint32 {
	ctx := (*C.drmModePlane)(this)
	return uint32(ctx.possible_crtcs)
}

func (this *Plane) GammaSize() uint32 {
	ctx := (*C.drmModePlane)(this)
	return uint32(ctx.gamma_size)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Plane) String() string {
	str := "<drm.plane"
	str += " id=" + fmt.Sprint(this.Id())
	if fb := this.Fb(); fb != 0 {
		str += " fb=" + fmt.Sprint(fb)
	}
	if crtc := this.Crtc(); crtc != 0 {
		str += " crtc=" + fmt.Sprint(crtc)
	}
	x, y := this.XY()
	str += " x,y={" + fmt.Sprint(x, ",", y) + "}"
	crtc_x, crtc_y := this.CrtcXY()
	str += " crtc_x,crtc_y={" + fmt.Sprint(crtc_x, ",", crtc_y) + "}"

	if gamma_size := this.GammaSize(); gamma_size != 0 {
		str += " gamma_size=" + fmt.Sprint(gamma_size)
	}
	if possible_crtcs := this.PossibleCrtcs(); possible_crtcs > 0 {
		str += fmt.Sprintf(" possible_crtcs=0b%032b", possible_crtcs)
	}
	return str + ">"
}
