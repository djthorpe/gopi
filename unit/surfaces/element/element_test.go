// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package element_test

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
	element "github.com/djthorpe/gopi/v2/unit/surfaces/element"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN DISPLAY

var (
	Display rpi.DXDisplayHandle
)

func init() {
	rpi.DXInit()
	if display, err := rpi.DXDisplayOpen(rpi.DXDisplayId(0)); err != nil {
		panic(fmt.Sprint(err))
	} else {
		Display = display
	}
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Element_000(t *testing.T) {
	t.Log("Test_Element_000")
}

func Test_Element_001(t *testing.T) {
	// Non-zero size is required
	if _, err := gopi.New(element.Config{}, NewLogger(t)); errors.Is(err, gopi.ErrBadParameter) == false {
		t.Error("Unexpected error return", err)
	}
	// Update parameter is required
	if element, err := gopi.New(element.Config{Size: gopi.Size{1, 1}}, NewLogger(t)); errors.Is(err, gopi.ErrBadParameter) == false {
		t.Error("Unexpected error return", err)
	} else {
		t.Log(element)
	}
	// Display parameter is required
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 1*time.Second)
	if element, err := gopi.New(element.Config{Size: gopi.Size{1, 1}, Update: update}, NewLogger(t)); errors.Is(err, gopi.ErrBadParameter) == false {
		t.Error("Unexpected error return", err)
	} else {
		t.Log(element)
	}
}

func Test_Element_002(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 1*time.Second)

	if element, err := gopi.New(element.Config{Size: gopi.Size{1, 1}, Update: update, Display: Display}, NewLogger(t)); err != nil {
		t.Error(err)
	} else {
		t.Log(element)
	}
}

func Test_Element_003(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 1*time.Second)

	// Create a 10x10 bitmap
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{10, 10}}, NewLogger(t)); err != nil {
		t.Error("Unexpected error return", err)
	} else if em, err := gopi.New(element.Config{Bitmap: bm.(bitmap.Bitmap), Update: update, Display: Display}, NewLogger(t)); err != nil {
		t.Error(err)
	} else if em.(element.Element).Size().W != 10 || em.(element.Element).Size().H != 10 {
		t.Error("Unexpected element size", em.(element.Element).Size())
	} else {
		t.Log(em)
	}
}

func Test_Element_004(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 1*time.Second)

	// Create a 10x10 bitmap but stretch to 20x10
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{10, 10}}, NewLogger(t)); err != nil {
		t.Error("Unexpected error return", err)
	} else if em, err := gopi.New(element.Config{Bitmap: bm.(bitmap.Bitmap), Size: gopi.Size{20, 0}, Opacity: 1.0, Update: update, Display: Display}, nil); err != nil {
		t.Error(err)
	} else if em.(element.Element).Size().W != 20 || em.(element.Element).Size().H != 10 {
		t.Error("Unexpected element size", em.(element.Element).Size())
	} else {
		t.Log(em)
	}

	// Create a 10x10 bitmap but stretch to 10x20
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{10, 10}}, NewLogger(t)); err != nil {
		t.Error("Unexpected error return", err)
	} else if em, err := gopi.New(element.Config{Bitmap: bm.(bitmap.Bitmap), Size: gopi.Size{0, 20}, Opacity: 1.0, Update: update, Display: Display}, nil); err != nil {
		t.Error(err)
	} else if em.(element.Element).Size().W != 10 || em.(element.Element).Size().H != 20 {
		t.Error("Unexpected element size", em.(element.Element).Size())
	} else {
		t.Log(em)
	}
}

func Test_Element_005(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 1*time.Second)

	// Create a 2x2 bitmap with gray
	bm, err := gopi.New(bitmap.Config{Size: gopi.Size{2, 2}}, NewLogger(t))
	bm.(bitmap.Bitmap).Set(0, 0, color.Gray{0})
	bm.(bitmap.Bitmap).Set(0, 1, color.Gray{128})
	bm.(bitmap.Bitmap).Set(1, 0, color.Gray{128})
	bm.(bitmap.Bitmap).Set(1, 1, color.Gray{255})

	if err != nil {
		t.Error("Unexpected error return", err)
	} else if em, err := gopi.New(element.Config{Bitmap: bm.(bitmap.Bitmap), Size: gopi.Size{100, 100}, Opacity: 1.0, Update: update, Display: Display}, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(em)
	}
}

