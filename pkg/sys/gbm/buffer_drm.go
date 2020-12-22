// +build gbm,drm

package gbm

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: gbm
#cgo pkg-config: libdrm
#include <xf86drm.h>
#include <xf86drmMode.h>
#include <gbm.h>
#include <stdio.h>

int _errno();
uint32_t gbm_bo_handle(struct gbm_bo* bo);

void gbm_destroy_buffer_callback(struct gbm_bo* bo, void* data) {
    int fd = gbm_device_get_fd(gbm_bo_get_device(bo));
    uint32_t fb = (uint32_t)data;

	printf("drmModeRmFB %d\n",fb);
    if(fb) {
		printf("->free\n");
		drmModeRmFB(fd,fb);
	}
}
void gbm_set_buffer_userdata(struct gbm_bo* bo,void* data) {
	printf("->gbm_bo_set_user_data %p\n",data);
	gbm_bo_set_user_data(bo,data,gbm_destroy_buffer_callback);
}
*/
import "C"
import (
	"os"
	"syscall"
	"unsafe"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MAX_PLANES = 4
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *GBMBuffer) NewFrameBuffer() (uint32, error) {
	ctx := (*C.struct_gbm_bo)(this)

	if fb := uintptr(C.gbm_bo_get_user_data(ctx)); fb != 0 {
		return uint32(fb), nil
	}

	fd := C.gbm_device_get_fd(C.gbm_bo_get_device(ctx))
	w, h := C.gbm_bo_get_width(ctx), C.gbm_bo_get_height(ctx)
	format := C.gbm_bo_get_format(ctx)
	flags := uint32(0)
	var strides, offsets, handles [MAX_PLANES]C.uint32_t

	count := int(C.gbm_bo_get_plane_count(ctx))
	if count > MAX_PLANES {
		return 0, gopi.ErrInternalAppError.WithPrefix("NewFrameBuffer")
	}
	for plane := 0; plane < count; plane++ {
		strides[plane] = C.gbm_bo_get_stride_for_plane(ctx, C.int(plane))
		handles[plane] = C.gbm_bo_handle(ctx)
		offsets[plane] = C.gbm_bo_get_offset(ctx, C.int(plane))
	}

	var fb C.uint32_t
	if ret := C.drmModeAddFB2(fd, w, h, format, &handles[0], &strides[0], &offsets[0], &fb, C.uint32_t(flags)); ret != 0 {
		return 0, os.NewSyscallError("drmModeAddFB2", syscall.Errno(C._errno()))
	}
	C.gbm_set_buffer_userdata(ctx, unsafe.Pointer(uintptr(fb)))
	return uint32(fb), nil
}

/*
func (this *GBMBuffer) NewFrameBuffer() (uint32, error) {
	ctx := (*C.struct_gbm_bo)(this)

	if fb := uintptr(C.gbm_bo_get_user_data(ctx)); fb != 0 {
		return uint32(fb), nil
	}

	fd := C.gbm_device_get_fd(C.gbm_bo_get_device(ctx))
	w, h := C.gbm_bo_get_width(ctx), C.gbm_bo_get_height(ctx)
	format := C.gbm_bo_get_format(ctx)
	flags := uint32(0)
	var strides, offsets, handles [MAX_PLANES]C.uint32_t
	var modifiers [MAX_PLANES]C.uint64_t

	count := int(C.gbm_bo_get_plane_count(ctx))
	if count > MAX_PLANES {
		return 0, gopi.ErrInternalAppError.WithPrefix("NewFrameBuffer")
	}
	for plane := 0; plane < count; plane++ {
		strides[plane] = C.gbm_bo_get_stride_for_plane(ctx, C.int(plane))
		handles[plane] = C.gbm_bo_handle(ctx)
		offsets[plane] = C.gbm_bo_get_offset(ctx, C.int(plane))
		modifiers[plane] = C.gbm_bo_get_modifier(ctx)
	}
	if modifiers[0] != 0 {
		flags = drm.DRM_MODE_FB_MODIFIERS
	}

	var fb C.uint32_t
	if ret := C.drmModeAddFB2WithModifiers(fd, w, h, format, &handles[0], &strides[0], &offsets[0], &modifiers[0], &fb, C.uint32_t(flags)); ret != 0 {
		return 0, os.NewSyscallError("drmModeAddFB2WithModifiers", syscall.Errno(C._errno()))
	}
	C.gbm_set_buffer_userdata(ctx, unsafe.Pointer(uintptr(fb)))
	return uint32(fb), nil
}
*/
