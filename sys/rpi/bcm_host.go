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
	"sync"
	"regexp"
	"strings"
	"strconv"

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

// OTP (One Time Programmable) memory constants
const (
	GENCMD_OTP_DUMP          = "otp_dump"
	GENCMD_OTP_DUMP_SERIAL   = 28
	GENCMD_OTP_DUMP_REVISION = 30
	GENCMD_COMMANDS          = "commands"
	GENCMD_MEASURE_TEMP      = "measure_temp"
	GENCMD_MEASURE_CLOCK     = "measure_clock arm core h264 isp v3d uart pwm emmc pixel vec hdmi dpi"
	GENCMD_MEASURE_VOLTS     = "measure_volts core sdram_c sdram_i sdram_p"
	GENCMD_CODEC_ENABLED     = "codec_enabled H264 MPG2 WVC1 MPG4 MJPG WMV9 VP8"
	GENCMD_MEMORY            = "get_mem arm gpu"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	gencmdBuffer = make([]byte, GENCMD_BUF_SIZE)
	bcmHostOnce sync.Once
	vcGenCmdOnce sync.Once
	vcGenCmdService int
)

var (
	REGEXP_OTP_DUMP *regexp.Regexp = regexp.MustCompile("(\\d\\d):([0123456789abcdefABCDEF]{8})")
	REGEXP_TEMP     *regexp.Regexp = regexp.MustCompile("temp=(\\d+\\.?\\d*)")
	REGEXP_CLOCK    *regexp.Regexp = regexp.MustCompile("frequency\\((\\d+)\\)=(\\d+)")
	REGEXP_VOLTAGE  *regexp.Regexp = regexp.MustCompile("volt=(\\d*\\.?\\d*)V")
	REGEXP_CODEC    *regexp.Regexp = regexp.MustCompile("(\\w+)=(enabled|disabled)")
	REGEXP_MEMORY   *regexp.Regexp = regexp.MustCompile("(\\w+)=(\\d+)M")
	REGEXP_COMMANDS *regexp.Regexp = regexp.MustCompile("commands=\"([^\"]+)\"")
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: BCMHOST

func BCMHostInit() error {
	bcmHostOnce.Do(func() {
		C.bcm_host_init()
	})
	return nil
}

func BCMHostTerminate(){
	C.bcm_host_deinit()
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
	vcGenCmdOnce.Do(func() {
		vcGenCmdService = int(C.vc_gencmd_init())
	})
	return vcGenCmdService
}

func VCGencmdTerminate() {
	C.vc_gencmd_stop()
}

// GeneralCommand executes a VideoCore "General Command" and return the results
// of that command as a string. See http://elinux.org/RPI_vcgencmd_usage for
// some examples of usage
func VCGeneralCommand(command string) (string, error) {
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))
	ptr := (*C.char)(&gencmdBuffer[0])	
	if ret := C.vc_gencmd_wrap(ptr,C.int(GENCMD_BUF_SIZE),ccommand); ret != 0 {
		return "", gopi.ErrUnexpectedResponse.WithPrefix(fmt.Sprint(ret))
	} else {
		return C.GoString(ptr), nil
	}
}

// Return list of all commands
func VCGeneralCommands() ([]string, error) {
	if value, err := VCGeneralCommand(GENCMD_COMMANDS); err != nil {
		return nil, err
	} else if matches := REGEXP_COMMANDS.FindStringSubmatch(value); len(matches) < 2 {
		return nil, gopi.ErrUnexpectedResponse
	} else {
		cmds := make([]string, 0)
		for _, cmd := range strings.Split(matches[1], ",") {
			cmds = append(cmds, strings.TrimSpace(cmd))
		}
		return cmds, nil
	}
}

// VCOTPDump returns OTP memory dump
func VCOTPDump() (map[byte]uint32, error) {
	if value, err := VCGeneralCommand(GENCMD_OTP_DUMP); err != nil {
		return nil, err
	} else if matches := REGEXP_OTP_DUMP.FindAllStringSubmatch(value, -1); len(matches) == 0 {
		return nil, gopi.ErrUnexpectedResponse
	} else {
		otp := make(map[byte]uint32, len(matches))
		for _, match := range matches {
			if len(match) != 3 {
				return nil, gopi.ErrUnexpectedResponse
			}
			if index, err := strconv.ParseUint(match[1], 10, 8); err != nil {
				return nil, gopi.ErrUnexpectedResponse
			} else if value, err := strconv.ParseUint(match[2], 16, 32); err != nil {
				return nil, gopi.ErrUnexpectedResponse
			} else {
				otp[byte(index)] = uint32(value)
			}
		}
		return otp, nil
	}
}

// VCGetSerialProduct returns the 64-bit serial number and 32-bit product number for the device
func VCGetSerialProduct() (uint64, uint32, error) {
	if otp, err := VCOTPDump(); err != nil {
		return 0, 0, err
	} else {
		return uint64(otp[GENCMD_OTP_DUMP_SERIAL]), uint32(otp[GENCMD_OTP_DUMP_REVISION]), nil
	}
}
