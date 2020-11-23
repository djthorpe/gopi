// +build ffmpeg

package ffmpeg_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/media/ffmpeg"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type MediaApp struct {
	gopi.Unit
	*ffmpeg.Manager
}

func Test_Discovery_001(t *testing.T) {
	tool.Test(t, nil, new(MediaApp), func(app *MediaApp) {
		if app.Manager == nil {
			t.Error("manager is nil")
		} else {
			t.Log(app.Manager)
		}
	})
}
