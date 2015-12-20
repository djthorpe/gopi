package egl

/*
	#cgo CFLAGS: -I/opt/vc/include
	#cgo LDFLAGS: -L/opt/vc/lib -llibbcm_host
	#include "bcm_host.h"
*/
import "C"

func BCMHostInit() {
	C.bcm_host_init()
}

func GraphicsGetDisplaySize(displayNumber uint16) (uint32,uint32) {
	var w,h uint32
	C.graphics_get_display_size((C.uint16_t)(displayNumber),(*C.uint32_t)(&w),(*C.uint32_t)(&h))
	return w,h
}