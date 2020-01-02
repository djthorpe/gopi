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

func Test_Model_000(t *testing.T) {
	t.Log("Test_Model_000")
}

func Test_Model_001(t *testing.T) {
	if _, product, err := rpi.VCGetSerialProduct(); err != nil {
		t.Error(err)
	} else {
		product := rpi.NewProductInfo(product)
		t.Log(product)
	}
}
