// +build drm

package drm_test

import (
	"testing"

	drm "github.com/djthorpe/gopi/v3/pkg/graphics/internal/drm"
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

		if primary := drm.NewPrimaryPlaneForCrtc(drm.Crtc()); primary != nil {
			t.Log("primary=", primary)
			if err := primary.Dispose(); err != nil {
				t.Error(err)
			}
		} else {
			t.Error("Expected primary plane")
		}

		if cursor := drm.NewCursorPlaneForCrtc(drm.Crtc()); cursor != nil {
			t.Log("cursor=", cursor)
			if err := cursor.Dispose(); err != nil {
				t.Error(err)
			}
		} else {
			t.Error("Expected cursor plane")
		}

		overlays := drm.NewOverlayPlanesForCrtc(drm.Crtc())
		for _, plane := range overlays {
			t.Log("overlay=", plane)
			if err := plane.Dispose(); err != nil {
				t.Error(err)
			}
		}

		if err := drm.Dispose(); err != nil {
			t.Error(err)
		}
	}
}
