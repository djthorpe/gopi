// +build ffmpeg

package ffmpeg_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/media/ffmpeg"
)

func Test_AudioProfile_001(t *testing.T) {
	profile := ffmpeg.NewAudioProfile(gopi.AUDIO_FMT_S16, 41000, gopi.AudioLayoutMono)
	if profile == nil {
		t.Error("Unexpected nil returned")
	} else {
		t.Log(profile)
	}
}
