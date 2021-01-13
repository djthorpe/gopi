// +build mmal

package mmal_test

import (
	"context"
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	mmal "github.com/djthorpe/gopi/v3/pkg/media/mmal"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

const (
	SAMPLE_FILE = "../../../etc/images/gopi-880x528.jpg"
)

type MMALApp struct {
	gopi.Unit
	*mmal.Manager
}

func (this *MMALApp) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func Test_MMALManager_001(t *testing.T) {
	tool.Test(t, nil, new(MMALApp), func(app *MMALApp) {
		if app.Manager == nil {
			t.Error("manager is nil")
		} else {
			t.Log(app.Manager)
		}
	})
}

func Test_MMALManager_002(t *testing.T) {
	tool.Test(t, nil, new(MMALApp), func(app *MMALApp) {
		if c, err := app.Manager.VideoDecoder(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		if c, err := app.Manager.VideoEncoder(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		if c, err := app.Manager.VideoRenderer(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		if c, err := app.Manager.ImageDecoder(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		if c, err := app.Manager.ImageEncoder(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		// Camera may not be created if it isn't enabled
		/*if c, err := app.Manager.Camera(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}*/
		if c, err := app.Manager.CameraInfo(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		if c, err := app.Manager.VideoSplitter(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		if c, err := app.Manager.AudioRenderer(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
		if c, err := app.Manager.Clock(); err != nil {
			t.Error(err)
		} else {
			t.Log(c)
		}
	})
}

/*
func Test_MMALManager_003(t *testing.T) {
	tool.Test(t, nil, new(MMALApp), func(app *MMALApp) {
		input, err := os.Open(SAMPLE_FILE)
		if err != nil {
			t.Fatal(err)
		}
		defer input.Close()
		output := new(bytes.Buffer)
		decoder, err := app.Manager.ImageDecoder()
		if err != nil {
			t.Fatal(err)
		}

		// TODO! Set format on input port
		if err := decoder.SetInputFormatJPEG(); err != nil {
			t.Fatal(err)
		}

		if _, err := app.Manager.CreateReaderForComponent(input, decoder, 0); err != nil {
			t.Error(err)
		} else if _, err := app.Manager.CreateWriterForComponent(output, decoder, 0); err != nil {
			t.Error(err)
		} else {
			// Execute the graph for up to five seconds
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := app.Manager.Exec(ctx); err != nil {
				t.Error(err)
			}
		}
	})
}
*/
