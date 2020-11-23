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
	ctx *ffmpeg.AVStream
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewStream(ctx *ffmpeg.AVStream) *stream {
	if ctx == nil {
		return nil
	} else {
		return &stream{ctx}
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *stream) Index() int {
	return this.ctx.Index()
}

func (this *stream) Flags() gopi.MediaFlag {
	flags := gopi.MEDIA_FLAG_NONE

	// Codec flags
	switch this.ctx.CodecPar().Type() {
	case ffmpeg.AVMEDIA_TYPE_VIDEO:
		if this.ctx.CodecPar().BitRate() > 0 {
			flags |= gopi.MEDIA_FLAG_VIDEO
		}
	case ffmpeg.AVMEDIA_TYPE_AUDIO:
		flags |= gopi.MEDIA_FLAG_AUDIO
	case ffmpeg.AVMEDIA_TYPE_UNKNOWN, ffmpeg.AVMEDIA_TYPE_DATA:
		flags |= gopi.MEDIA_FLAG_DATA
	case ffmpeg.AVMEDIA_TYPE_ATTACHMENT:
		flags |= gopi.MEDIA_FLAG_ATTACHMENT
	}

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

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stream) String() string {
	str := "<stream"
	str += " index=" + fmt.Sprint(this.Index())
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	return str + ">"
}
