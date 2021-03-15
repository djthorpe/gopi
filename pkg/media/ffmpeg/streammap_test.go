// +build ffmpeg

package ffmpeg_test

import (
	"testing"

	ffmpeg "github.com/djthorpe/gopi/v3/pkg/media/ffmpeg"
)

func Test_StreamMap_001(t *testing.T) {
	if streammap := ffmpeg.NewStreamMap(); streammap == nil {
		t.Error("Unexpected nil return")
	}
}
