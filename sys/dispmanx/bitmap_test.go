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

	// Frameworks
	bitmap "github.com/djthorpe/gopi/v2/sys/dispmanx"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

func init() {
	rpi.DXInit()
}

func Test_Bitmap_000(t *testing.T) {
	t.Log("Test_Bitmap_000")
}

func Test_Bitmap_001(t *testing.T) {
	for _, imageType := range []bitmap.ImageType{
		bitmap.IMAGE_TYPE_RGB565, bitmap.IMAGE_TYPE_RGB888, bitmap.IMAGE_TYPE_RGBA32,
	} {
		bm, err := bitmap.NewBitmap(imageType, 100, 100)
		if err != nil {
			t.Error(err)
		}
		t.Log(bm)
		if err := bm.Close(); err != nil {
			t.Error(err)
		}
	}
}

func Test_Bitmap_002(t *testing.T) {
	for _, imageType := range []bitmap.ImageType{
		bitmap.IMAGE_TYPE_RGB565, bitmap.IMAGE_TYPE_RGB888, bitmap.IMAGE_TYPE_RGBA32,
	} {
		bm, err := bitmap.NewBitmap(imageType, 100, 100)
		if err != nil {
			t.Error(err)
		}
		if err := bm.ClearToColor(color.White); err != nil {
			t.Error(err)
		}
		if err := bm.Close(); err != nil {
			t.Error(err)
		}
	}
}
