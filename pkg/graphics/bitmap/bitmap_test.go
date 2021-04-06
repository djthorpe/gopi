package bitmap_test

import (
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	bitmap "github.com/djthorpe/gopi/v3/pkg/graphics/bitmap"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	// Dependencies
	_ "github.com/djthorpe/gopi/v3/pkg/graphics/bitmap/rgba32"
	_ "github.com/djthorpe/gopi/v3/pkg/graphics/bitmap/rgba32dx"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
)

type App struct {
	gopi.Unit
	*bitmap.Bitmaps
}

const (
	PNG_FILEPATH = "../../../etc/images/gopi-800x388.png"
)

func Test_Bitmap_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		app.Require(app.Bitmaps)

		if bitmap, err := app.NewBitmap(gopi.SURFACE_FMT_RGBA32, 10, 10); err != nil {
			t.Error(err)
		} else {
			t.Log(bitmap)
		}
	})
}

func Test_Bitmap_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		app.Require(app.Bitmaps)

		bitmap, err := app.NewBitmap(gopi.SURFACE_FMT_RGBA32, 100, 100)
		if err != nil {
			t.Fatal(err)
		}
		bitmap.ClearToColor(color.RGBA{0xFF, 0x00, 0x00, 0xFF})
		if writer, err := ioutil.TempFile("", "png"); err != nil {
			t.Error(err)
		} else if err := png.Encode(writer, bitmap); err != nil {
			t.Error(err)
		} else {
			writer.Close()
			t.Log(writer.Name())
		}
	})
}

func Test_Bitmap_003(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		app.Require(app.Bitmaps)

		reader, err := os.Open(PNG_FILEPATH)
		if err != nil {
			t.Fatal(err)
		}
		defer reader.Close()
		if bitmap, _, err := image.Decode(reader); err != nil {
			t.Error(err)
		} else if dest, err := app.NewBitmap(gopi.SURFACE_FMT_RGBA32, uint32(bitmap.Bounds().Dx()), uint32(bitmap.Bounds().Dy())); err != nil {
			t.Error(err)
		} else {
			bounds := bitmap.Bounds()
			for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
				for x := bounds.Min.X; x <= bounds.Max.X; x++ {
					c := bitmap.At(x, y)
					dest.SetAt(c, x, y)
				}
			}

			if writer, err := ioutil.TempFile("", "png"); err != nil {
				t.Error(err)
			} else if err := png.Encode(writer, dest); err != nil {
				t.Error(err)
			} else {
				writer.Close()
				t.Log(writer.Name())
			}
		}
	})
}
