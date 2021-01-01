//+build mmal

package mmal

import (
	"fmt"
	"reflect"
	"unsafe"

	hw "github.com/djthorpe/gopi-hw"
	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>
#include <interface/mmal/util/mmal_util_params.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - PARAMETERS

func MMALPortParameterAllocGet(handle MMAL_PortHandle, name MMAL_ParameterType, size uint32) (MMAL_ParameterHandle, error) {
	var err C.MMAL_STATUS_T
	if param := C.mmal_port_parameter_alloc_get(handle, C.uint32_t(name), C.uint32_t(size), &err); MMAL_Status(err) != MMAL_SUCCESS {
		return nil, MMAL_Status(err)
	} else {
		return param, nil
	}
}

func MMALPortParameterAllocFree(handle MMAL_ParameterHandle) {
	C.mmal_port_parameter_free(handle)
}

func MMALPortParameterSetBool(handle MMAL_PortHandle, name MMAL_ParameterType, value bool) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_boolean(handle, C.uint32_t(name), mmal_to_bool(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetBool(handle MMAL_PortHandle, name MMAL_ParameterType) (bool, error) {
	var value C.MMAL_BOOL_T
	if status := MMAL_Status(C.mmal_port_parameter_get_boolean(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return value != C.MMAL_BOOL_T(0), nil
	} else {
		return false, status
	}
}

func MMALPortParameterSetUint64(handle MMAL_PortHandle, name MMAL_ParameterType, value uint64) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_uint64(handle, C.uint32_t(name), C.uint64_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetUint64(handle MMAL_PortHandle, name MMAL_ParameterType) (uint64, error) {
	var value C.uint64_t
	if status := MMAL_Status(C.mmal_port_parameter_get_uint64(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return uint64(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetInt64(handle MMAL_PortHandle, name MMAL_ParameterType, value int64) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_int64(handle, C.uint32_t(name), C.int64_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetInt64(handle MMAL_PortHandle, name MMAL_ParameterType) (int64, error) {
	var value C.int64_t
	if status := MMAL_Status(C.mmal_port_parameter_get_int64(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return int64(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetUint32(handle MMAL_PortHandle, name MMAL_ParameterType, value uint32) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_uint32(handle, C.uint32_t(name), C.uint32_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetUint32(handle MMAL_PortHandle, name MMAL_ParameterType) (uint32, error) {
	var value C.uint32_t
	if status := MMAL_Status(C.mmal_port_parameter_get_uint32(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return uint32(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetInt32(handle MMAL_PortHandle, name MMAL_ParameterType, value int32) error {
	if status := MMAL_Status(C.mmal_port_parameter_set_int32(handle, C.uint32_t(name), C.int32_t(value))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetInt32(handle MMAL_PortHandle, name MMAL_ParameterType) (int32, error) {
	var value C.int32_t
	if status := MMAL_Status(C.mmal_port_parameter_get_int32(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return int32(value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetString(handle MMAL_PortHandle, name MMAL_ParameterType, value string) error {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	if status := MMAL_Status(C.mmal_port_parameter_set_string(handle, C.uint32_t(name), cValue)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterSetBytes(handle MMAL_PortHandle, name MMAL_ParameterType, value []byte) error {
	ptr := (*C.uint8_t)(unsafe.Pointer(&value[0]))
	len := len(value)
	if status := MMAL_Status(C.mmal_port_parameter_set_bytes(handle, C.uint32_t(name), ptr, C.uint(len))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterSetRational(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALRationalNum) error {
	value_ := C.MMAL_RATIONAL_T{C.int32_t(value.Num), C.int32_t(value.Den)}
	if status := MMAL_Status(C.mmal_port_parameter_set_rational(handle, C.uint32_t(name), value_)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetRational(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALRationalNum, error) {
	var value C.MMAL_RATIONAL_T
	if status := MMAL_Status(C.mmal_port_parameter_get_rational(handle, C.uint32_t(name), &value)); status == MMAL_SUCCESS {
		return hw.MMALRationalNum{int32(value.num), int32(value.den)}, nil
	} else {
		return hw.MMALRationalNum{}, status
	}
}

func MMALPortParameterSetSeek(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_ParameterSeek) error {
	value.hdr.id = C.uint32_t(name)
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, (*C.MMAL_PARAMETER_HEADER_T)(&value.hdr))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterSetDisplayRegion(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_DisplayRegion) error {
	if status := MMAL_Status(C.mmal_util_set_display_region(handle, value)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetDisplayRegion(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_DisplayRegion, error) {
	var value (C.MMAL_DISPLAYREGION_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_DISPLAYREGION_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		display_region := MMAL_DisplayRegion(&value)
		display_region.set = MMAL_DISPLAY_SET_NONE
		return display_region, nil
	} else {
		return nil, status
	}
}

func MMALPortParameterGetVideoProfile(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALVideoProfile, error) {
	var value (C.MMAL_PARAMETER_VIDEO_PROFILE_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_VIDEO_PROFILE_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return hw.MMALVideoProfile{hw.MMALVideoEncProfile(value.profile[0].profile), hw.MMALVideoEncLevel(value.profile[0].level)}, nil
	} else {
		return hw.MMALVideoProfile{}, status
	}
}

func MMALPortParameterSetVideoProfile(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALVideoProfile) error {
	var value_ (C.MMAL_PARAMETER_VIDEO_PROFILE_T)
	value_.hdr.id = C.uint32_t(name)
	value_.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_VIDEO_PROFILE_T{}))
	value_.profile[0].profile = C.MMAL_VIDEO_PROFILE_T(value.Profile)
	value_.profile[0].level = C.MMAL_VIDEO_LEVEL_T(value.Level)
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, &value_.hdr)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALParamGetArrayVideoProfile(handle MMAL_ParameterHandle) []hw.MMALVideoProfile {
	// TODO
	fmt.Println("<TODO> SIZE=", handle.size) //
	return []hw.MMALVideoProfile{}
}

func MMALPortParameterGetCameraInfo(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_CameraInfo, error) {
	var value (C.MMAL_PARAMETER_CAMERA_INFO_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_CAMERA_INFO_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return MMAL_CameraInfo(&value), nil
	} else {
		return nil, status
	}
}

func MMALPortParameterGetCameraMeteringMode(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALCameraMeteringMode, error) {
	return 0, gopi.ErrNotImplemented // TODO
}

func MMALPortParameterSetCameraMeteringMode(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALCameraMeteringMode) error {
	return gopi.ErrNotImplemented // TODO
}

func MMALPortParameterGetCameraExposureMode(handle MMAL_PortHandle, name MMAL_ParameterType) (hw.MMALCameraExposureMode, error) {
	var value (C.MMAL_PARAMETER_EXPOSUREMODE_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_EXPOSUREMODE_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return hw.MMALCameraExposureMode(value.value), nil
	} else {
		return 0, status
	}
}

func MMALPortParameterSetCameraExposureMode(handle MMAL_PortHandle, name MMAL_ParameterType, value hw.MMALCameraExposureMode) error {
	var value_ (C.MMAL_PARAMETER_EXPOSUREMODE_T)
	value_.hdr.id = C.uint32_t(name)
	value_.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_EXPOSUREMODE_T{}))
	value_.value = C.MMAL_PARAM_EXPOSUREMODE_T(value)
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, &value_.hdr)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortParameterGetCameraAnnotation(handle MMAL_PortHandle, name MMAL_ParameterType) (MMAL_CameraAnnotation, error) {
	var value (C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T)
	value.hdr.id = C.uint32_t(name)
	value.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_get(handle, &value.hdr)); status == MMAL_SUCCESS {
		return MMAL_CameraAnnotation(&value), nil
	} else {
		return nil, status
	}
}

func MMALPortParameterSetCameraAnnotation(handle MMAL_PortHandle, name MMAL_ParameterType, value MMAL_CameraAnnotation) error {
	var value_ (C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T)
	value_.hdr.id = C.uint32_t(name)
	value_.hdr.size = C.uint32_t(unsafe.Sizeof(C.MMAL_PARAMETER_CAMERA_ANNOTATE_V4_T{}))
	if status := MMAL_Status(C.mmal_port_parameter_set(handle, &value_.hdr)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALParamGetArrayUint32(handle MMAL_ParameterHandle) []uint32 {
	var array []uint32

	// Data and length of the array
	data := uintptr(unsafe.Pointer(handle)) + unsafe.Sizeof(*handle)
	len := (uintptr(handle.size) - unsafe.Sizeof(*handle)) / C.sizeof_uint32_t

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&array)))
	sliceHeader.Cap = int(len)
	sliceHeader.Len = int(len)
	sliceHeader.Data = data

	// Return the array
	return array
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SEEK

func (this *MMAL_ParameterSeek) SetOffset(value int64) {
	this.offset = C.int64_t(value)
}

func (this *MMAL_ParameterSeek) SetFlags(value uint32) {
	this.flags = C.uint32_t(value)
}
