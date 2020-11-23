// +build ffmpeg

package ffmpeg

import (
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libavdevice
#include <libavdevice/avdevice.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	once_avdevice_init sync.Once
)

////////////////////////////////////////////////////////////////////////////////
// INIT

// Register all devices
func AVDeviceInit() {
	once_avdevice_init.Do(func() {
		C.avdevice_register_all()
	})
}

////////////////////////////////////////////////////////////////////////////////
// ENMERATE DEVICES

func AllAudioInputDevices() []*AVInputFormat {
	formats := make([]*AVInputFormat, 0)
	ptr := (*C.AVInputFormat)(nil)
	for {
		if ptr = C.av_input_audio_device_next(ptr); ptr == nil {
			break
		} else {
			formats = append(formats, (*AVInputFormat)(ptr))
		}
	}
	return formats
}

func AllVideoInputDevices() []*AVInputFormat {
	formats := make([]*AVInputFormat, 0)
	ptr := (*C.AVInputFormat)(nil)
	for {
		if ptr = C.av_input_video_device_next(ptr); ptr == nil {
			break
		} else {
			formats = append(formats, (*AVInputFormat)(ptr))
		}
	}
	return formats
}

func AllVideoOutputDevices() []*AVOutputFormat {
	formats := make([]*AVOutputFormat, 0)
	ptr := (*C.AVOutputFormat)(nil)
	for {
		if ptr = C.av_output_video_device_next(ptr); ptr == nil {
			break
		} else {
			formats = append(formats, (*AVOutputFormat)(ptr))
		}
	}
	return formats
}

func AllAudioOutputDevices() []*AVOutputFormat {
	formats := make([]*AVOutputFormat, 0)
	ptr := (*C.AVOutputFormat)(nil)
	for {
		if ptr = C.av_output_audio_device_next(ptr); ptr == nil {
			break
		} else {
			formats = append(formats, (*AVOutputFormat)(ptr))
		}
	}
	return formats
}
