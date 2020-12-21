// +build gbm

package gbm_test

import (
	"strings"
	"testing"

	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
)

func Test_Buffer_001(t *testing.T) {
	for _, node := range gbm.GBMDevices() {
		if strings.HasPrefix(node, "card") == false {
			t.Log("Skipping", node)
		} else if fh, err := gbm.OpenDevice(node); err != nil {
			t.Error(node, err)
		} else {
			defer fh.Close()
			if device, err := gbm.GBMCreateDevice(fh.Fd()); err != nil {
				t.Error(node, err)
			} else {
				defer device.Free()
				if buffer, err := device.BufferCreate(100, 100, gbm.GBM_FORMAT_XRGB8888, 0); err != nil {
					t.Error(err, "for GPU", node)
				} else {
					t.Log(buffer)
					buffer.Free()
				}
			}
		}
	}
}
