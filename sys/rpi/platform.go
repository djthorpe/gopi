// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: bcm_host
#include "bcm_host.h"
int vc_gencmd_wrap(char* response,int maxlen,const char* command) {
  return vc_gencmd(response,maxlen,command);
}
*/
import "C"

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
