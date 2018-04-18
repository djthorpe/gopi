// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"C"
	"regexp"
	"strconv"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
    #cgo LDFLAGS: -L/opt/vc/lib -lbcm_host
	#include "vc_vchi_gencmd.h"
	int vc_gencmd_wrap(char* response,int maxlen,const char* command) {
		return vc_gencmd(response,maxlen,command);
	}
*/
import "C"
import (
	"strings"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE TYPES

type VideoCore interface {
	gopi.Hardware

	// GeneralCommand executes a VideoCore "General Command"
	GeneralCommand(command string) (string, error)

	// Return list of comamnds
	GeneralCommands() ([]string, error)

	// Return OTP memory
	GetOTP() (map[byte]uint32, error)

	// GetSerialNumberUint64 returns the 64-bit serial number for the device
	GetSerialNumberUint64() (uint64, error)

	// GetRevisionUint32 returns the 32-bit revision code for the device
	GetRevisionUint32() (uint32, error)

	// GetCoreTemperatureCelcius gets CPU core temperature in celcius
	GetCoreTemperatureCelcius() (float64, error)
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
	REGEXP_OTP_DUMP *regexp.Regexp = regexp.MustCompile("(\\d\\d):([0123456789abcdefABCDEF]{8})")
	REGEXP_TEMP     *regexp.Regexp = regexp.MustCompile("temp=(\\d+\\.?\\d*)")
	REGEXP_CLOCK    *regexp.Regexp = regexp.MustCompile("frequency\\((\\d+)\\)=(\\d+)")
	REGEXP_VOLTAGE  *regexp.Regexp = regexp.MustCompile("volt=(\\d*\\.?\\d*)V")
	REGEXP_CODEC    *regexp.Regexp = regexp.MustCompile("(\\w+)=(enabled|disabled)")
	REGEXP_MEMORY   *regexp.Regexp = regexp.MustCompile("(\\w+)=(\\d+)M")
	REGEXP_COMMANDS *regexp.Regexp = regexp.MustCompile("commands=\"([^\"]+)\"")
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GeneralCommand executes a VideoCore "General Command" and return the results
// of that command as a string. See http://elinux.org/RPI_vcgencmd_usage for
// some examples of usage
func (this *hardware) GeneralCommand(command string) (string, error) {
	if this.service == GENCMD_SERVICE_NONE {
		var err error
		this.service, err = vcGencmdInit()
		if err != nil {
			this.service = GENCMD_SERVICE_NONE
			return "", err
		}
	}
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))
	cbuffer := make([]byte, GENCMD_BUF_SIZE)
	if int(C.vc_gencmd_wrap((*C.char)(unsafe.Pointer(&cbuffer[0])), C.int(GENCMD_BUF_SIZE), (*C.char)(unsafe.Pointer(ccommand)))) != 0 {
		return "", gopi.ErrAppError
	}
	return string(cbuffer), nil
}

// Return list of all commands
func (this *hardware) GeneralCommands() ([]string, error) {
	if value, err := this.GeneralCommand(GENCMD_COMMANDS); err != nil {
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

// Return OTP memory
func (this *hardware) GetOTP() (map[byte]uint32, error) {
	// retrieve OTP
	if value, err := this.GeneralCommand(GENCMD_OTP_DUMP); err != nil {
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

// GetSerialNumberUint64 returns the 64-bit serial number for the device
func (this *hardware) GetSerialNumberUint64() (uint64, error) {
	// Return cached version before fetching memory
	if this.serial != GENCMD_SERIAL_NONE {
		return this.serial, nil
	} else if otp, err := this.GetOTP(); err != nil {
		return GENCMD_SERIAL_NONE, err
	} else {
		this.serial = uint64(otp[GENCMD_OTP_DUMP_SERIAL])
		return this.serial, nil
	}
}

// GetRevisionUint32 returns the 32-bit revision code for the device
func (this *hardware) GetRevisionUint32() (uint32, error) {
	// Return cached version before fetching memory
	if this.revision != GENCMD_REVISION_NONE {
		return this.revision, nil
	} else if otp, err := this.GetOTP(); err != nil {
		return GENCMD_REVISION_NONE, err
	} else {
		this.revision = uint32(otp[GENCMD_OTP_DUMP_REVISION])
		return this.revision, nil
	}
}

// GetCoreTemperatureCelcius gets CPU core temperature in celcius
func (this *hardware) GetCoreTemperatureCelcius() (float64, error) {
	// Retrieve value as text
	if value, err := this.GeneralCommand(GENCMD_MEASURE_TEMP); err != nil {
		return 0.0, err
	} else if match := REGEXP_TEMP.FindStringSubmatch(value); len(match) != 2 {
		return 0.0, gopi.ErrUnexpectedResponse
	} else if value2, err := strconv.ParseFloat(match[1], 64); err != nil {
		return 0.0, gopi.ErrUnexpectedResponse
	} else {
		return value2, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func vcGencmdInit() (int, error) {
	service := int(C.vc_gencmd_init())
	if service < 0 {
		return -1, gopi.ErrAppError
	}
	return service, nil
}

func vcGencmdTerminate() error {
	C.vc_gencmd_stop()
	return nil
}
