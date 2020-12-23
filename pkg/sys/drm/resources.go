// +build drm

package drm

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libdrm
#include <xf86drm.h>
#include <xf86drmMode.h>

extern int _drm_errno();
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	ModeResources C.drmModeRes
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func GetResources(fd uintptr) (*ModeResources, error) {
	if res := C.drmModeGetResources(C.int(fd)); res == nil {
		return nil, os.NewSyscallError("drmModeGetResources", syscall.Errno(C._drm_errno()))
	} else {
		return (*ModeResources)(unsafe.Pointer(res)), nil
	}
}

func (this *ModeResources) Free() {
	ctx := (*C.drmModeRes)(this)
	C.drmModeFreeResources(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *ModeResources) Framebuffers() []uint32 {
	var result []uint32

	// Make fake slice
	ctx := (*C.drmModeRes)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_fbs)
	sliceHeader.Len = int(ctx.count_fbs)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.fbs))

	// Return result
	return result
}

func (this *ModeResources) CRTCs() []uint32 {
	var result []uint32

	// Make fake slice
	ctx := (*C.drmModeRes)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_crtcs)
	sliceHeader.Len = int(ctx.count_crtcs)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.crtcs))

	// Return result
	return result
}

func (this *ModeResources) Connectors() []uint32 {
	var result []uint32

	// Make fake slice
	ctx := (*C.drmModeRes)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_connectors)
	sliceHeader.Len = int(ctx.count_connectors)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.connectors))

	// Return result
	return result
}

func (this *ModeResources) Encoders() []uint32 {
	var result []uint32

	// Make fake slice
	ctx := (*C.drmModeRes)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_encoders)
	sliceHeader.Len = int(ctx.count_encoders)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.encoders))

	// Return result
	return result
}

func (this *ModeResources) Width() (uint32, uint32) {
	ctx := (*C.drmModeRes)(this)
	return uint32(ctx.min_width), uint32(ctx.max_width)
}

func (this *ModeResources) Height() (uint32, uint32) {
	ctx := (*C.drmModeRes)(this)
	return uint32(ctx.min_height), uint32(ctx.max_height)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ModeResources) String() string {
	str := "<drm.resources"
	if fb := this.Framebuffers(); len(fb) > 0 {
		str += " fb=" + fmt.Sprint(fb)
	}
	if crtc := this.CRTCs(); len(crtc) > 0 {
		str += " crtc=" + fmt.Sprint(crtc)
	}
	if connectors := this.Connectors(); len(connectors) > 0 {
		str += " connectors=" + fmt.Sprint(connectors)
	}
	if encoders := this.Encoders(); len(encoders) > 0 {
		str += " encoders=" + fmt.Sprint(encoders)
	}
	if min, max := this.Width(); min > 0 || max > 0 {
		str += fmt.Sprintf(" width{min,max}={%v,%v}", min, max)
	}
	if min, max := this.Height(); min > 0 || max > 0 {
		str += fmt.Sprintf(" height{min,max}={%v,%v}", min, max)
	}
	return str + ">"
}