func Test_Element_006(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 5*time.Second)

	size := gopi.Size{100, 100}
	if bm, err := gopi.New(bitmap.Config{Size: size, Mode: gopi.SURFACE_FLAG_RGB565}, NewLogger(t)); err != nil {
		t.Fatal(err)
	} else if em, err := gopi.New(element.Config{Bitmap: bm.(bitmap.Bitmap), Opacity: 1.0, Update: update, Display: Display}, nil); err != nil {
		t.Error(err)
	} else {
		blue := color.RGBA{0, 0, 255, 255}
		draw.Draw(bm.(bitmap.Bitmap), bm.(bitmap.Bitmap).Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)
		t.Log(em)
	}
}

func Test_Element_007(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 5*time.Second)

	im, err := LoadImage()
	if err != nil {
		t.Fatal(err)
	}

	bmsize := gopi.Size{float32(im.Bounds().Dx()), float32(im.Bounds().Dy())}
	bm, err := gopi.New(bitmap.Config{Size: bmsize, Mode: gopi.SURFACE_FLAG_RGBA32}, NewLogger(t))
	if err != nil {
		t.Fatal(err)
	}

	draw.Draw(bm.(bitmap.Bitmap), bm.(bitmap.Bitmap).Bounds(), im, image.Point{}, draw.Src)

	emorigin := gopi.Point{100, 100}
	if em, err := gopi.New(element.Config{
		Bitmap:  bm.(bitmap.Bitmap),
		Flags:   gopi.SURFACE_FLAG_ALPHA_FROM_SOURCE,
		Origin:  emorigin,
		Opacity: 1.0,
		Update:  update,
		Display: Display,
	}, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(em)
	}
}

func Test_Element_008(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 5*time.Second)

	bmsize := gopi.Size{200, 200}
	bm, err := gopi.New(bitmap.Config{Size: bmsize, Mode: gopi.SURFACE_FLAG_RGBA32}, NewLogger(t))
	if err != nil {
		t.Fatal(err)
	}
	bm.(bitmap.Bitmap).PaintCircle(color.White, bm.(bitmap.Bitmap).Centre(), 50)

	if em, err := gopi.New(element.Config{
		Bitmap:  bm.(bitmap.Bitmap),
		Flags:   gopi.SURFACE_FLAG_ALPHA_FROM_SOURCE,
		Opacity: 1.0,
		Update:  update,
		Display: Display,
	}, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(em)
	}
}

func Test_Element_009(t *testing.T) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 5*time.Second)

	bmsize := gopi.Size{200, 200}
	bm, err := gopi.New(bitmap.Config{Size: bmsize, Mode: gopi.SURFACE_FLAG_RGBA32}, NewLogger(t))
	if err != nil {
		t.Fatal(err)
	}

	red := color.RGBA{0xff, 0, 0, 0xff}
	bm.(bitmap.Bitmap).ClearToColor(red)
	bm.(bitmap.Bitmap).PaintLine(color.White, bm.(bitmap.Bitmap).NorthWest(), bm.(bitmap.Bitmap).SouthEast())
	bm.(bitmap.Bitmap).PaintLine(color.White, bm.(bitmap.Bitmap).SouthWest(), bm.(bitmap.Bitmap).NorthEast())

	emorigin := gopi.Point{100, 100}
	if em, err := gopi.New(element.Config{
		Bitmap:  bm.(bitmap.Bitmap),
		Origin:  emorigin,
		Flags:   gopi.SURFACE_FLAG_ALPHA_FROM_SOURCE,
		Opacity: 1.0,
		Update:  update,
		Display: Display,
	}, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(em)
	}
}

