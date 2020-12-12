// +build ffmpeg

package ffmpeg_test

import (
	"testing"

	ffmpeg "github.com/djthorpe/gopi/v3/pkg/media/ffmpeg"
	bindings "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

func Test_StreamMap_001(t *testing.T) {
	if streammap := ffmpeg.NewStreamMap(); streammap == nil {
		t.Error("Unexpected nil return")
	}
}

func Test_StreamMap_002(t *testing.T) {
	streammap := ffmpeg.NewStreamMap()
	if err := streammap.Set(nil, nil); err == nil {
		t.Fatal("Expected error return")
	}
	ctx := bindings.NewAVFormatContext()
	if ctx == nil {
		t.Fatal("Unexpected nil return")
	}
	if err := ctx.OpenInput(SAMPLE_FILE, nil); err != nil {
		t.Fatal(err)
	}
	defer ctx.CloseInput()
	if len(ctx.Streams()) == 0 {
		t.Fatal("Unexpected zero streams")
	}
	for _, stream := range ctx.Streams() {
		if err := streammap.Add(stream, nil); err != nil {
			t.Error(err, stream)
		}
		if out := streammap.Get(stream); out != nil {
			t.Error("Unexpected nil return", stream)
		}
		streammap.Set(stream, stream)
		if out := streammap.Get(stream); out != stream {
			t.Error("Unexpected nil return", stream)
		}
	}
}
