// +build drm

package drm

import (
	"os"
	"syscall"
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
	// This version of PageFlip blocks until page flip has occurred, or
	// returns an error otherwise
	if ret := C._drm_page_flip_wait(C.int(fd), C.uint32_t(crtc), C.uint32_t(fb)); ret != 0 {
		return os.NewSyscallError("drmModePageFlip", syscall.Errno(C._drm_errno()))
	} else {
		return nil
	}
}
