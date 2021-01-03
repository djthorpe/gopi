//+build mmal

package mmal

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALStreamFormat     (C.MMAL_ES_FORMAT_T)
	MMALStreamType       (C.MMAL_ES_TYPE_T)
	MMALStreamFlags      (C.uint32_t)
	MMALVideoFormat      (C.MMAL_VIDEO_FORMAT_T)
	MMALAudioFormat      (C.MMAL_AUDIO_FORMAT_T)
	MMALSubpictureFormat (C.MMAL_SUBPICTURE_FORMAT_T)
	MMALStreamFormatEvent (C.MMAL_EVENT_FORMAT_CHANGED_T)
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_ES_TYPE_UNKNOWN    MMALStreamType = C.MMAL_ES_TYPE_UNKNOWN    // Unknown elementary stream type
	MMAL_ES_TYPE_CONTROL    MMALStreamType = C.MMAL_ES_TYPE_CONTROL    // Elementary stream of control commands
	MMAL_ES_TYPE_AUDIO      MMALStreamType = C.MMAL_ES_TYPE_AUDIO      // Audio elementary stream
	MMAL_ES_TYPE_VIDEO      MMALStreamType = C.MMAL_ES_TYPE_VIDEO      //  Video elementary stream
	MMAL_ES_TYPE_SUBPICTURE MMALStreamType = C.MMAL_ES_TYPE_SUBPICTURE // Sub-picture elementary stream (e.g. subtitles, overlays)
)

