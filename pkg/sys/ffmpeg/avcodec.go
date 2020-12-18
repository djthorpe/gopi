// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"strconv"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libavcodec
#include <libavcodec/avcodec.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	AVCodec           C.struct_AVCodec
	AVCodecParameters C.struct_AVCodecParameters
)

////////////////////////////////////////////////////////////////////////////////
// GET CODEC FUNCTIONS

// AllCodecs returns all registered codecs
func AllCodecs() []*AVCodec {
	codecs := make([]*AVCodec, 0)
	ptr := unsafe.Pointer(nil)
	for {
		if codec := C.av_codec_iterate(&ptr); codec == nil {
			break
		} else {
			codecs = append(codecs, (*AVCodec)(codec))
		}
	}
	return codecs
}

func FindCodecById(id AVCodecId) *AVCodec {
	ptr := unsafe.Pointer(nil)
	for {
		if codec := C.av_codec_iterate(&ptr); codec == nil {
			break
		} else if AVCodecId(codec.id) == id {
			return (*AVCodec)(codec)
		}
	}
	return nil
}

func FindCodecByName(name string) *AVCodec {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	ptr := unsafe.Pointer(nil)
	for {
		if codec := C.av_codec_iterate(&ptr); codec == nil {
			break
		} else if C.strcmp(name_, codec.name) == 0 {
			return (*AVCodec)(codec)
		}
	}
	return nil
}

func FindDecoderById(id AVCodecId) *AVCodec {
	if codec := C.avcodec_find_decoder(uint32(id)); codec == nil {
		return nil
	} else {
		return (*AVCodec)(codec)
	}
}

func FindDecoderByName(name string) *AVCodec {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	if codec := C.avcodec_find_decoder_by_name(name_); codec == nil {
		return nil
	} else {
		return (*AVCodec)(codec)
	}
}

func FindEncoderById(id AVCodecId) *AVCodec {
	if codec := C.avcodec_find_encoder(uint32(id)); codec == nil {
		return nil
	} else {
		return (*AVCodec)(codec)
	}
}

func FindEncoderByName(name string) *AVCodec {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	if codec := C.avcodec_find_encoder_by_name(name_); codec == nil {
		return nil
	} else {
		return (*AVCodec)(codec)
	}
}

////////////////////////////////////////////////////////////////////////////////
// AVCodec

func (this *AVCodec) Name() string {
	return C.GoString(this.name)
}

func (this *AVCodec) Description() string {
	return C.GoString(this.long_name)
}

func (this *AVCodec) Type() AVMediaType {
	return AVMediaType(this._type)
}

func (this *AVCodec) Id() AVCodecId {
	return AVCodecId(this.id)
}

func (this *AVCodec) Capabilities() AVCodecCap {
	return AVCodecCap(this.capabilities)
}

func (this *AVCodec) WrapperName() string {
	return C.GoString(this.wrapper_name)
}

func (this *AVCodec) IsEncoder() bool {
	ctx := (*C.AVCodec)(this)
	if C.av_codec_is_encoder(ctx) == 0 {
		return false
	} else {
		return true
	}
}

func (this *AVCodec) IsDecoder() bool {
	ctx := (*C.AVCodec)(this)
	if C.av_codec_is_decoder(ctx) == 0 {
		return false
	} else {
		return true
	}
}

func (this *AVCodec) String() string {
	str := "<AVCodec"
	str += " name=" + strconv.Quote(this.Name())
	str += " description=" + strconv.Quote(this.Description())
	str += " type=" + fmt.Sprint(this.Type())
	if cap := this.Capabilities(); cap != AV_CODEC_CAP_NONE {
		str += " capabilities=" + fmt.Sprint(this.Capabilities())
	}
	if wn := this.WrapperName(); wn != "" {
		str += " wrapper_name=" + strconv.Quote(wn)
	}
	if encoder := this.IsEncoder(); encoder {
		str += " encoder=true"
	}
	if decoder := this.IsDecoder(); decoder {
		str += " decoder=true"
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// AVCODECPARAMETERS

// NewAVCodecParameters allocates a new AVCodecParameters and set
// its fields to default values (unknown/invalid/0)
func NewAVCodecParameters() *AVCodecParameters {
	return (*AVCodecParameters)(C.avcodec_parameters_alloc())
}

// Free AVCodecParameters
func (this *AVCodecParameters) Free() {
	ctx := (*C.AVCodecParameters)(unsafe.Pointer(this))
	C.avcodec_parameters_free(&ctx)
}

// Create a new Codec decoder context
func (this *AVCodecParameters) NewDecoderContext() (*AVCodecContext, *AVCodec) {
	if codec := FindDecoderById(this.Id()); codec == nil {
		return nil, nil
	} else if ctx := NewAVCodecContext(codec); ctx == nil {
		return nil, nil
	} else {
		return ctx, codec
	}
}

// FromContext fill the parameters based on the values from the
// supplied codec context
func (this *AVCodecParameters) FromContext(other *AVCodecContext) error {
	ctx := (*C.AVCodecParameters)(this)
	if err := AVError(C.avcodec_parameters_from_context(ctx, (*C.AVCodecContext)(other))); err != 0 {
		return err
	} else {
		return nil
	}
}

// ToContext fills the codec context based on the values
func (this *AVCodecParameters) ToContext(other *AVCodecContext) error {
	ctx := (*C.AVCodecParameters)(this)
	if err := AVError(C.avcodec_parameters_to_context((*C.AVCodecContext)(other), ctx)); err != 0 {
		return err
	} else {
		return nil
	}
}

// From fill the parameters based on the values from the supplied codec parameters
func (this *AVCodecParameters) CopyFrom(codecpar *AVCodecParameters) error {
	ctx := (*C.AVCodecParameters)(this)
	if err := AVError(C.avcodec_parameters_copy(ctx, (*C.AVCodecParameters)(codecpar))); err != 0 {
		return err
	} else {
		return nil
	}
}

func (this *AVCodecParameters) Type() AVMediaType {
	return AVMediaType(this.codec_type)
}

func (this *AVCodecParameters) Id() AVCodecId {
	return AVCodecId(this.codec_id)
}

func (this *AVCodecParameters) Tag() uint32 {
	return uint32(this.codec_tag)
}

func (this *AVCodecParameters) BitRate() int32 {
	return int32(this.bit_rate)
}

func (this *AVCodecParameters) Width() uint {
	return uint(this.width)
}

func (this *AVCodecParameters) Height() uint {
	return uint(this.height)
}

func (this *AVCodecParameters) String() string {
	str := "<AVCodecParameters"
	str += " type=" + fmt.Sprint(this.Type())
	str += " id=" + fmt.Sprint(this.Id())
	str += " tag=" + fmt.Sprintf("0x%08X", this.Tag())
	if br := this.BitRate(); br != 0 {
		str += " bit_rate=" + fmt.Sprint(br)
	}
	if w, h := this.Width(), this.Height(); w != 0 && h != 0 {
		str += " w,h={ " + fmt.Sprint(w, ",", h) + " }"
	}
	return str + ">"
}
