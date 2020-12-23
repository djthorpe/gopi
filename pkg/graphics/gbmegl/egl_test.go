// +build drm,gbm

package gbmegl_test

import (
	"context"
	"testing"
	"time"

	drm "github.com/djthorpe/gopi/v3/pkg/graphics/drm"
	gbmegl "github.com/djthorpe/gopi/v3/pkg/graphics/gbmegl"
)

func Test_EGL_000(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping, no primary device")
	}
	defer fh.Close()
	gbm, err := gbmegl.NewGBM(fh.Fd())
	if err != nil {
		t.Fatal(err)
	}
	egl, err := gbmegl.NewEGL(gbm)
	if err != nil {
		gbm.Dispose()
		t.Fatal(err)
	}

	t.Log("gbm=", gbm)
	t.Log("egl=", egl)

	if err := egl.Dispose(); err != nil {
		t.Error(err)
	}
	if err := gbm.Dispose(); err != nil {
		t.Error(err)
	}
}

func Test_EGL_001(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping, no primary device")
	}
	defer fh.Close()
	gbm, err := gbmegl.NewGBM(fh.Fd())
	if err != nil {
		t.Fatal(err)
	}
	egl, err := gbmegl.NewEGL(gbm)
	if err != nil {
		gbm.Dispose()
		t.Fatal(err)
	}

	for _, api := range egl.API() {
		if err := egl.BindAPI(api); err != nil {
			t.Error(err)
		} else if bound, err := egl.BoundAPI(); err != nil {
			t.Error(err)
		} else if bound != api {
			t.Error("Unexpected bind", bound, " (expected", api, ")")
		} else {
			t.Log(api, "=>", bound)
		}
	}

	if err := egl.Dispose(); err != nil {
		t.Error(err)
	}
	if err := gbm.Dispose(); err != nil {
		t.Error(err)
	}
}

func Test_EGL_002(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping, no primary device")
	}
	defer fh.Close()
	gbm, err := gbmegl.NewGBM(fh.Fd())
	if err != nil {
		t.Fatal(err)
	}
	egl, err := gbmegl.NewEGL(gbm)
	if err != nil {
		gbm.Dispose()
		t.Fatal(err)
	}

	for _, api := range egl.API() {
		if err := egl.BindAPI(api); err != nil {
			t.Error(err)
		}
		if config, context, err := egl.CreateContextForSurface(api, 1, 8, 8, 8, 8); err != nil {
			t.Error(err)
		} else {
			t.Log("config=", config)
			t.Log("context=", context)
		}
	}

	if err := egl.Dispose(); err != nil {
		t.Error(err)
	}
	if err := gbm.Dispose(); err != nil {
		t.Error(err)
	}
}

func Test_EGL_003(t *testing.T) {
	drm, err := drm.NewDRM("", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer drm.Dispose()
	gbm, err := gbmegl.NewGBM(drm.Fd())
	if err != nil {
		t.Fatal(err)
	}
	defer gbm.Dispose()
	egl, err := gbmegl.NewEGL(gbm)
	if err != nil {
		gbm.Dispose()
		t.Fatal(err)
	}
	defer egl.Dispose()

	for _, api := range egl.API() {
		if err := egl.BindAPI(api); err != nil {
			t.Error(err)
		}
		if surface, err := egl.CreateSurface(api, 1, 1920, 1080, gbmegl.GBM_BO_FORMAT_XRGB8888); err != nil {
			t.Log(err)
		} else {
			t.Log("surface=", surface)
			if err := egl.DestroySurface(surface); err != nil {
				t.Error(err)
			}
		}
	}
}

func Test_EGL_004(t *testing.T) {
	drm, err := drm.NewDRM("", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer drm.Dispose()
	gbm, err := gbmegl.NewGBM(drm.Fd())
	if err != nil {
		t.Fatal(err)
	}
	defer gbm.Dispose()
	egl, err := gbmegl.NewEGL(gbm)
	if err != nil {
		gbm.Dispose()
		t.Fatal(err)
	}
	defer egl.Dispose()

	for _, api := range egl.API() {
		if err := egl.BindAPI(api); err != nil {
			t.Error(err)
		}
		if surface, err := egl.CreateSurface(api, 1, 1920, 1080, gbmegl.GBM_BO_FORMAT_XRGB8888); err != nil {
			t.Log(err)
		} else {
			t.Log("surface=", surface)
			if err := egl.DestroySurface(surface); err != nil {
				t.Error(err)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := egl.Run(ctx); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		t.Error(err)
	}
}
