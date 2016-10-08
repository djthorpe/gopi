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

////////////////////////////////////////////////////////////////////////////////

const (
	BMC2835_VIDEOCORE_FILENAME  = "/dev/vchiq"
	BMC2835_DEVICETREE_FILENAME = "/proc/device-tree/soc/ranges"
)

////////////////////////////////////////////////////////////////////////////////

type RaspberryPi struct {
	revision uint32
	serial   uint64
	peribase uint32
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

func _BCMHostInit() error {
	if err := isWritablePath(BMC2835_VIDEOCORE_FILENAME); err != nil {
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

	// initialize broadcom host
	if err := _BCMHostInit(); err != nil {
		return nil, err
	}

	// initialize videocore device
	err := _VCGenCmdInit()
	if err != nil {
		_BCMHostTerminate()
		return nil, err
	}

	// return success
	return this, nil
}

// Close RaspberryPi object
func (this *RaspberryPi) Close() {
	_VCGenCmdStop()
	_BCMHostTerminate()
}

////////////////////////////////////////////////////////////////////////////////

/*
// Read /proc/device-tree/soc/ranges and determine the base address.
// Use the default Raspberry Pi 1 base address if this fails.
func (this *RaspberryPi) getBaseAddress() (uint32,error) {
	peripheralbase, err := this.PeripheralBase()
	if err != nil {
		return 0,err
	}
	ranges, err := os.Open(BMC2835_DEVICETREE_FILENAME)
	if err != nil {
		return uint32(peripheralbase + GPIO_BASE), nil
	}
	defer ranges.Close()
	b := make([]byte, 4)
	n, err := ranges.ReadAt(b, 4)
	if n != 4 || err != nil {
		return uint32(peripheralbase + GPIO_BASE), nil
	}
	buf := bytes.NewReader(b)
	var out uint32
	err = binary.Read(buf,binary.BigEndian,&out)
	if err != nil {
		return uint32(peripheralbase + GPIO_BASE), nil
	}
	return uint32(out + GPIO_BASE), nil
}
*/
