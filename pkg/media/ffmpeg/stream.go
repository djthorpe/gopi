// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type stream struct {
	sync.RWMutex
	*AudioProfile

	ctx   *ffmpeg.AVStream
	codec *codec
}

////////////////////////////////////////////////////////////////////////////////
// INIT

// NewStream returns a stream object, can optiomally copy codec
// parameters from another stream if source is not set to nil
func NewStream(ctx *ffmpeg.AVStream, source *stream) *stream {
	this := new(stream)

	if ctx == nil {
		return nil
	}
	if source == nil {
		if codec := NewCodecWithParameters(ctx.CodecPar()); codec == nil {
			return nil
		} else {
			this.ctx = ctx
			this.codec = codec
		}
	} else {
		if codec := NewCodecWithParameters(source.ctx.CodecPar()); codec == nil {
			return nil
		} else {
			this.ctx = ctx
			this.codec = codec
		}
	}

	// Return success
	return this
}

func (this *stream) Release() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Release
	var result error
	if this.codec != nil {
		if err := this.codec.Release(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.AudioProfile != nil {
		if err := this.AudioProfile.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Set instance variables to nil
	this.ctx = nil
	this.codec = nil
	this.AudioProfile = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// METHODS - STREAM

func (this *stream) Index() int {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return -1
	} else {
		return this.ctx.Index()
	}
}

func (this *stream) Flags() gopi.MediaFlag {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Return NONE if released
	flags := gopi.MEDIA_FLAG_NONE
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
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return nil
	} else {
		return this.codec
	}
}

func (this *stream) Profile() gopi.MediaProfile {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.AudioProfile != nil {
		return this.AudioProfile
	} else {
		return nil
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
		str += fmt.Sprint(" flags=", flags)
	}
	if codec := this.Codec(); codec != nil {
		str += fmt.Sprint(" codec=", codec)
	}
	if profile := this.Profile(); profile != nil {
		str += fmt.Sprint(" profile=", profile)
	}
	return str + ">"
}
