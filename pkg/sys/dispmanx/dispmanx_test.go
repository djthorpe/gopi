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
	rect := dx.NewRect(0, 0, 100, 100)
	if x, y := rect.Origin(); x != 0 || y != 0 {
		t.Error("Unexpected origin: ", rect)
	}
	if w, h := rect.Size(); w != 100 || h != 100 {
		t.Error("Unexpected size: ", rect)
	}
}

func Test_DX_003(t *testing.T) {
	rect := dx.NewRect(-100, -100, 0, 0)
	if x, y := rect.Origin(); x != -100 || y != -100 {
		t.Error("Unexpected origin: ", rect)
	}
	if w, h := rect.Size(); w != 0 || h != 0 {
		t.Error("Unexpected size: ", rect)
	}
}
func Test_DX_004(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Fatal("BCMHostInit Error:", err)
	}
	display, err := dx.DisplayOpen(0)
	if err != nil {
		t.Fatal("DisplayOpen Error:", err)
	}

	if update, err := dx.UpdateStart(0); err != nil {
		t.Error("UpdateStart:", err)
	} else if err := dx.UpdateSubmitSync(update); err != nil {
		t.Error("UpdateSubmitSync:", err)
	} else {
		t.Log("update=", update)
	}

	if err := dx.DisplayClose(display); err != nil {
		t.Fatal("DisplayClose Error:", err)
	}
}

func Test_DX_005(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Fatal("BCMHostInit Error:", err)
	}
	display, err := dx.DisplayOpen(0)
	if err != nil {
		t.Fatal("DisplayOpen Error:", err)
	}

	if resource, err := dx.ResourceCreate(dx.VC_IMAGE_RGB565, 100, 100); err != nil {
		t.Error("ResourceCreate:", err)
	} else {
		t.Log("Resource=", resource)
		dest := dx.NewRect(0, 0, 100, 100)
		src := dx.NewRect(0, 0, 100, 100)

		if update, err := dx.UpdateStart(0); err != nil {
			t.Error("UpdateStart:", err)
		} else if element, err := dx.ElementAdd(update, display, 100, dest, resource, src, 0, dx.NewAlphaFromSource(), nil, dx.DISPMANX_NO_ROTATE); err != nil {
			t.Error("ElementAdd:", err)
		} else if err := dx.UpdateSubmitSync(update); err != nil {
			t.Error("UpdateSubmitSync:", err)
		} else if update, err := dx.UpdateStart(0); err != nil {
			t.Error("UpdateStart:", err)
		} else if err := dx.ElementRemove(update, element); err != nil {
			t.Error("ElementRemove:", err)
		} else if err := dx.UpdateSubmitSync(update); err != nil {
			t.Error("UpdateSubmitSync:", err)
		}

		if err := dx.ResourceDelete(resource); err != nil {
			t.Error("ResourceDelete:", err)
		}
	}

	if err := dx.DisplayClose(display); err != nil {
		t.Fatal("DisplayClose Error:", err)
	}
}
