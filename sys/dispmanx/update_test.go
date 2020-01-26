// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package dispmanx_test

import (
	"image/color"
	"testing"
	"time"

	// Frameworks
	bitmap "github.com/djthorpe/gopi/v2/sys/dispmanx"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

const (
	DISPLAY_ID = 0
)

func init() {
	rpi.DXInit()
}

func Test_Update_000(t *testing.T) {
	t.Log("Test_Update_000")
}

func Test_Update_001(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	}
}

func Test_Update_002(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if update, err := bitmap.NewUpdate(display); err != nil {
		t.Error(err)
	} else if err := update.Close(); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	}
}

func Test_Update_003(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if update, err := bitmap.NewUpdate(display); err != nil {
		t.Error(err)
	} else {
		defer update.Close()
		defer rpi.DXDisplayClose(display)

		if err := update.Do(0, func() error {
			if bm, err := update.NewBitmap(bitmap.IMAGE_TYPE_RGB565, 100, 100); err != nil {
				return err
			} else {
				bm.ClearToColor(color.Gray16{0x8888})
				if _, err := update.AddElement(rpi.DXNewRect(100, 100, 100, 100), bm, 0, 0xFF); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			t.Error(err)
		}
		time.Sleep(2 * time.Second)
	}
}

func Test_Update_004(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if update, err := bitmap.NewUpdate(display); err != nil {
		t.Error(err)
	} else {
		defer update.Close()
		defer rpi.DXDisplayClose(display)

		bm1, err := update.NewBitmap(bitmap.IMAGE_TYPE_RGB565, 100, 100)
		if err != nil {
			t.Error(err)
		}
		bm2, err := update.NewBitmap(bitmap.IMAGE_TYPE_RGBA32, 100, 100)
		if err != nil {
			t.Error(err)
		}

		// Place first element
		if err := update.Do(0, func() error {
			bm1.ClearToColor(color.RGBA{0xFF, 0xFF, 0, 0xFF})
			if _, err := update.AddElement(rpi.DXNewRect(100, 100, 100, 100), bm1, 0, 0xFF); err != nil {
				return err
			}
			return nil
		}); err != nil {
			t.Error(err)
		}

		time.Sleep(2 * time.Second)

		// Place second element and change first
		if err := update.Do(0, func() error {
			bm1.ClearToColor(color.RGBA{0xFF, 0, 0xFF, 0xFF})
			bm2.ClearToColor(color.RGBA{0xFF, 0, 0, 0xFF})
			if _, err := update.AddElement(rpi.DXNewRect(200, 200, 200, 200), bm2, 0, 0xFF); err != nil {
				return err
			}
			return nil
		}); err != nil {
			t.Error(err)
		}

		time.Sleep(2 * time.Second)
	}
}

func Test_Update_005(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if update, err := bitmap.NewUpdate(display); err != nil {
		t.Error(err)
	} else {
		defer update.Close()
		defer rpi.DXDisplayClose(display)

		bm, err := update.NewBitmap(bitmap.IMAGE_TYPE_RGB565, 2, 2)
		if err != nil {
			t.Error(err)
		}
		bm.PaintPixel(color.RGBA{0xFF, 0, 0, 0xFF}, rpi.DXPoint{0, 0})
		bm.PaintPixel(color.RGBA{0, 0, 0xFF, 0xFF}, rpi.DXPoint{1, 0})
		bm.PaintPixel(color.RGBA{0, 0xFF, 0, 0xFF}, rpi.DXPoint{0, 1})
		bm.PaintPixel(color.RGBA{0, 0xFF, 0xFF, 0xFF}, rpi.DXPoint{1, 1})

		// Place element
		if err := update.Do(0, func() error {
			if _, err := update.AddElement(rpi.DXNewRect(100, 100, 200, 200), bm, 0, 0xFF); err != nil {
				return err
			}
			return nil
		}); err != nil {
			t.Error(err)
		}

		time.Sleep(2 * time.Second)

	}
}

func Test_Update_006(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if update, err := bitmap.NewUpdate(display); err != nil {
		t.Error(err)
	} else {
		defer update.Close()
		defer rpi.DXDisplayClose(display)

		bm, err := update.NewBitmap(bitmap.IMAGE_TYPE_RGB565, 200, 200)
		if err != nil {
			t.Error(err)
		}
		bm.ClearToColor(color.RGBA{0x80, 0x80, 0, 0xFF})
		bm.PaintCircle(color.RGBA{0xFF, 0, 0, 0xFF}, rpi.DXPoint{100, 100}, 50)

		// Place element
		if err := update.Do(0, func() error {
			if _, err := update.AddElement(rpi.DXNewRect(100, 100, 200, 200), bm, 0, 0xFF); err != nil {
				return err
			}
			return nil
		}); err != nil {
			t.Error(err)
		}

		time.Sleep(2 * time.Second)

	}
}

func Test_Update_007(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if update, err := bitmap.NewUpdate(display); err != nil {
		t.Error(err)
	} else {
		defer update.Close()
		defer rpi.DXDisplayClose(display)

		bm, err := update.NewBitmap(bitmap.IMAGE_TYPE_RGB565, 200, 300)
		if err != nil {
			t.Error(err)
		}

		// Clear to red
		bm.ClearToColor(color.RGBA{0xFF, 0, 0, 0xFF})

		// Cross in white
		bm.PaintLine(color.White, bm.NorthWest(), bm.SouthEast())
		bm.PaintLine(color.White, bm.SouthWest(), bm.NorthEast())

		// Place element
		if err := update.Do(0, func() error {
			if _, err := update.AddElement(rpi.DXNewRect(100, 100, 0, 0), bm, 0, 0xFF); err != nil {
				return err
			}
			return nil
		}); err != nil {
			t.Error(err)
		}

		time.Sleep(2 * time.Second)

	}
}

func Test_Update_008(t *testing.T) {
	if display, err := rpi.DXDisplayOpen(DISPLAY_ID); err != nil {
		t.Error(err)
	} else if update, err := bitmap.NewUpdate(display); err != nil {
		t.Error(err)
	} else {
		defer update.Close()
		defer rpi.DXDisplayClose(display)

		bm, err := update.NewBitmap(bitmap.IMAGE_TYPE_RGBA32, 300, 300)
		if err != nil {
			t.Error(err)
		}

		// Clear to transparent
		bm.ClearToColor(color.Transparent)

		// Place element
		if err := update.Do(0, func() error {
			if _, err := update.AddElement(rpi.DXNewRect(100, 100, 0, 0), bm, 0, 0xFF); err != nil {
				return err
			}
			return nil
		}); err != nil {
			t.Error(err)
		}

		// Add circles
		for r := uint32(10); r < uint32(200); r += 10 {
			if err := update.Do(0, func() error {
				bm.PaintCircle(color.RGBA{0xFF, 0, 0, 0xFF}, bm.Centre(), r)
				return nil
			}); err != nil {
				t.Error(err)
			}
			time.Sleep(500 * time.Millisecond)
		}

		time.Sleep(2 * time.Second)
	}
}
