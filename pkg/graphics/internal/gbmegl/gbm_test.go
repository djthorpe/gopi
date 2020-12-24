// +build drm,gbm

package gbmegl_test

import (
	"testing"

	drm "github.com/djthorpe/gopi/v3/pkg/graphics/internal/drm"
	gbmegl "github.com/djthorpe/gopi/v3/pkg/graphics/internal/gbmegl"
)

func Test_GBM_000(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping, no primary device")
	}
	defer fh.Close()
	if gbm, err := gbmegl.NewGBM(fh.Fd()); err != nil {
		t.Error(err)
	} else {
		t.Log(fh.Name(), gbm)
		if err := gbm.Dispose(); err != nil {
			t.Error(err)
		}
	}
}
