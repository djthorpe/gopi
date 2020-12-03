// +build ffmpeg

package ffmpeg

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type stream struct {
	ctx   *ffmpeg.AVStream
	codec *codec
}

type codec struct {
	ctx *ffmpeg.AVCodecParameters
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewStream(ctx *ffmpeg.AVStream) *stream {
	if ctx == nil {
		return nil
	} else if codec := NewCodec(ctx.CodecPar()); codec == nil {
		return nil
	} else {
		return &stream{ctx, codec}
	}
}

func NewCodec(ctx *ffmpeg.AVCodecParameters) *codec {
	if ctx == nil {
		return nil
	} else {
		return &codec{ctx}
	}
}

func (this *stream) Release() {
	this.ctx = nil
	this.codec = nil
}

func (this *codec) Release() {
	this.ctx = nil
}

////////////////////////////////////////////////////////////////////////////////
// METHODS - STREAM

func (this *stream) Index() int {
	if this.ctx == nil {
		return -1
	} else {
		return this.ctx.Index()
	}
}

func (this *stream) Flags() gopi.MediaFlag {
	flags := gopi.MEDIA_FLAG_NONE

	// Return NONE if released
	if this.ctx == nil {
		return flags
	}

	// Codec flags
	flags |= this.codec.Flags()

	// Disposition flags
	if this.ctx.Disposition()&ffmpeg.AV_DISPOSITION_ATTACHED_PIC != 0 {
		flags |= gopi.MEDIA_FLAG_ARTWORK
	}
	if this.ctx.Disposition()&ffmpeg.AV_DISPOSITION_CAPTIONS != 0 {
		flags |= gopi.MEDIA_FLAG_CAPTIONS
	}

	// Return flags
	return flags
}

func (this *stream) Codec() gopi.MediaCodec {
	if this.ctx == nil {
		return nil
	} else {
		return this.codec
	}
}

func (this *stream) NewContextWithOptions(options *ffmpeg.AVDictionary) *ffmpeg.AVCodecContext {
	if this.ctx == nil || this.codec == nil {
		return nil
	} else if ctx, codec := this.codec.ctx.NewDecoderContext(); ctx == nil || codec == nil {
		return nil
	} else if err := this.codec.ctx.ToContext(ctx); err != nil {
		ctx.Free()
		return nil
	} else if err := ctx.Open(codec, options); err != nil {
		ctx.Free()
		return nil
	} else {
		return ctx
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS - CODEC

func (this *codec) Flags() gopi.MediaFlag {
	flags := gopi.MEDIA_FLAG_NONE
	if this.ctx == nil {
		return flags
	}

	switch this.ctx.Type() {
	case ffmpeg.AVMEDIA_TYPE_VIDEO:
		if this.ctx.BitRate() > 0 {
			flags |= gopi.MEDIA_FLAG_VIDEO
		}
	case ffmpeg.AVMEDIA_TYPE_AUDIO:
		flags |= gopi.MEDIA_FLAG_AUDIO
	case ffmpeg.AVMEDIA_TYPE_UNKNOWN, ffmpeg.AVMEDIA_TYPE_DATA:
		flags |= gopi.MEDIA_FLAG_DATA
	case ffmpeg.AVMEDIA_TYPE_ATTACHMENT:
		flags |= gopi.MEDIA_FLAG_ATTACHMENT
	}

	// Return flags
	return flags
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stream) String() string {
	str := "<stream"
	str += " index=" + fmt.Sprint(this.Index())
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	str += " codec=" + fmt.Sprint(this.Codec())
	return str + ">"
}

func (this *codec) String() string {
	str := "<codec"
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	return str + ">"
}