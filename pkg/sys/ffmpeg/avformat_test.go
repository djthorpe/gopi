// +build ffmpeg

package ffmpeg_test

import (
	"io"
	"testing"

	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

const (
	SAMPLE_MP4 = "../../../etc/sample.mp4"
)

////////////////////////////////////////////////////////////////////////////////
// TEST ENUMS

func Test_avformat_000(t *testing.T) {
	t.Log("Test_avformat_000")
}

func Test_avformat_002(t *testing.T) {
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else {
		ctx.Free()
	}
}

func Test_avformat_003(t *testing.T) {
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else if err := ctx.OpenInput(SAMPLE_MP4, nil); err != nil {
		t.Error(err)
	} else {
		ctx.CloseInput()
	}
}

func Test_avformat_004(t *testing.T) {
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else if err := ctx.OpenInput(SAMPLE_MP4, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(ctx.Metadata())
		ctx.CloseInput()
	}
}

func Test_avformat_005(t *testing.T) {
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else if err := ctx.OpenInput(SAMPLE_MP4, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(ctx.Filename())
		ctx.CloseInput()
	}
}

func Test_avformat_006(t *testing.T) {
	ffmpeg.AVFormatInit()
	ffmpeg.AVFormatInit()
	ffmpeg.AVFormatDeinit()
	ffmpeg.AVFormatDeinit()
}

func Test_avformat_007(t *testing.T) {
	if iformats := ffmpeg.EnumerateInputFormats(); len(iformats) == 0 {
		t.Error("EnumerateInputFormats expected a return value")
	} else {
		for _, iformat := range iformats {
			t.Log(iformat)
		}
	}
}

func Test_avformat_008(t *testing.T) {
	if oformats := ffmpeg.EnumerateOutputFormats(); len(oformats) == 0 {
		t.Error("EnumerateOutputFormats expected a return value")
	} else {
		for _, oformat := range oformats {
			t.Log(oformat)
		}
	}
}

func Test_avformat_009(t *testing.T) {
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else if err := ctx.OpenInput(SAMPLE_MP4, nil); err != nil {
		t.Error(err)
	} else if streams := ctx.Streams(); len(streams) == 0 {
		t.Error("No streams found")
		ctx.CloseInput()
	} else {
		for _, stream := range streams {
			t.Log(stream)
		}
	}
}

func Test_avformat_010(t *testing.T) {
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else if err := ctx.OpenInput(SAMPLE_MP4, nil); err != nil {
		t.Error(err)
	} else {
		defer ctx.CloseInput()
		if packet := ffmpeg.NewAVPacket(); packet == nil {
			t.Error("Unexpected packet == nil")
		} else {
			packet.Free()
		}
	}
}

func Test_avformat_011(t *testing.T) {
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else if err := ctx.OpenInput(SAMPLE_MP4, nil); err != nil {
		t.Error(err)
	} else {
		defer ctx.CloseInput()
		if packet := ffmpeg.NewAVPacket(); packet == nil {
			t.Error("Unexpected packet == nil")
		} else {
			defer packet.Free()
			for {
				if err := ctx.ReadPacket(packet); err == io.EOF {
					break
				} else if err != nil {
					t.Error(err)
					break
				} else {
					packet.Release()
					t.Log(packet)
				}
			}
		}
	}
}
