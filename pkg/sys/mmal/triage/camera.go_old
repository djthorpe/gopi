//+build mmal

package mmal

import (
	"unsafe"

	// Frameworks
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// CAMERA ANNOTATION

func MMALCameraAnnotationEnabled(handle MMAL_CameraAnnotation) bool {
	return handle.enable == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetEnabled(handle MMAL_CameraAnnotation, value bool) {
	handle.enable = mmal_to_bool(value)
}

func MMALCameraAnnotationShowShutter(handle MMAL_CameraAnnotation) bool {
	return handle.show_shutter == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowShutter(handle MMAL_CameraAnnotation, value bool) {
	handle.show_shutter = mmal_to_bool(value)
}

func MMALCameraAnnotationShowAnalogGain(handle MMAL_CameraAnnotation) bool {
	return handle.show_analog_gain == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowAnalogGain(handle MMAL_CameraAnnotation, value bool) {
	handle.show_analog_gain = mmal_to_bool(value)
}

func MMALCameraAnnotationShowLens(handle MMAL_CameraAnnotation) bool {
	return handle.show_lens == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowLens(handle MMAL_CameraAnnotation, value bool) {
	handle.show_lens = mmal_to_bool(value)
}

func MMALCameraAnnotationShowCAF(handle MMAL_CameraAnnotation) bool {
	return handle.show_caf == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowCAF(handle MMAL_CameraAnnotation, value bool) {
	handle.show_caf = mmal_to_bool(value)
}

func MMALCameraAnnotationShowMotion(handle MMAL_CameraAnnotation) bool {
	return handle.show_motion == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowMotion(handle MMAL_CameraAnnotation, value bool) {
	handle.show_motion = mmal_to_bool(value)
}

func MMALCameraAnnotationShowFrameNum(handle MMAL_CameraAnnotation) bool {
	return handle.show_frame_num == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowFrameNum(handle MMAL_CameraAnnotation, value bool) {
	handle.show_frame_num = mmal_to_bool(value)
}

func MMALCameraAnnotationShowTextBackground(handle MMAL_CameraAnnotation) bool {
	return handle.enable_text_background == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetShowTextBackground(handle MMAL_CameraAnnotation, value bool) {
	handle.enable_text_background = mmal_to_bool(value)
}

func MMALCameraAnnotationUseCustomBackgroundColor(handle MMAL_CameraAnnotation) bool {
	return handle.custom_background_colour == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetUseCustomBackgroundColor(handle MMAL_CameraAnnotation, value bool) {
	handle.custom_background_colour = mmal_to_bool(value)
}

func MMALCameraAnnotationUseCustomColor(handle MMAL_CameraAnnotation) bool {
	return handle.custom_text_colour == MMAL_BOOL_TRUE
}

func MMALCameraAnnotationSetUseCustomColor(handle MMAL_CameraAnnotation, value bool) {
	handle.custom_text_colour = mmal_to_bool(value)
}

func MMALCameraAnnotationText(handle MMAL_CameraAnnotation) string {
	return C.GoString(&handle.text[0])
}

func MMALCameraAnnotationSetText(handle MMAL_CameraAnnotation, value string) {
	cstr := C.CString(value)
	defer C.free(unsafe.Pointer(cstr))
	C.strncpy(&handle.text[0], cstr, C.uint(len(value)))
}

func MMALCameraAnnotationTextSize(handle MMAL_CameraAnnotation) uint8 {
	return uint8(handle.text_size)
}

func MMALCameraAnnotationSetTextSize(handle MMAL_CameraAnnotation, value uint8) {
	handle.text_size = C.uint8_t(value)
}

////////////////////////////////////////////////////////////////////////////////
// CAMERA INFO

func MMALCameraInfoGetCamerasNum(handle MMAL_CameraInfo) uint32 {
	return uint32(handle.num_cameras)
}

func MMALCameraInfoGetFlashesNum(handle MMAL_CameraInfo) uint32 {
	return uint32(handle.num_flashes)
}

func MMALCameraInfoGetCameras(handle MMAL_CameraInfo) []MMAL_Camera {
	cameras := make([]MMAL_Camera, int(handle.num_cameras))
	for i := 0; i < len(cameras); i++ {
		cameras[i] = MMAL_Camera(handle.cameras[i])
	}
	return cameras
}

func MMALCameraInfoGetFlashes(handle MMAL_CameraInfo) []hw.MMALCameraFlashType {
	flashes := make([]hw.MMALCameraFlashType, int(handle.num_flashes))
	for i := 0; i < len(flashes); i++ {
		flashes[i] = hw.MMALCameraFlashType(handle.flashes[i].flash_type)
	}
	return flashes
}

func MMALCameraInfoGetCameraId(handle MMAL_Camera) uint32 {
	return uint32(handle.port_id)
}

func MMALCameraInfoGetCameraName(handle MMAL_Camera) string {
	return C.GoString(&handle.camera_name[0])
}

func MMALCameraInfoGetCameraMaxWidth(handle MMAL_Camera) uint32 {
	return uint32(handle.max_width)
}
func MMALCameraInfoGetCameraMaxHeight(handle MMAL_Camera) uint32 {
	return uint32(handle.max_height)
}

func MMALCameraInfoGetCameraLensPresent(handle MMAL_Camera) bool {
	return handle.lens_present == MMAL_BOOL_TRUE
}
