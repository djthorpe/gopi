// +build ffmpeg

package ffmpeg_test

import (
	"testing"

	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TEST

func Test_avdevice_000(t *testing.T) {
	ffmpeg.AVDeviceInit()
	t.Log("Test_avdevice_000")
}

func Test_avdevice_001(t *testing.T) {
	ffmpeg.AVDeviceInit()

	audioinput := ffmpeg.AllAudioInputDevices()
	t.Log(audioinput)
}

func Test_avdevice_002(t *testing.T) {
	ffmpeg.AVDeviceInit()

	inputs := ffmpeg.AllVideoInputDevices()
	t.Log(inputs)
}

func Test_avdevice_003(t *testing.T) {
	ffmpeg.AVDeviceInit()

	outputs := ffmpeg.AllAudioOutputDevices()
	t.Log(outputs)
}

func Test_avdevice_004(t *testing.T) {
	ffmpeg.AVDeviceInit()

	outputs := ffmpeg.AllVideoOutputDevices()
	t.Log(outputs)
}
