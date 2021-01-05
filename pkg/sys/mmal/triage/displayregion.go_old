//+build mmal

package mmal

import (
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
// DISPLAY REGION

func MMALDisplayRegionGetDisplayNum(handle MMAL_DisplayRegion) uint32 {
	return uint32(handle.display_num)
}

func MMALDisplayRegionGetFullScreen(handle MMAL_DisplayRegion) bool {
	return handle.fullscreen != 0
}

func MMALDisplayRegionSetFullScreen(handle MMAL_DisplayRegion, value bool) {
	handle.fullscreen = mmal_to_bool(value)
	handle.set |= MMAL_DISPLAY_SET_FULLSCREEN
}

func MMALDisplayRegionGetLayer(handle MMAL_DisplayRegion) int32 {
	return int32(handle.layer)
}

func MMALDisplayRegionGetAlpha(handle MMAL_DisplayRegion) uint32 {
	return uint32(handle.alpha)
}

func MMALDisplayRegionSetLayer(handle MMAL_DisplayRegion, value int32) {
	handle.layer = C.int32_t(value)
	handle.set |= MMAL_DISPLAY_SET_LAYER
}

func MMALDisplayRegionSetAlpha(handle MMAL_DisplayRegion, value uint32) {
	handle.alpha = C.uint32_t(value)
	handle.set |= MMAL_DISPLAY_SET_ALPHA
}

func MMALDisplayRegionGetTransform(handle MMAL_DisplayRegion) hw.MMALDisplayTransform {
	return hw.MMALDisplayTransform(handle.transform)
}

func MMALDisplayRegionGetMode(handle MMAL_DisplayRegion) hw.MMALDisplayMode {
	return hw.MMALDisplayMode(handle.mode)
}

func MMALDisplayRegionSetTransform(handle MMAL_DisplayRegion, value hw.MMALDisplayTransform) {
	handle.transform = C.MMAL_DISPLAYTRANSFORM_T(value)
	handle.set |= MMAL_DISPLAY_SET_TRANSFORM
}

func MMALDisplayRegionSetMode(handle MMAL_DisplayRegion, value hw.MMALDisplayMode) {
	handle.mode = C.MMAL_DISPLAYMODE_T(value)
	handle.set |= MMAL_DISPLAY_SET_MODE
}

func MMALDisplayRegionGetNoAspect(handle MMAL_DisplayRegion) bool {
	return handle.noaspect != 0
}

func MMALDisplayRegionGetCopyProtect(handle MMAL_DisplayRegion) bool {
	return handle.copyprotect_required != 0
}

func MMALDisplayRegionSetNoAspect(handle MMAL_DisplayRegion, value bool) {
	handle.noaspect = mmal_to_bool(value)
	handle.set |= MMAL_DISPLAY_SET_NOASPECT
}

func MMALDisplayRegionSetCopyProtect(handle MMAL_DisplayRegion, value bool) {
	handle.copyprotect_required = mmal_to_bool(value)
	handle.set |= MMAL_DISPLAY_SET_COPYPROTECT
}

func MMALDisplayRegionGetDestRect(handle MMAL_DisplayRegion) hw.MMALRect {
	return hw.MMALRect{int32(handle.dest_rect.x), int32(handle.dest_rect.y), uint32(handle.dest_rect.width), uint32(handle.dest_rect.height)}
}

func MMALDisplayRegionGetSrcRect(handle MMAL_DisplayRegion) hw.MMALRect {
	return hw.MMALRect{int32(handle.src_rect.x), int32(handle.src_rect.y), uint32(handle.src_rect.width), uint32(handle.src_rect.height)}
}

func MMALDisplayRegionSetDestRect(handle MMAL_DisplayRegion, value hw.MMALRect) {
	handle.dest_rect = C.MMAL_RECT_T{C.int32_t(value.X), C.int32_t(value.Y), C.int32_t(value.W), C.int32_t(value.H)}
	handle.set |= MMAL_DISPLAY_SET_DEST_RECT
}

func MMALDisplayRegionSetSrcRect(handle MMAL_DisplayRegion, value hw.MMALRect) {
	handle.src_rect = C.MMAL_RECT_T{C.int32_t(value.X), C.int32_t(value.Y), C.int32_t(value.W), C.int32_t(value.H)}
	handle.set |= MMAL_DISPLAY_SET_SRC_RECT
}
