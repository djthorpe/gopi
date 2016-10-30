/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
	"regexp"
	"strconv"
	"unsafe"
)

import (
	gopi "../.."      /* import "github.com/djthorpe/gopi" */
	util "../../util" /* import "github.com/djthorpe/gopi/util" */
)

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
    #cgo LDFLAGS: -L/opt/vc/lib -lbcm_host
    #include "bcm_host.h"
	#include "vc_vchi_gencmd.h"
	int vc_gencmd_wrap(char* response,int maxlen,const char* command) {
		return vc_gencmd(response,maxlen,command);
	}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Device struct{}

type DeviceState struct {
	log      *util.LoggerDevice // logger
	service  int                // service number
	serial   uint64
	revision uint32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GENCMD_BUF_SIZE      = 1024
	GENCMD_SERVICE_NONE  = -1
	GENCMD_SERIAL_NONE   = 0
	GENCMD_REVISION_NONE = 0
)

// OTP (One Time Programmable) memory constants
const (
	GENCMD_OTP_DUMP          = "otp_dump"
	GENCMD_OTP_DUMP_SERIAL   = 28
	GENCMD_OTP_DUMP_REVISION = 30
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	REGEXP_OTP_DUMP *regexp.Regexp = regexp.MustCompile("(\\d\\d):([0123456789abcdefABCDEF]{8})")
)

////////////////////////////////////////////////////////////////////////////////
// Open and close device

func (config Device) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug2("<rpi.Device>Open")
	if err := bcmHostInit(); err != nil {
		return nil, err
	}
	return &DeviceState{log, GENCMD_SERVICE_NONE, GENCMD_SERIAL_NONE, GENCMD_REVISION_NONE}, nil
}

func (this *DeviceState) Close() error {
	this.log.Debug2("<rpi.Device>Close")
	if this.service != GENCMD_SERVICE_NONE {
		if err := vcGencmdTerminate(); err != nil {
			bcmHostTerminate()
			return err
		}
	}
	if err := bcmHostTerminate(); err != nil {
		return err
	}
	return nil
}

func (this *DeviceState) String() string {
	serial, _ := this.GetSerialNumber()
	revision, _ := this.GetRevision()
	model, pcb, _ := this.GetModel()
	processor, _ := this.GetProcessor()
	warranty_bit, _ := this.GetWarrantyBit()
	return fmt.Sprintf("<rpi.Device>{ serial_number=%08X revision=%04X model=%v pcb=%v processor=%v warranty_bit=%v }", serial, revision, model, pcb, processor, warranty_bit)
}

////////////////////////////////////////////////////////////////////////////////
// Get Device Information

func (this *DeviceState) GetPeripheralAddress() uint32 {
	return bcmHostGetPeripheralAddress()
}

func (this *DeviceState) GetPeripheralSize() uint32 {
	return bcmHostGetPeripheralSize()
}

// Return the 64-bit serial number for the device
func (this *DeviceState) GetSerialNumber() (uint64, error) {
	// Return cached version
	if this.serial != GENCMD_SERIAL_NONE {
		return this.serial, nil
	}

	// Get embedded memory
	otp, err := this.GetOTP()
	if err != nil {
		return GENCMD_SERIAL_NONE, err
	}
	// Cache and return serial number
	this.serial = uint64(otp[GENCMD_OTP_DUMP_SERIAL])
	return this.serial, nil
}

// Return the 32-bit revision code for the device
func (this *DeviceState) GetRevision() (uint32, error) {
	// Return cached version
	if this.revision != GENCMD_REVISION_NONE {
		return this.revision, nil
	}

	// Get embedded memory
	otp, err := this.GetOTP()
	if err != nil {
		return GENCMD_REVISION_NONE, err
	}

	// Cache and return revision number
	this.revision = uint32(otp[GENCMD_OTP_DUMP_REVISION])
	return this.revision, nil
}

// Return the size of a particular display
func (this *DeviceState) GetDisplaySize(display uint16) (uint32, uint32) {
	return bcmGHostGetDisplaySize(display)
}

////////////////////////////////////////////////////////////////////////////////
// General Command Interface

// Execute a VideoCore "General Command" and return the results of
// that command. See http://elinux.org/RPI_vcgencmd_usage for some example
// usage
func (this *DeviceState) GeneralCommand(command string) (string, error) {
	if this.service == GENCMD_SERVICE_NONE {
		var err error
		this.service, err = vcGencmdInit(this.log)
		if err != nil {
			this.service = GENCMD_SERVICE_NONE
			return "", err
		}
	}
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))
	cbuffer := make([]byte, GENCMD_BUF_SIZE)
	if int(C.vc_gencmd_wrap((*C.char)(unsafe.Pointer(&cbuffer[0])), C.int(GENCMD_BUF_SIZE), (*C.char)(unsafe.Pointer(ccommand)))) != 0 {
		return "", this.log.Error("General Command Error")
	}
	return string(cbuffer), nil
}

// Return OTP memory
func (this *DeviceState) GetOTP() (map[byte]uint32, error) {
	// retrieve OTP
	value, err := this.GeneralCommand(GENCMD_OTP_DUMP)
	if err != nil {
		return nil, err
	}

	// find matches in the text
	matches := REGEXP_OTP_DUMP.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return nil, this.log.Error("Bad Response from %v", GENCMD_OTP_DUMP)
	}
	otp := make(map[byte]uint32, len(matches))
	for _, match := range matches {
		if len(match) != 3 {
			return nil, this.log.Error("Bad Response from %v", GENCMD_OTP_DUMP)
		}
		index, err := strconv.ParseUint(match[1], 10, 8)
		if err != nil {
			return nil, err
		}
		value, err := strconv.ParseUint(match[2], 16, 32)
		if err != nil {
			return nil, err
		}
		otp[byte(index)] = uint32(value)
	}

	return otp, nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

func bcmHostInit() error {
	C.bcm_host_init()
	return nil
}

func bcmHostTerminate() error {
	C.bcm_host_deinit()
	return nil
}

func vcGencmdInit(log *util.LoggerDevice) (int, error) {
	service := int(C.vc_gencmd_init())
	if service < 0 {
		return -1, log.Error("vc_gencmd_init failed")
	}
	return service, nil
}

func vcGencmdTerminate() error {
	C.vc_gencmd_stop()
	return nil
}

func bcmHostGetPeripheralAddress() uint32 {
	return uint32(C.bcm_host_get_peripheral_address())
}

func bcmHostGetPeripheralSize() uint32 {
	return uint32(C.bcm_host_get_peripheral_size())
}

func bcmHostGetSDRAMAddress() uint32 {
	return uint32(C.bcm_host_get_sdram_address())
}

func bcmGHostGetDisplaySize(display uint16) (uint32, uint32) {
	var w, h uint32
	C.graphics_get_display_size((C.uint16_t)(display), (*C.uint32_t)(&w), (*C.uint32_t)(&h))
	return w, h
}
