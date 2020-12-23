// +build drm

package drm_test

import (
	"testing"

	drm "github.com/djthorpe/gopi/v3/pkg/graphics/drm"
)

func Test_DRM_000(t *testing.T) {
	if drm, err := drm.NewDRM("", 0); err != nil {
		t.Error(err)
	} else {
		t.Log(drm)
		if err := drm.Dispose(); err != nil {
			t.Error(err)
		}
	}
}

func Test_DRM_001(t *testing.T) {
	if drm, err := drm.NewDRM("", 0); err != nil {
		t.Error(err)
	} else {
		t.Log(drm)
		for _, plane := range drm.NewPlanes() {
			t.Log(plane)
			if err := plane.Dispose(); err != nil {
				t.Error(err)
			}
		}

		if err := drm.Dispose(); err != nil {
			t.Error(err)
		}
	}
}
