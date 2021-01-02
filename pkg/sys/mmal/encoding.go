//+build mmal

package mmal

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALEncodingType uint32
	MMALColorSpace   uint32
	MMALBufferEvent  uint32
)

var (
	MMAL_ENCODING_UNKNOWN MMALEncodingType = C.MMAL_ENCODING_UNKNOWN
)

////////////////////////////////////////////////////////////////////////////////
// VIDEO ENCODINGS

var (
	MMAL_ENCODING_H264   MMALEncodingType = C.MMAL_ENCODING_H264
	MMAL_ENCODING_MVC    MMALEncodingType = C.MMAL_ENCODING_MVC
	MMAL_ENCODING_H263   MMALEncodingType = C.MMAL_ENCODING_H263
	MMAL_ENCODING_MP4V   MMALEncodingType = C.MMAL_ENCODING_MP4V
	MMAL_ENCODING_MP2V   MMALEncodingType = C.MMAL_ENCODING_MP2V
	MMAL_ENCODING_MP1V   MMALEncodingType = C.MMAL_ENCODING_MP1V
	MMAL_ENCODING_WMV3   MMALEncodingType = C.MMAL_ENCODING_WMV3
	MMAL_ENCODING_WMV2   MMALEncodingType = C.MMAL_ENCODING_WMV2
	MMAL_ENCODING_WMV1   MMALEncodingType = C.MMAL_ENCODING_WMV1
	MMAL_ENCODING_WVC1   MMALEncodingType = C.MMAL_ENCODING_WVC1
	MMAL_ENCODING_VP8    MMALEncodingType = C.MMAL_ENCODING_VP8
	MMAL_ENCODING_VP7    MMALEncodingType = C.MMAL_ENCODING_VP7
	MMAL_ENCODING_VP6    MMALEncodingType = C.MMAL_ENCODING_VP6
	MMAL_ENCODING_THEORA MMALEncodingType = C.MMAL_ENCODING_THEORA
	MMAL_ENCODING_SPARK  MMALEncodingType = C.MMAL_ENCODING_SPARK
	MMAL_ENCODING_MJPEG  MMALEncodingType = C.MMAL_ENCODING_MJPEG
)

////////////////////////////////////////////////////////////////////////////////
// IMAGE ENCODINGS

var (
	MMAL_ENCODING_JPEG MMALEncodingType = C.MMAL_ENCODING_JPEG
	MMAL_ENCODING_GIF  MMALEncodingType = C.MMAL_ENCODING_GIF
	MMAL_ENCODING_PNG  MMALEncodingType = C.MMAL_ENCODING_PNG
	MMAL_ENCODING_PPM  MMALEncodingType = C.MMAL_ENCODING_PPM
	MMAL_ENCODING_TGA  MMALEncodingType = C.MMAL_ENCODING_TGA
	MMAL_ENCODING_BMP  MMALEncodingType = C.MMAL_ENCODING_BMP
)

////////////////////////////////////////////////////////////////////////////////
// UNCOMPRESSED ENCODINGS

