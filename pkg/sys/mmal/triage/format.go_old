//+build mmal

package mmal

import (
	"unsafe"

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
// PUBLIC METHODS - STREAM FORMATS

func MMALStreamFormatAlloc() MMAL_StreamFormat {
	return MMAL_StreamFormat(C.mmal_format_alloc())
}

func MMALStreamFormatFree(handle MMAL_StreamFormat) {
	C.mmal_format_free(handle)
}

func MMALStreamFormatExtraDataAlloc(handle MMAL_StreamFormat, size uint) error {
	if status := MMAL_Status(C.mmal_format_extradata_alloc(handle, C.uint(size))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALStreamFormatCopy(dest, src MMAL_StreamFormat) error {
	C.mmal_format_copy(dest, src)
	return nil
}

func MMALStreamFormatFullCopy(dest, src MMAL_StreamFormat) error {
	if status := MMAL_Status(C.mmal_format_full_copy(dest, src)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALStreamFormatCompare(dest, src MMAL_StreamFormat) MMAL_StreamCompareFlags {
	return MMAL_StreamCompareFlags(C.mmal_format_compare(dest, src))
}

func MMALStreamFormatType(handle MMAL_StreamFormat) MMAL_StreamType {
	return MMAL_StreamType(handle._type)
}

func MMALStreamFormatBitrate(handle MMAL_StreamFormat) uint32 {
	return uint32(handle.bitrate)
}

func MMALStreamFormatSetBitrate(handle MMAL_StreamFormat, value uint32) {
	handle.bitrate = C.uint32_t(value)
}

func MMALStreamFormatEncoding(handle MMAL_StreamFormat) (hw.MMALEncodingType, hw.MMALEncodingType) {
	return hw.MMALEncodingType(handle.encoding), hw.MMALEncodingType(handle.encoding_variant)
}

func MMALStreamFormatSetEncoding(handle MMAL_StreamFormat, value, variant hw.MMALEncodingType) {
	handle.encoding = C.MMAL_FOURCC_T(value)
	handle.encoding_variant = C.MMAL_FOURCC_T(variant)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION - VIDEO STREAM FORMAT

func MMALStreamFormatVideoWidthHeight(handle MMAL_StreamFormat) (uint32, uint32) {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	return uint32(video.width), uint32(video.height)
}

func MMALStreamFormatVideoSetWidthHeight(handle MMAL_StreamFormat, w, h uint32) {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	video.width = C.uint32_t(w)
	video.height = C.uint32_t(h)
}

func MMALStreamFormatVideoCrop(handle MMAL_StreamFormat) hw.MMALRect {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	return hw.MMALRect{int32(video.crop.x), int32(video.crop.y), uint32(video.crop.width), uint32(video.crop.height)}
}

func MMALStreamFormatVideoSetCrop(handle MMAL_StreamFormat, value hw.MMALRect) {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	video.crop = C.MMAL_RECT_T{C.int32_t(value.X), C.int32_t(value.Y), C.int32_t(value.W), C.int32_t(value.H)}
}

func MMALStreamFormatVideoFrameRate(handle MMAL_StreamFormat) hw.MMALRationalNum {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	return hw.MMALRationalNum{int32(video.frame_rate.num), int32(video.frame_rate.den)}
}

func MMALStreamFormatVideoSetFrameRate(handle MMAL_StreamFormat, value hw.MMALRationalNum) {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	video.frame_rate.num = C.int32_t(value.Num)
	video.frame_rate.den = C.int32_t(value.Den)
}

func MMALStreamFormatVideoPixelAspectRatio(handle MMAL_StreamFormat) hw.MMALRationalNum {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	return hw.MMALRationalNum{int32(video.par.num), int32(video.par.den)}
}

func MMALStreamFormatVideoSetPixelAspectRatio(handle MMAL_StreamFormat, value hw.MMALRationalNum) {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	video.par.num = C.int32_t(value.Num)
	video.par.den = C.int32_t(value.Den)
}

func MMALStreamFormatVideoColorSpace(handle MMAL_StreamFormat) hw.MMALEncodingType {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	return hw.MMALEncodingType(video.color_space)
}

func MMALStreamFormatVideoSetColorSpace(handle MMAL_StreamFormat, value hw.MMALEncodingType) {
	video := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(handle.es))
	video.color_space = C.MMAL_FOURCC_T(value)
}
