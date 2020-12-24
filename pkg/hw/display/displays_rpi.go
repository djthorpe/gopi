// +build rpi

package display

import (
	gopi "github.com/djthorpe/gopi/v3"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Displays) Enumerate() []gopi.Display {
	displays := make([]gopi.Display, 0, rpi.TV_MAX_ATTACHED_DISPLAYS)
	for i := uint32(0); i < rpi.TV_MAX_ATTACHED_DISPLAYS; i++ {
		if display, err := this.Open(i); err == nil {
			displays = append(displays, display)
		}
	}
	return displays
}
