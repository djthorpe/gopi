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
	"fmt"
	"strings"
	"unsafe"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: bcm_host
#include <interface/vmcs_host/vc_tvservice.h>
#include <interface/vmcs_host/vc_hdmi.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	VCHIInstance       C.VCHI_INSTANCE_T
	VCHIConnection     C.VCHI_CONNECTION_T
	TVDisplayState     C.TV_DISPLAY_STATE_T
	TVDisplayInfo      C.TV_DEVICE_ID_T
	TVError            uint
	TVDisplayStateFlag uint32
)

////////////////////////////////////////////////////////////////////////////////
// CONST

const (
	ErrNone TVError = iota
	ErrConnectError
	ErrDeviceError
)

const (
	TV_MAX_ATTACHED_DISPLAYS = 16
)

const (
	TV_STATE_HDMI_UNPLUGGED         TVDisplayStateFlag = (1 << iota) // HDMI cable is detached
	TV_STATE_HDMI_ATTACHED                                           // HDMI cable is attached but not powered on
	TV_STATE_HDMI_DVI                                                // HDMI is on but in DVI mode (no audio)
	TV_STATE_HDMI_HDMI                                               //HDMI is on and HDMI mode is active
	TV_STATE_HDMI_HDCP_UNAUTH                                        // HDCP authentication is broken (e.g. Ri mismatched) or not active
	TV_STATE_HDMI_HDCP_AUTH                                          // HDCP is active
	TV_STATE_HDMI_HDCP_KEY_DOWNLOAD                                  // HDCP key download successful/fail
	TV_STATE_HDMI_HDCP_SRM_DOWNLOAD                                  // HDCP revocation list download successful/fail
	TV_STATE_HDMI_CHANGING_MODE                                      // HDMI is starting to change mode, clock has not yet been set
	TV_STATE_SDTV_UNPLUGGED         TVDisplayStateFlag = 1 << 16     // SDTV cable unplugged, subject to platform support
	TV_STATE_SDTV_ATTACHED          TVDisplayStateFlag = 1 << 17     // SDTV cable is plugged in
	TV_STATE_SDTV_NTSC              TVDisplayStateFlag = 1 << 18     // SDTV is in NTSC mode
	TV_STATE_SDTV_PAL               TVDisplayStateFlag = 1 << 19     // SDTV is in PAL mode
	TV_STATE_SDTV_CP_INACTIVE       TVDisplayStateFlag = 1 << 20     // Copy protection disabled
	TV_STATE_SDTV_CP_ACTIVE         TVDisplayStateFlag = 1 << 21     // Copy protection enabled
	TV_STATE_LCD_ATTACHED_DEFAULT   TVDisplayStateFlag = 1 << 22     // LCD display is attached and default
	TV_STATE_MIN                                       = TV_STATE_HDMI_UNPLUGGED
	TV_STATE_MAX                                       = TV_STATE_LCD_ATTACHED_DEFAULT
	TV_STATE_NONE                   TVDisplayStateFlag = 0
)

/*
vc_tv_get_device_id_id
*/

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

func VCHI_Init() VCHIInstance {
	var instance VCHIInstance
	if err := C.vchi_initialise((*C.VCHI_INSTANCE_T)(unsafe.Pointer(&instance))); err != 0 {
		return nil
	} else {
		return instance
	}
}

func VCHI_TVInit(instance VCHIInstance) (VCHIConnection, error) {
	var connection VCHIConnection
	if err := C.vchi_connect(nil, 0, C.VCHI_INSTANCE_T(instance)); err != 0 {
		return connection, ErrConnectError
	} else if err := C.vc_vchi_tv_init(C.VCHI_INSTANCE_T(instance), (**C.VCHI_CONNECTION_T)(unsafe.Pointer(&connection)), 1); err != 0 {
		return connection, ErrConnectError
	} else {
		return connection, nil
	}
}

func VCHI_TVStop(instance VCHIInstance) error {
	C.vc_vchi_tv_stop()
	C.vchi_disconnect(C.VCHI_INSTANCE_T(instance))
	return nil
}

func VCHI_TVGetAttachedDevices() ([]DXDisplayId, error) {
	var devices C.TV_ATTACHED_DEVICES_T
	if err := C.vc_tv_get_attached_devices((*C.TV_ATTACHED_DEVICES_T)(unsafe.Pointer(&devices))); err != 0 {
		return nil, ErrDeviceError
	}
	displays := make([]DXDisplayId, 0, TV_MAX_ATTACHED_DISPLAYS)
	for device := 0; device < int(devices.num_attached); device++ {
		displays = append(displays, DXDisplayId(devices.display_number[device]))
	}
	return displays, nil
}

func VCHI_TVGetDisplayState(display DXDisplayId) (TVDisplayState, error) {
	var state C.TV_DISPLAY_STATE_T
	if err := C.vc_tv_get_display_state_id(C.uint32_t(display), (*C.TV_DISPLAY_STATE_T)(unsafe.Pointer(&state))); err != 0 {
		return TVDisplayState(state), ErrDeviceError
	} else {
		return TVDisplayState(state), nil
	}
}

func VCHI_TVHDMIPowerOnPreferred(display DXDisplayId) error {
	if err := C.vc_tv_hdmi_power_on_preferred_id(C.uint32_t(display));  err != 0 {
		return ErrDeviceError
	} else {
		return nil
	}
}

