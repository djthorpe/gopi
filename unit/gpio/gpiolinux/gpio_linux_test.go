// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiolinux_test

import (
	"testing"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/gpio/gpiolinux"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_GPIO_Linux_000(t *testing.T) {
	t.Log("Test_GPIO_Linux_000")
}
