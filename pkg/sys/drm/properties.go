// +build drm

package drm

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
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
	Properties   C.drmModeObjectProperties
	Property     C.drmModePropertyRes
	PropertyEnum C.struct_drm_mode_property_enum
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DRM_MODE_OBJECT_CRTC      = C.DRM_MODE_OBJECT_CRTC
	DRM_MODE_OBJECT_CONNECTOR = C.DRM_MODE_OBJECT_CONNECTOR
	DRM_MODE_OBJECT_ENCODER   = C.DRM_MODE_OBJECT_ENCODER
	DRM_MODE_OBJECT_MODE      = C.DRM_MODE_OBJECT_MODE
	DRM_MODE_OBJECT_PROPERTY  = C.DRM_MODE_OBJECT_PROPERTY
	DRM_MODE_OBJECT_FB        = C.DRM_MODE_OBJECT_FB
	DRM_MODE_OBJECT_BLOB      = C.DRM_MODE_OBJECT_BLOB
	DRM_MODE_OBJECT_PLANE     = C.DRM_MODE_OBJECT_PLANE
	DRM_MODE_OBJECT_ANY       = C.DRM_MODE_OBJECT_ANY
)

const (
	DRM_MODE_PROP_PENDING       = C.DRM_MODE_PROP_PENDING
	DRM_MODE_PROP_RANGE         = C.DRM_MODE_PROP_RANGE
	DRM_MODE_PROP_IMMUTABLE     = C.DRM_MODE_PROP_IMMUTABLE
	DRM_MODE_PROP_ENUM          = C.DRM_MODE_PROP_ENUM
	DRM_MODE_PROP_BLOB          = C.DRM_MODE_PROP_BLOB
	DRM_MODE_PROP_EXTENDED_TYPE = C.DRM_MODE_PROP_EXTENDED_TYPE
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func GetPlaneProperties(fd uintptr, id uint32) *Properties {
	ctx := C.drmModeObjectGetProperties(C.int(fd), C.uint32_t(id), C.uint32_t(DRM_MODE_OBJECT_PLANE))
	return (*Properties)(unsafe.Pointer(ctx))
}

func GetConnectorProperties(fd uintptr, id uint32) *Properties {
	ctx := C.drmModeObjectGetProperties(C.int(fd), C.uint32_t(id), C.uint32_t(DRM_MODE_OBJECT_CONNECTOR))
	return (*Properties)(unsafe.Pointer(ctx))
}

func GetCrtcProperties(fd uintptr, id uint32) *Properties {
	ctx := C.drmModeObjectGetProperties(C.int(fd), C.uint32_t(id), C.uint32_t(DRM_MODE_OBJECT_CRTC))
	return (*Properties)(unsafe.Pointer(ctx))
}

func GetEncoderProperties(fd uintptr, id uint32) *Properties {
	ctx := C.drmModeObjectGetProperties(C.int(fd), C.uint32_t(id), C.uint32_t(DRM_MODE_OBJECT_ENCODER))
	return (*Properties)(unsafe.Pointer(ctx))
}

func GetFrameBufferProperties(fd uintptr, id uint32) *Properties {
	ctx := C.drmModeObjectGetProperties(C.int(fd), C.uint32_t(id), C.uint32_t(DRM_MODE_OBJECT_FB))
	return (*Properties)(unsafe.Pointer(ctx))
}

func NewProperty(fd uintptr, id uint32) *Property {
	ctx := C.drmModeGetProperty(C.int(fd), C.uint32_t(id))
	return (*Property)(unsafe.Pointer(ctx))
}

func (this *Properties) Free() {
	ctx := (*C.drmModeObjectProperties)(this)
	C.drmModeFreeObjectProperties(ctx)
}

func (this *Property) Free() {
	ctx := (*C.drmModePropertyRes)(this)
	C.drmModeFreeProperty(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Properties) Keys() []uint32 {
	var result []uint32
	ctx := (*C.drmModeObjectProperties)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_props)
	sliceHeader.Len = int(ctx.count_props)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.props))
	return result
}

func (this *Properties) Values() []uint64 {
	var result []uint64
	ctx := (*C.drmModeObjectProperties)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_props)
	sliceHeader.Len = int(ctx.count_props)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.prop_values))
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTY

func (this *Property) Id() uint32 {
	ctx := (*C.drmModePropertyRes)(this)
	return uint32(ctx.prop_id)
}

func (this *Property) Flags() uint32 {
	ctx := (*C.drmModePropertyRes)(this)
	return uint32(ctx.flags)
}

func (this *Property) Name() string {
	ctx := (*C.drmModePropertyRes)(this)
	return C.GoString(&ctx.name[0])
}

func (this *Property) Values() []uint64 {
	var result []uint64
	ctx := (*C.drmModePropertyRes)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_values)
	sliceHeader.Len = int(ctx.count_values)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.values))
	return result
}

func (this *Property) Blobs() []uint32 {
	var result []uint32
	ctx := (*C.drmModePropertyRes)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_blobs)
	sliceHeader.Len = int(ctx.count_blobs)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.blob_ids))
	return result
}

func (this *Property) Enums() []PropertyEnum {
	var result []PropertyEnum
	ctx := (*C.drmModePropertyRes)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_enums)
	sliceHeader.Len = int(ctx.count_enums)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.enums))
	return result
}

////////////////////////////////////////////////////////////////////////////////
// ENUM

func (this PropertyEnum) Value() uint64 {
	ctx := (C.struct_drm_mode_property_enum)(this)
	return uint64(ctx.value)
}

func (this PropertyEnum) Name() string {
	ctx := (C.struct_drm_mode_property_enum)(this)
	return C.GoString(&ctx.name[0])
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Properties) String() string {
	str := "<drm.properties"
	values := this.Values()
	for i, k := range this.Keys() {
		str += fmt.Sprintf(" %d=>%d", k, values[i])
	}
	return str + ">"
}

func (this *Property) String() string {
	str := "<drm.property"
	str += " id=" + fmt.Sprint(this.Id())
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if flags := this.Flags(); flags != 0 {
		str += fmt.Sprintf(" flags=0x%08X", flags)
	}
	if values := this.Values(); len(values) != 0 {
		str += " values=" + fmt.Sprint(values)
	}
	if blobs := this.Blobs(); len(blobs) != 0 {
		str += " blobs=" + fmt.Sprint(blobs)
	}
	if enums := this.Enums(); len(enums) != 0 {
		str += " enums=" + fmt.Sprint(enums)
	}
	return str + ">"
}

func (this PropertyEnum) String() string {
	str := "<"
	str += fmt.Sprint(this.Value(), ":")
	str += strconv.Quote(this.Name())
	return str + ">"
}
