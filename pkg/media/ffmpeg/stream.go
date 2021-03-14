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

////////////////////////////////////////////////////////////////////////////////
// INIT

// NewStream returns a stream object, can optiomally copy codec
// parameters from another stream if source is not set to nil
func NewStream(ctx *ffmpeg.AVStream, source *stream) *stream {
	if ctx == nil {
		return nil
	}
	if source == nil {
		if codec := NewCodecWithParameters(ctx.CodecPar()); codec == nil {
			return nil
		} else {
			return &stream{ctx, codec}
		}
	} else {
		if codec := NewCodecWithParameters(source.ctx.CodecPar()); codec == nil {
			return nil
		} else {
			return &stream{ctx, codec}
		}
	}
}

func (this *stream) Release() {
	this.ctx = nil
	this.codec = nil
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

	// Remove encoder/decoder flags
	flags ^= (gopi.MEDIA_FLAG_ENCODER | gopi.MEDIA_FLAG_DECODER)

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
// STRINGIFY

func (this *stream) String() string {
	str := "<ffmpeg.stream"
	str += " index=" + fmt.Sprint(this.Index())
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	str += " codec=" + fmt.Sprint(this.Codec())
	return str + ">"
}