func VCHI_TVHDMIPowerOn(display DXDisplayId, width, height, framerate uint32, interlaced bool) error {
	interlaced_ := C.HDMI_INTERLACED_T(C.HDMI_NONINTERLACED)
	if interlaced {
		interlaced_ = C.HDMI_INTERLACED_T(C.HDMI_INTERLACED)
	}	
	if err := C.vc_tv_hdmi_power_on_best_id(C.uint32_t(display),C.uint32_t(width),C.uint32_t(height),C.uint32_t(framerate),interlaced_,C.HDMI_MODE_MATCH_RESOLUTION);  err != 0 {
		return ErrDeviceError
	} else {
		return nil
	}
}

/*
func VCHI_TVSDPowerOn(display DXDisplayId) error {
	if err := C.vc_tv_sdtv_power_on_id(C.uint32_t(display),mode,);  err != 0 {
		return ErrDeviceError
	} else {
		return nil
	}
}
*/

func VCHI_TVPowerOff(display DXDisplayId) error {
	if err := C.vc_tv_power_off_id(C.uint32_t(display));  err != 0 {
		return ErrDeviceError
	} else {
		return nil
	}
}

func VCHI_TVGetDisplayInfo(display DXDisplayId) (TVDisplayInfo,error) {
	var info TVDisplayInfo
	if err := C.vc_tv_get_device_id_id(C.uint32_t(display),(*C.TV_DEVICE_ID_T)(unsafe.Pointer(&info))); err != 0 {
		return info,ErrDeviceError
	} else {
		return info,nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// TVDisplayInfo

func (this TVDisplayInfo) Vendor() string {
	ptr := C.TV_DEVICE_ID_T(this).vendor
	return C.GoString(&ptr[0])
}

func (this TVDisplayInfo) Product() string {
	ptr := C.TV_DEVICE_ID_T(this).monitor_name
	return C.GoString(&ptr[0])
}


func (this TVDisplayInfo) Serial() uint32 {
	return uint32(C.TV_DEVICE_ID_T(this).serial_num)
}

func (this TVDisplayInfo) String() string {
	return "<TVDisplayInfo" +
		" vendor=" + strconv.Quote(this.Vendor()) +
		" product=" + strconv.Quote(this.Product()) +
		" serial=" + fmt.Sprint(this.Serial) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// TVDisplayState

func (this TVDisplayState) Flags() TVDisplayStateFlag {
	return TVDisplayStateFlag(this.state)
}

func (this TVDisplayState) String() string {
	return "<TVDisplayState state=" + fmt.Sprint(TVDisplayStateFlag(this.state)) + ">"
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e TVError) Error() string {
	switch e {
	case ErrNone:
		return "ErrNone"
	case ErrConnectError:
		return "ErrConnectError"
	case ErrDeviceError:
		return "ErrDeviceError"
	default:
		return "[?? Invalid TVError value]"
	}
}

func (f TVDisplayStateFlag) String() string {
	str := ""
	if f == TV_STATE_NONE {
		return f.String()
	}
	for v := TV_STATE_MIN; v <= TV_STATE_MAX; v <<= 1 {
		if f&v != 0 {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (f TVDisplayStateFlag) StringFlag() string {
	switch f {
	case TV_STATE_HDMI_UNPLUGGED:
		return "TV_STATE_HDMI_UNPLUGGED"
	case TV_STATE_HDMI_ATTACHED:
		return "TV_STATE_HDMI_ATTACHED"
	case TV_STATE_HDMI_DVI:
		return "TV_STATE_HDMI_DVI"
	case TV_STATE_HDMI_HDMI:
		return "TV_STATE_HDMI_HDMI"
	case TV_STATE_HDMI_HDCP_UNAUTH:
		return "TV_STATE_HDMI_HDCP_UNAUTH"
	case TV_STATE_HDMI_HDCP_AUTH:
		return "TV_STATE_HDMI_HDCP_AUTH"
	case TV_STATE_HDMI_HDCP_KEY_DOWNLOAD:
		return "TV_STATE_HDMI_HDCP_KEY_DOWNLOAD"
	case TV_STATE_HDMI_HDCP_SRM_DOWNLOAD:
		return "TV_STATE_HDMI_HDCP_SRM_DOWNLOAD"
	case TV_STATE_HDMI_CHANGING_MODE:
		return "TV_STATE_HDMI_CHANGING_MODE"
	case TV_STATE_SDTV_UNPLUGGED:
		return "TV_STATE_SDTV_UNPLUGGED"
	case TV_STATE_SDTV_ATTACHED:
		return "TV_STATE_SDTV_ATTACHED"
	case TV_STATE_SDTV_NTSC:
		return "TV_STATE_SDTV_NTSC"
	case TV_STATE_SDTV_PAL:
		return "TV_STATE_SDTV_PAL"
	case TV_STATE_SDTV_CP_INACTIVE:
		return "TV_STATE_SDTV_CP_INACTIVE"
	case TV_STATE_SDTV_CP_ACTIVE:
		return "TV_STATE_SDTV_CP_ACTIVE"
	case TV_STATE_LCD_ATTACHED_DEFAULT:
		return "TV_STATE_LCD_ATTACHED_DEFAULT"
	default:
		return "[?? Invalid TVDisplayStateFlags value]"
	}
}
