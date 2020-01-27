// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap_test

import (
	"errors"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
)

func init() {
	rpi.DXInit()
}

func Test_Bitmap_000(t *testing.T) {
	t.Log("Test_Bitmap_000")
}

func Test_Bitmap_001(t *testing.T) {
	if _, err := gopi.New(bitmap.Config{}, nil); errors.Is(err, gopi.ErrBadParameter) == false {
		t.Error("Unexpected error return", err)
	}
	if _, err := gopi.New(bitmap.Config{Size: gopi.Size{0, 1}}, nil); errors.Is(err, gopi.ErrBadParameter) == false {
		t.Error("Unexpected error return", err)
	}
	if _, err := gopi.New(bitmap.Config{Size: gopi.Size{1, 0}}, nil); errors.Is(err, gopi.ErrBadParameter) == false {
		t.Error("Unexpected error return", err)
	}
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{1, 1}}, nil); err != nil {
		t.Error("Unexpected error return", err)
	} else if bm.(bitmap.Bitmap).Mode() != gopi.SURFACE_FLAG_RGBA32 {
		t.Error("Unexpected bitmap mode", bm.(bitmap.Bitmap).Mode())
	}
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{1, 1}, Mode: gopi.SURFACE_FLAG_RGB565}, nil); err != nil {
		t.Error("Unexpected error return", err)
	} else {
		t.Log(bm)
	}
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{1, 1}, Mode: gopi.SURFACE_FLAG_RGB888}, nil); err != nil {
		t.Error("Unexpected error return", err)
	} else {
		t.Log(bm)
	}
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{1, 1}, Mode: gopi.SURFACE_FLAG_RGBA32}, nil); err != nil {
		t.Error("Unexpected error return", err)
	} else {
		t.Log(bm)
	}
}

func Test_Bitmap_002(t *testing.T) {
	if bm, err := gopi.New(bitmap.Config{Size: gopi.Size{100, 100}}, nil); err != nil {
		t.Error(err)
	} else if size := bm.(bitmap.Bitmap).Size(); size.W != 100 || size.H != 100 {
		t.Error("Unexpected size:", size)
	} else if err := bm.Close(); err != nil {
		t.Error(err)
	}
}
