// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap_test

import (
	"encoding/hex"
	"image/color"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
)

func init() {
	rpi.DXInit()
}

func Test_Paint_000(t *testing.T) {
	t.Log("Test_Paint_000")
}

func Test_Paint_001(t *testing.T) {
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{1, 1}}, nil); err != nil {
		t.Error("Unexpected error return", err)
	} else {
		bm.(bitmap.Bitmap).ClearToColor(color.White)
		if bytes, stride := bm.(bitmap.Bitmap).Bytes(); bytes == nil {
			t.Error("Bytes failed")
		} else {
			t.Log("bytes=", hex.EncodeToString(bytes), "stride=", stride)
		}
	}
}
