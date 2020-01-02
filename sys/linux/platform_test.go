// +build linux,rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/sys/linux"
)

func Test_Platform_000(t *testing.T) {
	t.Log("Test_Platform_000")
}

func Test_Platform_001(t *testing.T) {
	if serial := linux.SerialNumber(); serial == "" {
		t.Error("Unexpected response from SerialNumber")
	} else {
		t.Log("serial", serial)
	}
}
