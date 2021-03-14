// +build ffmpeg

package ffmpeg

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type codec struct {
	ctx   *ffmpeg.AVCodecParameters
	codec *ffmpeg.AVCodec
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewCodec(ctx *ffmpeg.AVCodec) *codec {
	if ctx == nil {
		return nil
	} else {
		return &codec{nil, ctx}
	}
}

func NewCodecWithParameters(ctx *ffmpeg.AVCodecParameters) *codec {
	if ctx == nil {
		return nil
	} else if c := ffmpeg.FindCodecById(ctx.Id()); c == nil {
		return nil
	} else {
		return &codec{ctx, c}
	}
}

func (this *codec) Release() {
	this.ctx = nil
	this.codec = nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *codec) Name() string {
	return this.codec.Name()
}

func (this *codec) Description() string {
	return this.codec.Description()
}

func (this *codec) Flags() gopi.MediaFlag {
	flags := gopi.MEDIA_FLAG_NONE

	switch {
	case this.ctx != nil:
		switch this.ctx.Type() {
		case ffmpeg.AVMEDIA_TYPE_VIDEO:
			if this.ctx.BitRate() > 0 {
				flags |= gopi.MEDIA_FLAG_VIDEO
			}
		case ffmpeg.AVMEDIA_TYPE_AUDIO:
			flags |= gopi.MEDIA_FLAG_AUDIO
		case ffmpeg.AVMEDIA_TYPE_SUBTITLE:
			flags |= gopi.MEDIA_FLAG_SUBTITLE
		case ffmpeg.AVMEDIA_TYPE_UNKNOWN, ffmpeg.AVMEDIA_TYPE_DATA:
			flags |= gopi.MEDIA_FLAG_DATA
		case ffmpeg.AVMEDIA_TYPE_ATTACHMENT:
			flags |= gopi.MEDIA_FLAG_ATTACHMENT
		}
	case this.codec != nil:
		switch this.codec.Type() {
		case ffmpeg.AVMEDIA_TYPE_VIDEO:
			flags |= gopi.MEDIA_FLAG_VIDEO
		case ffmpeg.AVMEDIA_TYPE_AUDIO:
			flags |= gopi.MEDIA_FLAG_AUDIO
		case ffmpeg.AVMEDIA_TYPE_SUBTITLE:
			flags |= gopi.MEDIA_FLAG_SUBTITLE
		case ffmpeg.AVMEDIA_TYPE_UNKNOWN, ffmpeg.AVMEDIA_TYPE_DATA:
			flags |= gopi.MEDIA_FLAG_DATA
		case ffmpeg.AVMEDIA_TYPE_ATTACHMENT:
			flags |= gopi.MEDIA_FLAG_ATTACHMENT
		}
	}

	// Encode and decode flags
	if this.codec.IsEncoder() {
		flags |= gopi.MEDIA_FLAG_ENCODER
	}
	if this.codec.IsDecoder() {
		flags |= gopi.MEDIA_FLAG_DECODER
	}

	// Return flags
	return flags
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *codec) String() string {
	str := "<ffmpeg.codec"
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	if this.ctx != nil {
		str += " ctx_type=" + fmt.Sprint(this.ctx.Type())
	}
	if this.codec != nil {
		str += " codec_type=" + fmt.Sprint(this.codec.Type())
	}
	return str + ">"
}
