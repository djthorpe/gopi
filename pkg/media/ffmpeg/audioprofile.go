// +build ffmpeg

package ffmpeg

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type AudioProfile struct {
	fmt      ffmpeg.AVSampleFormat
	rate     uint
	channels uint
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
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *AudioProfile) Flags() gopi.MediaFlag {
	return gopi.MEDIA_FLAG_AUDIO
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
