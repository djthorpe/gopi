// +build drm

package drm

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libdrm
#include <xf86drm.h>
#include <xf86drmMode.h>
#include <errno.h>
#include <stdlib.h>

int _drm_errno() { return errno; }

void _drm_page_flip_handler(int fd, unsigned int frame,unsigned int sec, unsigned int usec, void* data) {
	(void)fd, (void)frame, (void)sec, (void)usec;
	int* pflag = (int* )data;
    *pflag = 1;
}

int _drm_page_flip_wait(int fd,uint32_t crtc_id,uint32_t fb_id) {
	int flag = 0;
	int ret = 0;
	fd_set fds;
	drmEventContext evctx = {
		.version = 2,
		.page_flip_handler = _drm_page_flip_handler,
    };
	ret = drmModePageFlip(fd,crtc_id,fb_id,DRM_MODE_PAGE_FLIP_EVENT,&flag);
	if(ret != 0) {
		return ret;
	}
	struct timeval timeout = {
        .tv_sec = 3,
        .tv_usec = 0
    };
	while(flag == 0) {
	    FD_ZERO(&fds);
	    FD_SET(0, &fds);
		FD_SET(fd, &fds);
		ret = select(fd + 1, &fds, NULL, NULL, &timeout);
	    if (ret < 0) {
			return ret;
	    } else if (ret == 0) {
			return ETIMEDOUT;
	    } else if (FD_ISSET(0, &fds)) {
			return EINTR;
		}
		ret = drmHandleEvent(fd, &evctx);
	    if(ret != 0) {
			return ret;
		}
	}

	return 0;
}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	ModeResources C.drmModeRes
	ModeConnector C.drmModeConnector
	ModeEncoder   C.drmModeEncoder
	ModeCRTC      C.drmModeCrtc
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func GetResources(fd uintptr) (*ModeResources, error) {
	if res := C.drmModeGetResources(C.int(fd)); res == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("GetResources")
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

func GetEncoder(fd uintptr, id uint32) (*ModeEncoder, error) {
	if enc := C.drmModeGetEncoder(C.int(fd), C.uint32_t(id)); enc == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("GetEncoder")
	} else {
		return (*ModeEncoder)(unsafe.Pointer(enc)), nil
	}
}

func GetCRTC(fd uintptr, id uint32) (*ModeCRTC, error) {
	if crtc := C.drmModeGetCrtc(C.int(fd), C.uint32_t(id)); crtc == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("GetCRTC")
	} else {
		return (*ModeCRTC)(unsafe.Pointer(crtc)), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// ModeEncoder

func (this *ModeEncoder) Free() {
	ctx := (*C.drmModeEncoder)(this)
	C.drmModeFreeEncoder(ctx)
}

func (this *ModeEncoder) Id() uint32 {
	ctx := (*C.drmModeEncoder)(this)
	return uint32(ctx.encoder_id)
}

func (this *ModeEncoder) Crtc() uint32 {
	ctx := (*C.drmModeEncoder)(this)
	return uint32(ctx.crtc_id)
}

func (this *ModeEncoder) String() string {
	str := "<drm.encoder"
	str += " id=" + fmt.Sprint(this.Id())
	if crtc := this.Crtc(); crtc != 0 {
		str += " crtc=" + fmt.Sprint(crtc)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// ModeCRTC

func (this *ModeCRTC) Free() {
	ctx := (*C.drmModeCrtc)(this)
	C.drmModeFreeCrtc(ctx)
}

func (this *ModeCRTC) Id() uint32 {
	ctx := (*C.drmModeCrtc)(this)
	return uint32(ctx.crtc_id)
}

func (this *ModeCRTC) String() string {
	str := "<drm.crtc"
	str += " id=" + fmt.Sprint(this.Id())
	return str + ">"
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

func (this *ModeConnector) Status() ModeConnection {
	ctx := (*C.drmModeConnector)(this)
	return ModeConnection(ctx.connection)
}

func (this *ModeConnector) Type() ConnectorType {
	ctx := (*C.drmModeConnector)(this)
	return ConnectorType(ctx.connector_type)
}

func (this *ModeConnector) Dimensions() (uint32, uint32) {
	ctx := (*C.drmModeConnector)(this)
	return uint32(ctx.mmWidth), uint32(ctx.mmHeight)
}

func (this *ModeConnector) Modes() []ModeInfo {
	var result []ModeInfo

	// Make fake slice
	ctx := (*C.drmModeConnector)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_modes)
	sliceHeader.Len = int(ctx.count_modes)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.modes))

	return result
}

func (this *ModeConnector) Encoder() uint32 {
	ctx := (*C.drmModeConnector)(this)
	return uint32(ctx.encoder_id)
}

func (this *ModeConnector) String() string {
	str := "<drm.connector"
	str += " id=" + fmt.Sprint(this.Id())
	if c := this.Status(); c != ModeConnectionNone {
		str += " status=" + fmt.Sprint(c)
	}
	if c := this.Type(); c != DRM_MODE_CONNECTOR_Unknown {
		str += " type=" + fmt.Sprint(c)
	}
	if enc := this.Encoder(); enc != 0 {
		str += " encoder=" + fmt.Sprint(enc)
	}
	if w, h := this.Dimensions(); w > 0 && h > 0 {
		str += fmt.Sprintf(" dimensions={%vmm,%vmm}", w, h)
	}
	str += fmt.Sprintf(" modes=%v", this.Modes())
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

////////////////////////////////////////////////////////////////////////////////
// Frame Buffers

func AddFrameBuffer(fd uintptr, width, height uint32, depth, bpp uint8, stride uint32, handle uintptr) (uint32, error) {
	var id C.uint32_t
	if ret := C.drmModeAddFB(C.int(fd), C.uint32_t(width), C.uint32_t(height), C.uint8_t(depth), C.uint8_t(bpp), C.uint32_t(stride), C.uint32_t(handle), &id); ret != 0 {
		return 0, os.NewSyscallError("drmModeAddFB", syscall.Errno(C._drm_errno()))
	} else {
		return uint32(id), nil
	}
}

func RemoveFrameBuffer(fd uintptr, fb uint32) error {
	if ret := C.drmModeRmFB(C.int(fd), C.uint32_t(fb)); ret != 0 {
		return os.NewSyscallError("drmModeRmFB", syscall.Errno(C._drm_errno()))
	} else {
		return nil
	}
}

func SetCrtc(fd uintptr, crtc, connector uint32, buffer uint32, x, y uint32, mode *ModeInfo) error {
	if ret := C.drmModeSetCrtc(C.int(fd), C.uint32_t(crtc), C.uint32_t(buffer), C.uint32_t(x), C.uint32_t(y), (*C.uint32_t)(&connector), 1, (*C.drmModeModeInfo)(mode)); ret != 0 {
		return os.NewSyscallError("drmModeSetCrtc", syscall.Errno(C._drm_errno()))
	} else {
		return nil
	}
}

func PageFlip(fd uintptr, crtc, fb uint32) error {
	if ret := C._drm_page_flip_wait(C.int(fd), C.uint32_t(crtc), C.uint32_t(fb)); ret != 0 {
		return os.NewSyscallError("drmModePageFlip", syscall.Errno(C._drm_errno()))
	} else {
		return nil
	}
}
