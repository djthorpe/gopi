// +build drm

package drm_test

import (
	"strings"
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

func Test_DRM_000(t *testing.T) {
	devices := drm.Devices()
	if devices == nil {
		t.Error("Unexpected error with drm.Devices call")
	}
	for _, node := range devices {
		if strings.HasPrefix(node, "card") == false {
			t.Log("Skipping", node)
			continue
		}
		if fh, err := drm.OpenDevice(node); err != nil {
			t.Error(err)
		} else if err := fh.Close(); err != nil {
			t.Error(err)
		} else {
			t.Log("Opened device", node)
		}
	}
}

func Test_DRM_001(t *testing.T) {
	fh, err := drm.OpenDevice("card1")
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
	fh, err := drm.OpenDevice("card1")
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
