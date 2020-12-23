// +build drm

package drm_test

import (
	"testing"

	drm "github.com/djthorpe/gopi/v3/pkg/graphics/drm"
)

func Test_Resources_000(t *testing.T) {
	t.Log("Test", t.Name())
	if res, err := drm.NewResources(0); err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
}
