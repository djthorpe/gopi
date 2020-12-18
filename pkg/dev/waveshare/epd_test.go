// +build linux

package waveshare_test

import (
	"context"
	"image"
	"os"
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/dev/waveshare"
	"github.com/djthorpe/gopi/v3/pkg/tool"

	_ "image/jpeg"

	_ "github.com/djthorpe/gopi/v3/pkg/hw/gpio"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/spi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*waveshare.EPD
}

const (
	SAMPLE_IMAGE = "../../../etc/images/gopi-880x528.jpg"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_EPD_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.EPD == nil {
			t.Error("nil EPD unit")
		} else {
			t.Log(app.EPD)
		}
	})
}

func Test_EPD_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if err := app.EPD.Clear(context.Background()); err != nil {
			t.Error(err)
		}
	})
}

func Test_EPD_003(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		fh, err := os.Open(SAMPLE_IMAGE)
		if err != nil {
			t.Fatal(err)
		}
		defer fh.Close()
		if img, _, err := image.Decode(fh); err != nil {
			t.Fatal(err)
		} else if err := app.EPD.Draw(context.Background(), img); err != nil {
			t.Error(err)
		}
	})
}
