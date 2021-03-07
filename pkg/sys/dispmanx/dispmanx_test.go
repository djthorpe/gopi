//+build dispmanx

package dispmanx_test

import (
	"testing"

	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

func Test_DX_001(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Fatal("BCMHostInit Error:", err)
	}
	display, err := dx.DisplayOpen(0)
	if err != nil {
		t.Fatal("DisplayOpen Error:", err)
	} else {
		t.Log("Display[0]=", display)
	}
	if err := dx.DisplayClose(display); err != nil {
		t.Fatal("DisplayClose Error:", err)
	}
}

func Test_DX_002(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Fatal("BCMHostInit Error:", err)
	}
	display, err := dx.DisplayOpen(0)
	if err != nil {
		t.Fatal("DisplayOpen Error:", err)
	} else {
		t.Log("Display[0]=", display)
	}
	if info, err := dx.DisplayGetInfo(display); err != nil {
		t.Error("DisplayGetInfo error", err)
	} else {
		t.Log("Info=", info)
	}
	if err := dx.DisplayClose(display); err != nil {
		t.Fatal("DisplayClose Error:", err)
	}
}

func Test_DX_003(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Fatal("BCMHostInit Error:", err)
	}
	display, err := dx.DisplayOpen(0)
	if err != nil {
		t.Fatal("DisplayOpen Error:", err)
	} else {
		t.Log("Display[0]=", display)
	}
	bitmap, err := dx.DisplaySnapshot(display)
	if err != nil {
		t.Error("DisplayGetInfo error", err)
	} else if err := dx.ResourceDelete(bitmap); err != nil {
		t.Error("ResourceDelete error", err)
	}
	if err := dx.DisplayClose(display); err != nil {
		t.Fatal("DisplayClose Error:", err)
	}
}
