// +build ffmpeg

package ffmpeg_test

import (
	"io"
	"testing"

	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

func Test_avcodeccontext_001(t *testing.T) {
	if ctx := ffmpeg.NewAVCodecContext(nil); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else {
		t.Log(ctx)
		ctx.Free()
	}
}

func Test_avcodeccontext_002(t *testing.T) {
	if codec := ffmpeg.FindEncoderByName("rawvideo"); codec == nil {
		t.Fatal("FindEncoderByName failed")
	} else if ctx := ffmpeg.NewAVCodecContext(codec); ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else {
		t.Log(ctx)
		ctx.Free()
	}
}

func Test_avcodeccontext_009(t *testing.T) {
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

func Test_avcodeccontext_005(t *testing.T) {
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
					t.Log(packet)
					packet.Release()
				}
			}
		}
	}
}

/*

func Test_avformat_012(t *testing.T) {
	ctx := ffmpeg.NewAVFormatContext()
	if ctx == nil {
		t.Fatal("NewAVFormatContext failed")
	} else if err := ctx.OpenInput(SAMPLE_MP4, nil); err != nil {
		t.Fatal(err)
	}
	defer ctx.CloseInput()
	packet := ffmpeg.NewAVPacket()
	if packet == nil {
		t.Fatal("Unexpected packet == nil")
	}
	defer packet.Free()
	frame := ffmpeg.NewAVFrame()
	if frame == nil {
		t.Fatal("Unexpected frame == nil")
	}
	defer frame.Free()
	for {
		err := ctx.ReadPacket(packet)
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
			break
		}
		for {

			if err := ctx.DecodeFrame(frame); errors.Is(err, syscall.EAGAIN) {
				continue
			} else if errors.Is(err, syscall.EINVAL) {
				break
			} else if err != nil {
				t.Error(err)
			} else {
				t.Log(frame)
			}
		}
		packet.Release()
	}
}

*/
