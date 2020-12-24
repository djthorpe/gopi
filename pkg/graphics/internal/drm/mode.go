// +build drm

package drm

import (
	"fmt"

	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Mode struct {
	drm.ModeInfo
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewMode(info drm.ModeInfo) *Mode {
	return &Mode{info}
}

func (this *Mode) Dispose() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Mode) String() string {
	str := "<drm.mode"
	str += " " + fmt.Sprint(this.ModeInfo)
	return str + ">"
}
