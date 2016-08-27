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
*/
import "C"

import (
	"os"
)

////////////////////////////////////////////////////////////////////////////////

const (
	VIDEOCORE_DEVICE = "/dev/vchiq"
)

////////////////////////////////////////////////////////////////////////////////

type RaspberryPi struct {
	revision uint32
	serial   uint64
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

func _BCMHostInit() error {
	// TODO: Ensure /dev/vchiq is readable and writable
	_, err := os.Stat(VIDEOCORE_DEVICE)
	if err != nil {
		return ErrorVchiq
	}
	C.bcm_host_init()
	return nil
}

func _BCMHostTerminate() {
	C.bcm_host_deinit()
}

func _GraphicsGetDisplaySize(displayNumber uint16) (uint32, uint32) {
	var w, h uint32
	C.graphics_get_display_size((C.uint16_t)(displayNumber), (*C.uint32_t)(&w), (*C.uint32_t)(&h))
	return w, h
}

func _VCGenCmdInit() error {
	if C.vc_gencmd_init() >= 0 {
		return nil
	}
	return ErrorInit
}

func _VCGenCmdStop() {
	C.vc_gencmd_stop()
}

////////////////////////////////////////////////////////////////////////////////

// Create a new RaspberryPi object
func New() (*RaspberryPi, error) {
	// create this object
	this := new(RaspberryPi)
	// initialize
	_BCMHostInit()
	err := _VCGenCmdInit()
	if err != nil {
		return nil, err
	}
	return this, nil
}

// Close RaspberryPi object
func (this *RaspberryPi) Close() {
	_VCGenCmdStop()
	_BCMHostTerminate()
}
