// +build drm

package drm

import (
	"os"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func OpenPrimaryDevice() (*os.File, error) {
	for _, device := range drm.Devices() {
		if path := device.AvailableNode(drm.DRM_NODE_PRIMARY); path != "" {
			if fh, err := drm.OpenDevice(device, drm.DRM_NODE_PRIMARY); err == nil {
				return fh, nil
			}
		}
	}
	return nil, gopi.ErrNotFound.WithPrefix("OpenPrimaryDevice")
}
