// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type AudioProfile struct {
	sync.RWMutex

	fmt      ffmpeg.AVSampleFormat
	rate     uint
	channels uint
	layout   ffmpeg.AVChannelLayout
	ctx      *ffmpeg.SwrContext
	frame    *ffmpeg.AVFrame
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewAudioProfile(fmt gopi.AudioFormat, rate uint, layout gopi.AudioChannelLayout) *AudioProfile {
	this := new(AudioProfile)
	if fmt := toSampleFormat(fmt); fmt == ffmpeg.AV_SAMPLE_FMT_NONE {
		return nil
	} else if rate == 0 {
		return nil
	} else if layout.Channels == 0 {
		return nil
	} else {
		this.fmt = fmt
		this.rate = rate
		this.channels = layout.Channels
	}

	// Return success
	return this
}

func (this *AudioProfile) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Free resources
	if this.ctx != nil {
		this.ctx.Free()
	}
	if this.frame != nil {
		this.frame.Free()
	}

	// Release resources
	this.ctx = nil
	this.frame = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *AudioProfile) Flags() gopi.MediaFlag {
	return gopi.MEDIA_FLAG_AUDIO
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Resample returns resampled audio frames which adhere to profile. The returned value
// is either a frame or nil if the resampling operation is still in process. After all
// frames have been called into this method a final call with nil is required
// to flush the last resampled frame.
func (this *AudioProfile) Resample(src *ffmpeg.AVFrame) (*ffmpeg.AVFrame, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If Resample is called with nil then Flush frame
	if src == nil && this.ctx != nil && this.frame != nil {
		if err := this.ctx.FlushFrame(this.frame); err != nil {
			return nil, err
		} else if this.frame.NumSamples() > 0 {
			return this.frame, nil
		} else {
			return nil, nil
		}
	}

	// Return error if called with nil
	if src == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("Resample")
	}

	// Check incoming frame parameters and create context
	if src_fmt := src.SampleFormat(); src_fmt == ffmpeg.AV_SAMPLE_FMT_NONE {
		return nil, gopi.ErrBadParameter.WithPrefix("Resample")
	} else if this.ctx == nil {
		// Initialize frame and context
		if dest := ffmpeg.NewAudioFrame(this.fmt, int(this.rate), this.layout); dest == nil {
			return nil, gopi.ErrInternalAppError.WithPrefix("Resample")
		} else if ctx := ffmpeg.NewSwrContextEx(src_fmt, this.fmt, src.SampleRate(), int(this.rate), src.ChannelLayout(), this.layout); ctx == nil {
			dest.Free()
			return nil, gopi.ErrInternalAppError.WithPrefix("Resample")
		} else if err := this.ctx.ConfigFrame(dest, src); err != nil {
			dest.Free()
			ctx.Free()
			return nil, err
		} else if ctx.IsInitialized() == false {
			dest.Free()
			ctx.Free()
			return nil, gopi.ErrUnexpectedResponse.WithPrefix("Resample")
		} else {
			this.ctx = ctx
			this.frame = dest
		}
	}

	// Resample frame and return the frame if there is data else return nil
	if err := this.ctx.ConvertFrame(this.frame, src); err != nil {
		return nil, err
	} else if this.frame.NumSamples() > 0 {
		return this.frame, nil
	} else {
		return nil, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *AudioProfile) String() string {
	str := "<ffmpeg.profile"
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += fmt.Sprint(" flags=", flags)
	}
	if this.fmt != ffmpeg.AV_SAMPLE_FMT_NONE {
		str += fmt.Sprint(" fmt=", this.fmt)
	}
	if this.rate != 0 {
		str += fmt.Sprint(" sample_rate=", this.rate)
	}
	if this.channels != 0 {
		str += fmt.Sprint(" channels=", this.channels)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func toSampleFormat(fmt gopi.AudioFormat) ffmpeg.AVSampleFormat {
	switch fmt {
	case gopi.AUDIO_FMT_U8:
		return ffmpeg.AV_SAMPLE_FMT_U8
	case gopi.AUDIO_FMT_U8P:
		return ffmpeg.AV_SAMPLE_FMT_U8P
	case gopi.AUDIO_FMT_S16:
		return ffmpeg.AV_SAMPLE_FMT_S16
	case gopi.AUDIO_FMT_S16P:
		return ffmpeg.AV_SAMPLE_FMT_S16P
	case gopi.AUDIO_FMT_S32:
		return ffmpeg.AV_SAMPLE_FMT_S32
	case gopi.AUDIO_FMT_S32P:
		return ffmpeg.AV_SAMPLE_FMT_S32P
	case gopi.AUDIO_FMT_F32:
		return ffmpeg.AV_SAMPLE_FMT_FLT
	case gopi.AUDIO_FMT_F32P:
		return ffmpeg.AV_SAMPLE_FMT_FLTP
	case gopi.AUDIO_FMT_F64:
		return ffmpeg.AV_SAMPLE_FMT_DBL
	case gopi.AUDIO_FMT_F64P:
		return ffmpeg.AV_SAMPLE_FMT_DBLP
	case gopi.AUDIO_FMT_S64:
		return ffmpeg.AV_SAMPLE_FMT_S64
	case gopi.AUDIO_FMT_S64P:
		return ffmpeg.AV_SAMPLE_FMT_S64P
	default:
		return ffmpeg.AV_SAMPLE_FMT_NONE
	}
}