const (
	MMAL_ES_FORMAT_FLAG_FRAMED MMALStreamFlags = C.MMAL_ES_FORMAT_FLAG_FRAMED
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func MMALCreateStreamFormat() *MMALStreamFormat {
	ctx := C.mmal_format_alloc()
	return (*MMALStreamFormat)(ctx)
}

func (this *MMALStreamFormat) Free() {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	C.mmal_format_free(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *MMALStreamFormat) ExtradataAlloc(size uint32) error {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	if err := Error(C.mmal_format_extradata_alloc(ctx, C.uint32_t(size))); err == MMAL_SUCCESS {
		return nil
	} else {
		return err
	}
}

func (this *MMALStreamFormat) ExtradataRead(r io.ReadSeeker, size uint32) error {
	if err := this.ExtradataAlloc(size); err != nil {
		return err
	} else if _, err := r.Read(this.Extradata()); err != nil {
		return err
	} else if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES - STREAM

func (this *MMALStreamFormat) Type() MMALStreamType {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	return MMALStreamType(ctx._type)
}

func (this *MMALStreamFormat) SetType(value MMALStreamType) {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	ctx._type = C.MMAL_ES_TYPE_T(value)
}

func (this *MMALStreamFormat) Encoding() MMALEncodingType {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	return MMALEncodingType(ctx.encoding)
}

func (this *MMALStreamFormat) SetEncoding(value MMALEncodingType) {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	ctx.encoding = C.MMAL_FOURCC_T(value)
}

func (this *MMALStreamFormat) Variant() MMALEncodingType {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	return MMALEncodingType(ctx.encoding_variant)
}

func (this *MMALStreamFormat) SetVariant(value MMALEncodingType) {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	ctx.encoding_variant = C.MMAL_FOURCC_T(value)
}

func (this *MMALStreamFormat) Flags() MMALStreamFlags {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	return MMALStreamFlags(ctx.flags)
}

func (this *MMALStreamFormat) SetFlags(f MMALStreamFlags) {
	ctx := (*C.MMAL_ES_FORMAT_T)(this)
	ctx.flags = C.uint32_t(f)
}

func (this *MMALStreamFormat) Video() *MMALVideoFormat {
	if this.Type() == MMAL_ES_TYPE_VIDEO {
		ctx := (*C.MMAL_VIDEO_FORMAT_T)(unsafe.Pointer(this.es))
		return (*MMALVideoFormat)(ctx)
	} else {
		return nil
	}
}

func (this *MMALStreamFormat) Audio() *MMALAudioFormat {
	if this.Type() == MMAL_ES_TYPE_AUDIO {
		ctx := (*C.MMAL_AUDIO_FORMAT_T)(unsafe.Pointer(this.es))
		return (*MMALAudioFormat)(ctx)
	} else {
		return nil
	}
}

func (this *MMALStreamFormat) Subpicture() *MMALSubpictureFormat {
	if this.Type() == MMAL_ES_TYPE_SUBPICTURE {
		ctx := (*C.MMAL_SUBPICTURE_FORMAT_T)(unsafe.Pointer(this.es))
		return (*MMALSubpictureFormat)(ctx)
	} else {
		return nil
	}
}

func (this *MMALStreamFormat) Extradata() []byte {
	var result []byte

	ctx := (*C.MMAL_ES_FORMAT_T)(this)

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.extradata_size)
	sliceHeader.Len = int(ctx.extradata_size)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.extradata))

	// Return data
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES - VIDEO

func (this *MMALVideoFormat) Size() (uint32, uint32) {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	return uint32(ctx.width), uint32(ctx.height)
}

func (this *MMALVideoFormat) SetSize(w, h uint32) {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	ctx.width = C.uint32_t(w)
	ctx.height = C.uint32_t(h)
}

func (this *MMALVideoFormat) Crop() MMALRect {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	return MMALRect(ctx.crop)
}

func (this *MMALVideoFormat) SetCrop(crop MMALRect) {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	ctx.crop = (C.MMAL_RECT_T)(crop)
}

func (this *MMALVideoFormat) FrameRate() MMALRational {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	return MMALRational(ctx.frame_rate)
}

func (this *MMALVideoFormat) SetFrameRate(r MMALRational) {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	ctx.frame_rate = (C.MMAL_RATIONAL_T)(r)
}

func (this *MMALVideoFormat) Par() MMALRational {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	return MMALRational(ctx.par)
}

func (this *MMALVideoFormat) SetPar(r MMALRational) {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	ctx.par = (C.MMAL_RATIONAL_T)(r)
}

func (this *MMALVideoFormat) ColorSpace() MMALColorSpace {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	return MMALColorSpace(ctx.color_space)
}

func (this *MMALVideoFormat) SetColorSpace(cs MMALColorSpace) {
	ctx := (*C.MMAL_VIDEO_FORMAT_T)(this)
	ctx.color_space = (C.MMAL_FOURCC_T)(cs)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES - AUDIO

func (this *MMALAudioFormat) Channels() uint32 {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	return uint32(ctx.channels)
}

func (this *MMALAudioFormat) SetChannels(ch uint32) {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	ctx.channels = C.uint32_t(ch)
}

func (this *MMALAudioFormat) SampleRate() uint32 {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	return uint32(ctx.sample_rate)
}

func (this *MMALAudioFormat) SetSampleRate(sr uint32) {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	ctx.sample_rate = C.uint32_t(sr)
}

func (this *MMALAudioFormat) BitsPerSample() uint32 {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	return uint32(ctx.bits_per_sample)
}

func (this *MMALAudioFormat) SetBitsPerSample(bits uint32) {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	ctx.bits_per_sample = C.uint32_t(bits)
}

func (this *MMALAudioFormat) BlockAlign() uint32 {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	return uint32(ctx.block_align)
}

func (this *MMALAudioFormat) SetBlockAlign(ba uint32) {
	ctx := (*C.MMAL_AUDIO_FORMAT_T)(this)
	ctx.block_align = C.uint32_t(ba)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES - SUBPICTURE

func (this *MMALSubpictureFormat) Offset() (uint32, uint32) {
	ctx := (*C.MMAL_SUBPICTURE_FORMAT_T)(this)
	return uint32(ctx.x_offset), uint32(ctx.y_offset)
}


////////////////////////////////////////////////////////////////////////////////
// PROPERTIES - EVENT

func (this *MMALStreamFormatEvent) BufferMin() (uint32,uint32) {
	ctx := (*C.MMAL_EVENT_FORMAT_CHANGED_T)(this)
	return uint32(ctx.buffer_num_min), uint32(ctx.buffer_size_min)
}

func (this *MMALStreamFormatEvent) BufferPreferred() (uint32,uint32) {
	ctx := (*C.MMAL_EVENT_FORMAT_CHANGED_T)(this)
	return uint32(ctx.buffer_num_recommended), uint32(ctx.buffer_size_recommended)	
}

func (this *MMALStreamFormatEvent) Format() (*MMALStreamFormat) {
	ctx := (*C.MMAL_EVENT_FORMAT_CHANGED_T)(this)
	return (*MMALStreamFormat)(ctx.format)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *MMALStreamFormat) String() string {
	str := "<mmal.format"
	if t := this.Type(); t != 0 {
		str += " type=" + fmt.Sprint(t)
	}
	if enc := this.Encoding(); enc != 0 {
		str += " enc=" + fmt.Sprintf("%q",enc)
		if encvar := this.Variant(); encvar != 0 {
			str += ",var=" + fmt.Sprintf("%q",encvar)
		}
	}
	if f := this.Flags(); f != 0 {
		str += " flags=" + fmt.Sprint(f)
	}
	switch this.Type() {
	case MMAL_ES_TYPE_VIDEO:
		if w, h := this.Video().Size(); w > 0 && h > 0 {
			str += fmt.Sprintf(" size={ %d,%d }", w, h)
		}
		if crop := this.Video().Crop(); crop.IsZero() == false {
			str += " crop=" + fmt.Sprint(crop)
		}
		if r := this.Video().FrameRate(); r.IsZero() == false {
			str += " frame_rate=" + fmt.Sprint(r)
		}
		if p := this.Video().Par(); p.IsZero() == false {
			str += " par=" + fmt.Sprint(p)
		}
		if cs := this.Video().ColorSpace(); cs != 0 {
			str += " color_space=" + fmt.Sprint(cs)
		}
	case MMAL_ES_TYPE_AUDIO:
	case MMAL_ES_TYPE_SUBPICTURE:
	}
	return str + ">"
}

func (s MMALStreamType) String() string {
	switch s {
	case MMAL_ES_TYPE_UNKNOWN:
		return "MMAL_ES_TYPE_UNKNOWN"
	case MMAL_ES_TYPE_CONTROL:
		return "MMAL_ES_TYPE_CONTROL"
	case MMAL_ES_TYPE_AUDIO:
		return "MMAL_ES_TYPE_AUDIO"
	case MMAL_ES_TYPE_VIDEO:
		return "MMAL_ES_TYPE_VIDEO"
	case MMAL_ES_TYPE_SUBPICTURE:
		return "MMAL_ES_TYPE_SUBPICTURE"
	default:
		return "[?? Invalid MMALStreamType value]"
	}
}

func (f MMALStreamFlags) String() string {
	switch f {
	case MMAL_ES_FORMAT_FLAG_FRAMED:
		return "MMAL_ES_FORMAT_FLAG_FRAMED"
	default:
		return "[?? Invalid MMALStreamFlags value]"
	}
}


func (this *MMALStreamFormatEvent) String() string {
	str := "<mmal.formatchanged"
	if n,s := this.BufferMin(); n >0 && s > 0 {
		str += fmt.Sprintf(" min={%v,%v}",n,s)
	}
	if n,s := this.BufferPreferred(); n >0 && s > 0 {
		str += fmt.Sprintf(" preferred={%v,%v}",n,s)
	}
	str += fmt.Sprintf(" format=%v",this.Format())
	return str + ">"
}