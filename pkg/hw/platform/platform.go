package platform

import (
	"fmt"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Platform struct {
	gopi.Unit
	*Implementation // Implementation-specific members
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Platform) String() string {
	str := "<platform"
	if p := this.Product(); p != "" {
		str += fmt.Sprintf(" product=%q", p)
	}
	if t := this.Type(); t != gopi.PLATFORM_NONE {
		str += " type=" + fmt.Sprint(this.Type())
	}
	if sn := this.SerialNumber(); sn != "" {
		str += " serial_number=" + fmt.Sprint(sn)
	}
	if ut := this.Uptime(); ut != 0 {
		str += " uptime=" + fmt.Sprint(ut.Truncate(time.Second))
	}
	if av1, av5, av15 := this.LoadAverages(); av1 != 0 || av5 != 0 || av15 != 0 {
		str += fmt.Sprintf(" load_avg={ %.2f, %.2f, %.2f }", av1, av5, av15)
	}
	if d := this.NumberOfDisplays(); d != 0 {
		str += " number_displays=" + fmt.Sprint(d)
	}
	if a := this.AttachedDisplays(); len(a) != 0 {
		str += " attached_displays=" + fmt.Sprint(a)
	}
	return str + ">"
}
