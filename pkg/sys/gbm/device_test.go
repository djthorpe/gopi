// +build gbm

package gbm_test

import (
	"strings"
	"testing"

	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
)

func Test_Device_001(t *testing.T) {
	for _, node := range gbm.GBMDevices() {
		t.Log(node)
	}
}

func Test_Device_002(t *testing.T) {
	for _, node := range gbm.GBMDevices() {
		if strings.HasPrefix(node, "card") == false {
			t.Log("Skipping", node)
		} else if fh, err := gbm.OpenDevice(node); err != nil {
			t.Error(node, err)
		} else {
			defer fh.Close()
			if device, err := gbm.GBMCreateDevice(fh.Fd()); err != nil {
				t.Error(node, "Unable to create device", err)
			} else {
				t.Log(device)
				device.Free()
			}
		}
	}
}
