// +build rpi
// +build !darwin

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

func Test_Platform_000(t *testing.T) {
	t.Log("Test_Platform_000")
}

func Test_Platform_001(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Error("Unexpected response from BCMHostInit")
	} else if err := rpi.BCMHostTerminate(); err != nil {
		t.Error("Unexpected response from BCMHostTerminate")
	}
}

func Test_Platform_002(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Error("Unexpected response from BCMHostInit")
	} else {
		defer rpi.BCMHostTerminate()
		if addr := rpi.BCMHostGetPeripheralAddress(); addr == 0 {
			t.Error("Unexpected response from BCMHostGetPeripheralAddress")
		} else {
			t.Logf("BCMHostGetPeripheralAddress => %08X", addr)
		}
		if size := rpi.BCMHostGetPeripheralSize(); size == 0 {
			t.Error("Unexpected response from BCMHostGetPeripheralSize")
		} else {
			t.Logf("BCMHostGetPeripheralSize => %08X", size)
		}
		if addr := rpi.BCMHostGetSDRAMAddress(); addr == 0 {
			t.Error("Unexpected response from BCMHostGetSDRAMAddress")
		} else {
			t.Logf("BCMHostGetSDRAMAddress => %08X", addr)
		}
	}
}

func Test_Platform_003(t *testing.T) {
	if service := rpi.VCGencmdInit(); service < 0 {
		t.Error("Unexpected response from VCGencmdInit")
	} else {
		t.Logf("VCGencmdInit => %08X", service)
	}
	rpi.VCGencmdTerminate()
}

func Test_Platform_004(t *testing.T) {
	if service := rpi.VCGencmdInit(); service < 0 {
		t.Error("Unexpected response from VCGencmdInit")
	} else if dump, err := rpi.VCGeneralCommand("commands"); err != nil {
		t.Error("Unexpected response from VCGeneralCommand", err)
	} else {
		t.Logf("VCGeneralCommand => %v", dump)
		rpi.VCGencmdTerminate()
	}
}
