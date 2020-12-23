// +build gbm

package gbm

import (
	"os"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: gbm
#include <gbm.h>

int _errno();
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	GBMSurface C.struct_gbm_surface
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *GBMDevice) SurfaceCreate(width, height uint32, format GBMBufferFormat, flags GBMBufferFlags) (*GBMSurface, error) {
	ctx := (*C.struct_gbm_device)(this)
	if surface := C.gbm_surface_create(ctx, C.uint32_t(width), C.uint32_t(height), C.uint32_t(format), C.uint32_t(flags)); surface != nil {
		return (*GBMSurface)(surface), nil
	} else {
		return nil, os.NewSyscallError("gbm_surface_create", syscall.Errno(C._errno()))
	}
}

func (this *GBMDevice) SurfaceCreateWithModifiers(width, height uint32, format GBMBufferFormat, modifiers []uint64) (*GBMSurface, error) {
	ctx := (*C.struct_gbm_device)(this)
	data := (*C.uint64_t)(nil)
	count := (C.uint)(0)
	if len(modifiers) > 0 {
		data = (*C.uint64_t)(unsafe.Pointer(&modifiers[0]))
	}
	if surface := C.gbm_surface_create_with_modifiers(ctx, C.uint32_t(width), C.uint32_t(height), C.uint32_t(format), data, count); surface != nil {
		return (*GBMSurface)(surface), nil
	} else {
		return nil, os.NewSyscallError("gbm_surface_create", syscall.Errno(C._errno()))
	}
}

func (this *GBMSurface) Free() {
	ctx := (*C.struct_gbm_surface)(this)
	C.gbm_surface_destroy(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *GBMSurface) RetainBuffer() *GBMBuffer {
	ctx := (*C.struct_gbm_surface)(this)
	if buf := C.gbm_surface_lock_front_buffer(ctx); buf == nil {
		return nil
	} else {
		return (*GBMBuffer)(buf)
	}
}

func (this *GBMSurface) ReleaseBuffer(buf *GBMBuffer) {
	ctx := (*C.struct_gbm_surface)(this)
	C.gbm_surface_release_buffer(ctx, (*C.struct_gbm_bo)(buf))
}

func (this *GBMSurface) HasFreeBuffers() bool {
	ctx := (*C.struct_gbm_surface)(this)
	if ret := C.gbm_surface_has_free_buffers(ctx); ret == 0 {
		return false
	} else {
		return true
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GBMSurface) String() string {
	str := "<gbm.surface"
	return str + ">"
}
