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

	if r, err := drm.GetResources(fh.Fd()); err != nil {
		t.Error(err)
	} else {
		defer r.Free()
		t.Log(r)
	}
}

func Test_DRM_002(t *testing.T) {
	fh, err := drm.OpenDevice(1)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	r, err := drm.GetResources(fh.Fd())
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()

	connectors := r.Connectors()
	if len(connectors) == 0 {
		t.Log("Skipping test as no connectors")
		t.SkipNow()
	}
	for _, id := range connectors {
		if connector, err := drm.GetConnector(fh.Fd(), id); err != nil {
			t.Error(err)
		} else {
			defer connector.Free()
			t.Log(connector)
		}
	}
}
