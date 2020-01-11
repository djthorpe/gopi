// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/sys/rpi"
)

func Test_Display_000(t *testing.T) {
	t.Log("Test_Display_000")
}

func Test_Display_001(t *testing.T) {
	rpi.DXInit()
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log(display)
	}
}

func Test_Display_002(t *testing.T) {
	rpi.DXInit()
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if info, err := rpi.DXDisplayGetInfo(display); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log(info)
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST RECT

func Test_Rect_001(t *testing.T) {
	r := rpi.DXNewRect(0, 0, 0, 0)
	if size := rpi.DXRectSize(r); size.W != 0 || size.H != 0 {
		t.Error("Unexpected values for rect size")
	} else if origin := rpi.DXRectOrigin(r); origin.X != 0 || origin.Y != 0 {
		t.Error("Unexpected values for rect size")
	} else {
		t.Log("rect", rpi.DXRectString(r))
		t.Log("size", size)
		t.Log("origin", origin)
	}
}
func Test_Rect_002(t *testing.T) {
	r := rpi.DXNewRect(-100, -99, 100, 99)
	if size := rpi.DXRectSize(r); size.W != 100 || size.H != 99 {
		t.Error("Unexpected values for rect size")
	} else if origin := rpi.DXRectOrigin(r); origin.X != -100 || origin.Y != -99 {
		t.Error("Unexpected values for rect size")
	} else {
		t.Log("rect", rpi.DXRectString(r))
		t.Log("size", size)
		t.Log("origin", origin)
	}
}

func Test_Rect_003(t *testing.T) {
	r := rpi.DXNewRect(0, 0, 0, 0)
	if err := rpi.DXRectSet(r, -100, -99, 100, 99); err != nil {
		t.Error(err)
	} else if size := rpi.DXRectSize(r); size.W != 100 || size.H != 99 {
		t.Error("Unexpected values for rect size")
	} else if origin := rpi.DXRectOrigin(r); origin.X != -100 || origin.Y != -99 {
		t.Error("Unexpected values for rect size")
	} else {
		t.Log("rect", rpi.DXRectString(r))
		t.Log("size", size)
		t.Log("origin", origin)
	}
}

func Test_Rect_004(t *testing.T) {
	r1 := rpi.DXNewRect(0, 0, 10, 10)
	r2 := rpi.DXNewRect(-10, -10, 20, 20)
	r3 := rpi.DXRectIntersection(r1, r2)
	t.Log("r1", rpi.DXRectString(r1))
	t.Log("r2", rpi.DXRectString(r2))
	t.Log("r1 u r2", rpi.DXRectString(r3))
	if size := rpi.DXRectSize(r3); size.W != 10 || size.H != 10 {
		t.Error("Expected intersection of rectangles to be of size 10")
	}
}

func Test_Rect_005(t *testing.T) {
	r1 := rpi.DXNewRect(0, 0, 10, 10)
	r2 := rpi.DXNewRect(-10, -10, 50, 50)
	r3 := rpi.DXRectIntersection(r1, r2)
	t.Log("r1", rpi.DXRectString(r1))
	t.Log("r2", rpi.DXRectString(r2))
	t.Log("r1 u r2", rpi.DXRectString(r3))
	if size := rpi.DXRectSize(r3); size.W != 10 || size.H != 10 {
		t.Error("Expected intersection of rectangles to be of size 10")
	}
}

////////////////////////////////////////////////////////////////////////////////
// TEST RESOURCES

func Test_Resources_001(t *testing.T) {
	rpi.DXInit()
	if resource, err := rpi.DXResourceCreate(rpi.DX_IMAGE_TYPE_RGBA32, rpi.DXSize{100, 100}); err != nil {
		t.Error(err)
	} else if err := rpi.DXResourceDelete(resource); err != nil {
		t.Error(err)
	} else {
		t.Log(resource)
	}
}
