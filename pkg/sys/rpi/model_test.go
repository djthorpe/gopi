// +build rpi
// +build !darwin

package rpi_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v3/pkg/sys/rpi"
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
