// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"unsafe"	
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: bcm_host
#include "bcm_host.h"
#include <stdio.h>
int vc_gencmd_wrap(char* response,int maxlen,const char* command) {
	return vc_gencmd(response,maxlen,command);
}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GENCMD_BUF_SIZE     = 1024
)

var (
	gencmdBuffer = make([]byte, GENCMD_BUF_SIZE)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: BCMHOST

func BCMHostInit() error {
	C.bcm_host_init()
	return nil
}

func BCMHostTerminate() error {
	C.bcm_host_deinit()
	return nil
}

func BCMHostGetPeripheralAddress() uint32 {
	return uint32(C.bcm_host_get_peripheral_address())
}

func BCMHostGetPeripheralSize() uint32 {
	return uint32(C.bcm_host_get_peripheral_size())
}

func BCMHostGetSDRAMAddress() uint32 {
	return uint32(C.bcm_host_get_sdram_address())
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: VIDEOCORE

func VCGencmdInit() int {
	return int(C.vc_gencmd_init())
}

func VCGencmdTerminate() {
	C.vc_gencmd_stop()
}

// GeneralCommand executes a VideoCore "General Command" and return the results
// of that command as a string. See http://elinux.org/RPI_vcgencmd_usage for
// some examples of usage
func VCGeneralCommand(command string) (string, error) {
	fmt.Println(command)
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))
	ret := C.vc_gencmd_wrap(
		(*C.char)(&gencmdBuffer[0]),
		C.int(GENCMD_BUF_SIZE),
		(*C.char)(ccommand))
	if ret != 0 {
		fmt.Println("buf",string(gencmdBuffer))
		return "", gopi.ErrUnexpectedResponse.WithPrefix(fmt.Sprint(ret))
	}
	return string(gencmdBuffer), nil
}
