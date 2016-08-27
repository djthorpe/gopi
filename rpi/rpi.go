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
	GENCMD_BUFFER_SIZE      = 1024
	GEMCMD_COMMANDS         = "commands"
	GENCMD_MEASURE_TEMP     = "measure_temp"
	GENCMD_MEASURE_CLOCK    = "measure_clock arm core h264 isp v3d uart pwm emmc pixel vec hdmi dpi"
	GENCMD_MEASURE_VOLTS    = "measure_volts core sdram_c sdram_i sdram_p"
	GENCMD_CODEC_ENABLED    = "codec_enabled H264 MPG2 WVC1 MPG4 MJPG WMV9 VP8"
	GENCMD_MEMORY           = "get_mem arm gpu"
	GENCMD_OTPDUMP          = "otp_dump"
	GENCMD_OTPDUMP_SERIAL   = 28
	GENCMD_OTPDUMP_REVISION = 30
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
