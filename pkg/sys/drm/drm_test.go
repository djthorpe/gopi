// +build drm

package drm_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

func Test_DRM_000(t *testing.T) {
	if devices := drm.Devices(); devices == nil {
		t.Error("Unexpected error with drm.Devices call")
	} else {
		for _, bus := range devices {
			if fh, err := drm.OpenDevice(bus); err != nil {
				t.Error(err)
			} else if err := fh.Close(); err != nil {
				t.Error(err)
			} else {
				t.Log("Opened device", bus)
			}
		}
	}
}

func Test_DRM_001(t *testing.T) {
	fh, err := drm.OpenDevice(1)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if err := drm.GetResources(fh.Fd()); err != nil {
		t.Error(err)
	}
}
