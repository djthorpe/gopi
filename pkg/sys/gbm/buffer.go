// +build gbm

package gbm

import (
	"fmt"
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
	GBMBuffer C.struct_gbm_bo
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *GBMDevice) BufferCreate(width, height uint32, format GBMFormat, flags GBMBufferFlags) (*GBMBuffer, error) {
	ctx := (*C.struct_gbm_device)(this)
	if buf := C.gbm_bo_create(ctx, C.uint32_t(width), C.uint32_t(height), C.uint32_t(format), C.uint32_t(flags)); buf != nil {
		return (*GBMBuffer)(buf), nil
	} else {
		return nil, os.NewSyscallError("gbm_bo_create", syscall.Errno(C._errno()))
	}
}

func (this *GBMDevice) BufferImport(foreign_type GBMBufferType, ptr uintptr, flags GBMBufferFlags) (*GBMBuffer, error) {
	ctx := (*C.struct_gbm_device)(this)
	if buf := C.gbm_bo_import(ctx, C.uint32_t(foreign_type), unsafe.Pointer(ptr), C.uint32_t(flags)); buf != nil {
		return (*GBMBuffer)(buf), nil
	} else {
		return nil, os.NewSyscallError("gbm_bo_import", syscall.Errno(C._errno()))
	}
}

func (this *GBMBuffer) Free() {
	ctx := (*C.struct_gbm_bo)(this)
	C.gbm_bo_destroy(ctx)
}

func (this *GBMBuffer) Width() uint32 {
	ctx := (*C.struct_gbm_bo)(this)
	return uint32(C.gbm_bo_get_width(ctx))
}

func (this *GBMBuffer) Height() uint32 {
	ctx := (*C.struct_gbm_bo)(this)
	return uint32(C.gbm_bo_get_height(ctx))
}

func (this *GBMBuffer) Stride() uint32 {
	ctx := (*C.struct_gbm_bo)(this)
	return uint32(C.gbm_bo_get_stride(ctx))
}

func (this *GBMBuffer) PlaneCount() uint {
	ctx := (*C.struct_gbm_bo)(this)
	return uint(C.gbm_bo_get_plane_count(ctx))
}

func (this *GBMBuffer) StrideForPlane(plane uint) uint32 {
	ctx := (*C.struct_gbm_bo)(this)
	return uint32(C.gbm_bo_get_stride_for_plane(ctx, C.int(plane)))
}

func (this *GBMBuffer) Format() GBMFormat {
	ctx := (*C.struct_gbm_bo)(this)
	return GBMFormat(C.gbm_bo_get_format(ctx))
}

func (this *GBMBuffer) BitsPerPixel() uint32 {
	ctx := (*C.struct_gbm_bo)(this)
	return uint32(C.gbm_bo_get_bpp(ctx))
}

func (this *GBMBuffer) Offset(plane uint) uint32 {
	ctx := (*C.struct_gbm_bo)(this)
	return uint32(C.gbm_bo_get_offset(ctx, C.int(plane)))
}

func (this *GBMBuffer) Write(data []byte) error {
	ctx := (*C.struct_gbm_bo)(this)
	ptr := unsafe.Pointer(&data[0])
	size := C.size_t(len(data))
	if ret := C.gbm_bo_write(ctx, ptr, size); ret == 0 {
		return nil
	} else {
		return os.NewSyscallError("gbm_bo_write", syscall.Errno(C._errno()))
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GBMBuffer) String() string {
	str := "<gbm.buffer"
	if w, h := this.Width(), this.Height(); w != 0 || h != 0 {
		str += fmt.Sprint(" size={", w, ",", h, "}")
	}
	if format := this.Format(); format != 0 {
		str += fmt.Sprint(" format=", format)
	}
	if stride := this.Stride(); stride != 0 {
		str += fmt.Sprint(" stride=", stride)
	}
	if bpp := this.BitsPerPixel(); bpp != 0 {
		str += fmt.Sprint(" bits_per_pixel=", bpp)
	}
	if planes := this.PlaneCount(); planes != 0 {
		str += fmt.Sprint(" plane_count=", planes)
		for plane := uint(0); plane < planes; plane++ {
			str += fmt.Sprint(" plane_", plane, "=<stride=", this.StrideForPlane(plane), " offset=", this.Offset(plane), ">")
		}
	}
	return str + ">"
}
