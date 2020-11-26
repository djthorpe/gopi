// +build ffmpeg

package ffmpeg_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/media/ffmpeg"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type MediaApp struct {
	gopi.Unit
	*ffmpeg.Manager
}

const (
	SAMPLE_FILE = "../../../etc/sample.mp4"
)

func Test_MediaManager_001(t *testing.T) {
	tool.Test(t, nil, new(MediaApp), func(app *MediaApp) {
		if app.Manager == nil {
			t.Error("manager is nil")
		} else {
			t.Log(app.Manager)
		}
	})
}

func Test_MediaManager_002(t *testing.T) {
	tool.Test(t, nil, new(MediaApp), func(app *MediaApp) {
		if file, err := app.Manager.OpenFile(SAMPLE_FILE); err != nil {
			t.Error(err)
		} else {
			t.Log(file)
		}
	})
}

func Test_MediaManager_003(t *testing.T) {
	tool.Test(t, nil, new(MediaApp), func(app *MediaApp) {
		if file, err := app.Manager.OpenFile(SAMPLE_FILE); err != nil {
			t.Error(err)
		} else if err := app.Manager.Close(file); err != nil {
			t.Error(err)
		}
	})
}

func Test_MediaManager_004(t *testing.T) {
	tool.Test(t, nil, new(MediaApp), func(app *MediaApp) {
		file, err := app.Manager.OpenFile(SAMPLE_FILE)
		if err != nil {
			t.Error(err)
		}
		defer app.Manager.Close(file)
		if err := file.DecodeIterator(nil, func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
			t.Log(ctx, packet)
			return nil
		}); err != nil {
			t.Error(err)
		}
	})
}

func Test_MediaManager_005(t *testing.T) {
	tool.Test(t, nil, new(MediaApp), func(app *MediaApp) {
		file, err := app.Manager.OpenFile(SAMPLE_FILE)
		if err != nil {
			t.Error(err)
		}
		defer app.Manager.Close(file)
		if err := file.DecodeIterator(nil, func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
			return file.DecodeFrameIterator(ctx, packet, func(frame gopi.MediaFrame) error {
				t.Log("=>", frame)
				return nil
			})
			return nil
		}); err != nil {
			t.Error(err)
		}
	})
}
