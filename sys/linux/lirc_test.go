// +build linux

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

func Test_LIRC_000(t *testing.T) {
	t.Log("Test_LIRC_000")
}
func Test_LIRC_001(t *testing.T) {
	for f := linux.LIRC_FEATURES_MIN; f <= linux.LIRC_FEATURES_MAX; f <<= 1 {
		t.Log(f)
	}
}