func Test_Element_010(t *testing.T) {
	if bm1, err := NewBitmapRGBA32(t, gopi.Size{200, 200}, color.RGBA{0xff, 0, 0, 0xff}); err != nil {
		t.Fatal(err)
	} else {
		bm1.PaintLine(color.White, bm1.NorthWest(), bm1.SouthEast())
		bm1.PaintLine(color.White, bm1.SouthWest(), bm1.NorthEast())
		if em, err := AddElementWithBitmap(t, bm1); err != nil {
			t.Fatal(err)
		} else if err := ResizeElement(t, em, gopi.Size{400, 400}); err != nil {
			t.Fatal(err)
		} else if err := ResizeElement(t, em, gopi.Size{40, 40}); err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Element_011(t *testing.T) {
	if bm1, err := NewBitmapRGBA32(t, gopi.Size{200, 200}, color.RGBA{0xff, 0, 0, 0xff}); err != nil {
		t.Fatal(err)
	} else if bm2, err := NewBitmapRGBA32(t, gopi.Size{400, 400}, color.RGBA{0x00, 0, 0xFF, 0xff}); err != nil {
		t.Fatal(err)
	} else {
		bm1.PaintLine(color.White, bm1.NorthWest(), bm1.SouthEast())
		bm1.PaintLine(color.White, bm1.SouthWest(), bm1.NorthEast())
		bm2.PaintLine(color.White, bm1.NorthWest(), bm1.SouthEast())
		bm2.PaintLine(color.White, bm1.SouthWest(), bm1.NorthEast())

		// Retain the bitmaps as we'll release them as they are replaced
		bm1.Retain()
		defer bm1.Release()
		bm2.Retain()
		defer bm2.Release()

		// Switch between bitmaps
		if em, err := AddElementWithBitmap(t, bm1); err != nil {
			t.Fatal(err)
		} else {
			for i := 0; i < 100; i++ {
				if err := ReplaceBitmap(t, em, bm2); err != nil {
					t.Fatal(err)
				}
				if err := ReplaceBitmap(t, em, bm1); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// MAKE NEW BITMAP WITH COLOR

func NewBitmapRGBA32(t *testing.T, size gopi.Size, c color.Color) (bitmap.Bitmap, error) {
	if bm, err := gopi.New(bitmap.Config{Size: size, Mode: gopi.SURFACE_FLAG_RGBA32}, NewLogger(t)); err != nil {
		return nil, err
	} else {
		bm.(bitmap.Bitmap).ClearToColor(c)
		return bm.(bitmap.Bitmap), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// ADD ELEMENT WITH BITMAP

func AddElementWithBitmap(t *testing.T, bm bitmap.Bitmap) (element.Element, error) {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 2*time.Second)

	emorigin := gopi.Point{100, 100}
	if em, err := gopi.New(element.Config{
		Bitmap:  bm.(bitmap.Bitmap),
		Origin:  emorigin,
		Flags:   gopi.SURFACE_FLAG_ALPHA_FROM_SOURCE,
		Opacity: 1.0,
		Update:  update,
		Display: Display,
	}, nil); err != nil {
		return nil, err
	} else {
		return em.(element.Element), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// RESIZE ELEMENT

func ResizeElement(t *testing.T, em element.Element, size gopi.Size) error {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 2*time.Second)

	return em.SetSize(update, size)
}

func ReplaceBitmap(t *testing.T, em element.Element, bm bitmap.Bitmap) error {
	update := UpdateStart(t)
	defer UpdateEnd(t, update, 100*time.Millisecond)

	return em.SetBitmap(update, bm)
}

////////////////////////////////////////////////////////////////////////////////
// RETURN AN IMAGE

func LoadImage() (image.Image, error) {
	// Load in an image
	if wd, err := os.Getwd(); err != nil {
		return nil, err
	} else {
		imagePath := filepath.Join(wd, "..", "..", "..", "etc", "images", "gopi-400x194-white.png")
		if fh, err := os.Open(imagePath); err != nil {
			return nil, err
		} else {
			defer fh.Close()
			if im, err := png.Decode(fh); err != nil {
				return nil, err
			} else {
				return im, nil
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// MAKE UPDATE

func UpdateStart(t *testing.T) rpi.DXUpdate {
	if handle, err := rpi.DXUpdateStart(0); err != nil {
		t.Fatal(err)
		return 0
	} else {
		return handle
	}
}

func UpdateEnd(t *testing.T, handle rpi.DXUpdate, sleep time.Duration) {
	if err := rpi.DXUpdateSubmitSync(handle); err != nil {
		t.Error(err)
	}
	time.Sleep(sleep)
}

////////////////////////////////////////////////////////////////////////////////
// LOGGER

type logger struct{ testing.T }

func NewLogger(t *testing.T) gopi.Logger {
	return &logger{*t}
}

func (this *logger) Clone(string) gopi.Logger {
	return this
}
func (this *logger) Name() string {
	return this.T.Name()
}
func (this *logger) Error(err error) error {
	this.T.Error(err)
	return err
}
func (this *logger) Warn(args ...interface{}) {
	this.T.Log(args...)
}
func (this *logger) Info(args ...interface{}) {
	this.T.Log(args...)
}
func (this *logger) Debug(args ...interface{}) {
	this.T.Log(args...)
}
func (this *logger) IsDebug() bool {
	return true
}
func (this *logger) Close() error {
	return nil
}
func (this *logger) String() string {
	return this.Name()
}
