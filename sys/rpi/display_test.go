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