var (
	MMAL_ENCODING_I420        MMALEncodingType = C.MMAL_ENCODING_I420
	MMAL_ENCODING_I420_SLICE  MMALEncodingType = C.MMAL_ENCODING_I420_SLICE
	MMAL_ENCODING_YV12        MMALEncodingType = C.MMAL_ENCODING_YV12
	MMAL_ENCODING_I422        MMALEncodingType = C.MMAL_ENCODING_I422
	MMAL_ENCODING_I422_SLICE  MMALEncodingType = C.MMAL_ENCODING_I422_SLICE
	MMAL_ENCODING_YUYV        MMALEncodingType = C.MMAL_ENCODING_YUYV
	MMAL_ENCODING_YVYU        MMALEncodingType = C.MMAL_ENCODING_YVYU
	MMAL_ENCODING_UYVY        MMALEncodingType = C.MMAL_ENCODING_UYVY
	MMAL_ENCODING_VYUY        MMALEncodingType = C.MMAL_ENCODING_VYUY
	MMAL_ENCODING_NV12        MMALEncodingType = C.MMAL_ENCODING_NV12
	MMAL_ENCODING_NV21        MMALEncodingType = C.MMAL_ENCODING_NV21
	MMAL_ENCODING_ARGB        MMALEncodingType = C.MMAL_ENCODING_ARGB
	MMAL_ENCODING_ARGB_SLICE  MMALEncodingType = C.MMAL_ENCODING_ARGB_SLICE
	MMAL_ENCODING_RGBA        MMALEncodingType = C.MMAL_ENCODING_RGBA
	MMAL_ENCODING_RGBA_SLICE  MMALEncodingType = C.MMAL_ENCODING_RGBA_SLICE
	MMAL_ENCODING_ABGR        MMALEncodingType = C.MMAL_ENCODING_ABGR
	MMAL_ENCODING_ABGR_SLICE  MMALEncodingType = C.MMAL_ENCODING_ABGR_SLICE
	MMAL_ENCODING_BGRA        MMALEncodingType = C.MMAL_ENCODING_BGRA
	MMAL_ENCODING_BGRA_SLICE  MMALEncodingType = C.MMAL_ENCODING_BGRA_SLICE
	MMAL_ENCODING_RGB16       MMALEncodingType = C.MMAL_ENCODING_RGB16
	MMAL_ENCODING_RGB16_SLICE MMALEncodingType = C.MMAL_ENCODING_RGB16_SLICE
	MMAL_ENCODING_RGB24       MMALEncodingType = C.MMAL_ENCODING_RGB24
	MMAL_ENCODING_RGB24_SLICE MMALEncodingType = C.MMAL_ENCODING_RGB24_SLICE
	MMAL_ENCODING_RGB32       MMALEncodingType = C.MMAL_ENCODING_RGB32
	MMAL_ENCODING_RGB32_SLICE MMALEncodingType = C.MMAL_ENCODING_RGB32_SLICE
	MMAL_ENCODING_BGR16       MMALEncodingType = C.MMAL_ENCODING_BGR16
	MMAL_ENCODING_BGR16_SLICE MMALEncodingType = C.MMAL_ENCODING_BGR16_SLICE
	MMAL_ENCODING_BGR24       MMALEncodingType = C.MMAL_ENCODING_BGR24
	MMAL_ENCODING_BGR24_SLICE MMALEncodingType = C.MMAL_ENCODING_BGR24_SLICE
	MMAL_ENCODING_BGR32       MMALEncodingType = C.MMAL_ENCODING_BGR32
	MMAL_ENCODING_BGR32_SLICE MMALEncodingType = C.MMAL_ENCODING_BGR32_SLICE
)

////////////////////////////////////////////////////////////////////////////////
// BUFFER EVENT

var (
	// Generic Error
	MMAL_EVENT_ERROR MMALBufferEvent = C.MMAL_EVENT_ERROR

	// End-of-stream event. Data contains a MMAL_EVENT_END_OF_STREAM_T
	MMAL_EVENT_EOS MMALBufferEvent = C.MMAL_EVENT_EOS

	// Format changed event. Data contains a MMAL_EVENT_FORMAT_CHANGED_T
	MMAL_EVENT_FORMAT_CHANGED MMALBufferEvent = C.MMAL_EVENT_FORMAT_CHANGED

	// Parameter changed event. Data contains a MMAL_EVENT_PARAMETER_CHANGED_T
	MMAL_EVENT_PARAMETER_CHANGED MMALBufferEvent = C.MMAL_EVENT_PARAMETER_CHANGED
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (e MMALEncodingType) StrideToWidth(stride uint32) uint32 {
	return uint32(C.mmal_encoding_stride_to_width(C.uint32_t(e), C.uint32_t(stride)))
}

func (e MMALEncodingType) WidthToStride(width uint32) uint32 {
	return uint32(C.mmal_encoding_width_to_stride(C.uint32_t(e), C.uint32_t(width)))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e MMALEncodingType) String() string {
	buf := [5]C.char{}
	size := C.size_t(len(buf))
	return C.GoString(C.mmal_4cc_to_string(&buf[0], size, C.uint32_t(e)))
}

func (c MMALColorSpace) String() string {
	buf := [5]C.char{}
	size := C.size_t(len(buf))
	return C.GoString(C.mmal_4cc_to_string(&buf[0], size, C.uint32_t(c)))
}

func (e MMALBufferEvent) String() string {
	switch e {
	case MMAL_EVENT_ERROR:
		return "MMAL_EVENT_ERROR"
	case MMAL_EVENT_EOS:
		return "MMAL_EVENT_EOS"
	case MMAL_EVENT_FORMAT_CHANGED:
		return "MMAL_EVENT_FORMAT_CHANGED"
	case MMAL_EVENT_PARAMETER_CHANGED:
		return "MMAL_EVENT_PARAMETER_CHANGED"
	default:
		return "[?? Invalid MMALBufferEvent value]"
	}
}
