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

extern int _drm_errno();
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// Get and Set Client Capabilities

func GetCap(fd uintptr, key uint64) (uint64, error) {
	var value C.uint64_t
	if ret := C.drmGetCap(C.int(fd), C.uint64_t(key), &value); ret != 0 {
		return 0, os.NewSyscallError("drmGetCap", syscall.Errno(C._drm_errno()))
	} else {
		return uint64(value), nil
	}
}

func SetClientCap(fd uintptr, key, value uint64) error {
	if ret := C.drmSetClientCap(C.int(fd), C.uint64_t(key), C.uint64_t(value)); ret != 0 {
		return os.NewSyscallError("drmSetClientCap", syscall.Errno(C._drm_errno()))
	} else {
		return nil
	}
}
