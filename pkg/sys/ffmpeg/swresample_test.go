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

func Test_swresample_003(t *testing.T) {
	ctx := ffmpeg.NewSwrContext()
	if ctx == nil {
		t.Error("Unexpected nil return from NewSwrContext")
	}
	defer ctx.Free()

	in := ffmpeg.NewAudioFrame(ffmpeg.AV_SAMPLE_FMT_U8, 44100, ffmpeg.AV_CH_LAYOUT_MONO)
	if in == nil {
		t.Fatal("Unexpected nil return for NewAudioFrame")
	} else {
		t.Log("in=", in)
	}
	defer in.Free()

	out := ffmpeg.NewAudioFrame(ffmpeg.AV_SAMPLE_FMT_U8, 11025, ffmpeg.AV_CH_LAYOUT_STEREO)
	if out == nil {
		t.Fatal("Unexpected nil return for NewAudioFrame")
	} else {
		t.Log("out=", out)
	}
	defer out.Free()

	if err := ctx.ConfigFrame(out, in); err != nil {
		t.Error(err)
	} else if ctx.IsInitialized() == false {
		t.Error("Expected IsInitialized=true")
	} else {
		t.Log(ctx)
	}

	// Set number of samples to 10
	if err := in.GetAudioBuffer(10); err != nil {
		t.Error(err)
	}

	// Write out 10 lots of 10 zero samples
	for i := 0; i < 10; i++ {
		if err := ctx.ConvertFrame(out, in); err != nil {
			t.Error(err)
		} else {
			t.Log("out=", out.Buffer(0), " samples=", out.NumSamples())
		}
	}
	if err := ctx.FlushFrame(out); err != nil {
		t.Error(err)
	} else {
		t.Log("out=", out.Buffer(0), " samples=", out.NumSamples())
	}
}
