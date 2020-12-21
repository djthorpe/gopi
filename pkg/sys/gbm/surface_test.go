// +build gbm

package gbm_test

import (
	"strings"
	"testing"

	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
)

func Test_Surface_001(t *testing.T) {
	for _, node := range gbm.GBMDevices() {
		if strings.HasPrefix(node, "card") == false {
			t.Log("Skipping", node)
		} else if fh, err := gbm.OpenDevice(node); err != nil {
			t.Error(node, err)
		} else {
			defer fh.Close()
			device, err := gbm.GBMCreateDevice(fh.Fd())
			if err != nil {
				t.Fatal(node, "Unable to create device", err)
			}
			defer device.Free()
			if surface, err := device.SurfaceCreate(100, 100, gbm.GBM_FORMAT_RGB888, gbm.GBM_BO_USE_SCANOUT); err != nil {
				t.Error(err, "for GPU", node)
			} else {
				t.Log(surface)
				surface.Free()
			}
		}
	}
}
