// +build ffmpeg

package ffmpeg_test

import (
	"testing"

	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TEST ENUMS

func Test_swresample_000(t *testing.T) {
	t.Log("Test_swresample_000")
}

func Test_swresample_001(t *testing.T) {
	if ctx := ffmpeg.NewSwrContext(); ctx == nil {
		t.Error("Unexpected nil return from NewSwrContext")
	} else {
		defer ctx.Free()
		t.Log(ctx)
	}
}

func Test_swresample_002(t *testing.T) {
	if ctx := ffmpeg.NewSwrContextEx(
		ffmpeg.AV_SAMPLE_FMT_U8,
		ffmpeg.AV_SAMPLE_FMT_U8,
		44100,
		44100,
		ffmpeg.AV_CH_LAYOUT_MONO,
		ffmpeg.AV_CH_LAYOUT_MONO,
	); ctx == nil {
		t.Error("Unexpected nil return from NewSwrContext")
	} else {
		defer ctx.Free()
		if err := ctx.Init(); err != nil {
			t.Error(err)
		} else if ctx.IsInitialized() == false {
			t.Error("Expected IsInitialized=true")
		} else {
			t.Log(ctx)
		}
	}
}
