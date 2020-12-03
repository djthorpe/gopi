// +build drm

package drm

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
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
	ModeResources C.drmModeRes
	ModeConnector C.drmModeConnector
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DRM_DEV = "/dev/dri/card"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - OPEN

func Device(bus uint) string {
	return fmt.Sprintf("%v%v", DRM_DEV, bus)
}

func Devices() []uint {
	matches, err := filepath.Glob(DRM_DEV + "*")
	if err != nil {
		return nil
	}
	var result []uint
	for _, match := range matches {
		if bus, err := strconv.ParseUint(strings.TrimPrefix(match, DRM_DEV), 10, 32); err == nil {
			result = append(result, uint(bus))
		}
	}
	return result
}

func OpenDevice(bus uint) (*os.File, error) {
	if file, err := os.OpenFile(Device(bus), os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func GetResources(fd uintptr) (*ModeResources, error) {
	if res := C.drmModeGetResources(C.int(fd)); res == nil {
		return nil, gopi.ErrInternalAppError
	} else {
		return (*ModeResources)(unsafe.Pointer(res)), nil
	}
}

func GetConnector(fd uintptr, id uint32) (*ModeConnector, error) {
	if conn := C.drmModeGetConnector(C.int(fd), C.uint32_t(id)); conn == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("GetConnector")
	} else {
		return (*ModeConnector)(unsafe.Pointer(conn)), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// ModeConnector

func (this *ModeConnector) Free() {
	ctx := (*C.drmModeConnector)(this)
	C.drmModeFreeConnector(ctx)
}

func (this *ModeConnector) Id() uint32 {
	ctx := (*C.drmModeConnector)(this)
	return uint32(ctx.connector_id)
}

func (this *ModeConnector) String() string {
	str := "<drm.modeconnector"
	str += " id=" + fmt.Sprint(this.Id())
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// ModeResources

func (this *ModeResources) Free() {
	ctx := (*C.drmModeRes)(this)
	C.drmModeFreeResources(ctx)
}

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

func (this *ModeResources) String() string {
	str := "<drm.moderesources"
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
