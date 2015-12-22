/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
    #cgo LDFLAGS: -L/opt/vc/lib -lbcm_host
    #include "bcm_host.h"
	#include "vc_vchi_gencmd.h"

	int vc_gencmd_wrapper(char* response,int maxlen,const char* command) {
		return vc_gencmd(response,maxlen,command);
	}
*/
import "C"

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

const (
	GENCMD_BUFFER_SIZE = 1024
)

////////////////////////////////////////////////////////////////////////////////

func BCMHostInit() {
	C.bcm_host_init()
}

func BCMHostTerminate() {
	C.bcm_host_deinit()
}

func GraphicsGetDisplaySize(displayNumber uint16) (uint32, uint32) {
	var w, h uint32
	C.graphics_get_display_size((C.uint16_t)(displayNumber), (*C.uint32_t)(&w), (*C.uint32_t)(&h))
	return w, h
}

func VCGenCmdInit() error {
	if C.vc_gencmd_init() >= 0 {
		return nil
	}
	return ErrorInit
}

func VCGenCmdStop() {
	C.vc_gencmd_stop()
}

func VCGenCmd(command string) (string,error) {
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))
	cbuffer := make([]byte,GENCMD_BUFFER_SIZE)
	if int(C.vc_gencmd_wrapper((*C.char)(unsafe.Pointer(&cbuffer[0])),C.int(GENCMD_BUFFER_SIZE),(*C.char)(unsafe.Pointer(ccommand)))) != 0 {
		return "",ErrorGenCmd
	}
	return string(cbuffer),nil
}


//,,)