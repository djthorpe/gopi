// +build rpi

package display

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/djthorpe/gopi/v3"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

type display struct {
	sync.Mutex
	rpi.DXDisplayId
	rpi.TVDisplayInfo
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDisplay(id rpi.DXDisplayId) *display {
	this := new(display)

	if id == 0 {
		return nil
	} else {
		this.DXDisplayId = id
	}

	if info, err := rpi.VCHI_TVGetDisplayInfo(id); err != nil {
		return nil
	} else {
		this.TVDisplayInfo = info
	}

	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *display) Id() uint32 {
	return uint32(this.DXDisplayId)
}

func (this *display) Name() string {
	return this.TVDisplayInfo.Product()
}

func (this *display) Flags() gopi.DisplayFlag {
	flags := gopi.DISPLAY_FLAG_NONE
	state, err := rpi.VCHI_TVGetDisplayState(this.DXDisplayId)
	if err != nil {
		return flags
	}
	if state.Flags()&rpi.TV_STATE_HDMI_UNPLUGGED == rpi.TV_STATE_HDMI_UNPLUGGED {
		flags |= gopi.DISPLAY_FLAG_UNPLUGGED
	}
	if state.Flags()&rpi.TV_STATE_HDMI_ATTACHED == rpi.TV_STATE_HDMI_ATTACHED {
		flags |= gopi.DISPLAY_FLAG_ATTACHED
	}
	if state.Flags()&rpi.TV_STATE_HDMI_DVI == rpi.TV_STATE_HDMI_DVI {
		flags |= gopi.DISPLAY_FLAG_DVI
	}
	if state.Flags()&rpi.TV_STATE_HDMI_HDMI == rpi.TV_STATE_HDMI_HDMI {
		flags |= gopi.DISPLAY_FLAG_HDMI
	}
	if state.Flags()&rpi.TV_STATE_SDTV_UNPLUGGED == rpi.TV_STATE_SDTV_UNPLUGGED {
		flags |= gopi.DISPLAY_FLAG_UNPLUGGED
	}
	if state.Flags()&rpi.TV_STATE_SDTV_ATTACHED == rpi.TV_STATE_SDTV_ATTACHED {
		flags |= gopi.DISPLAY_FLAG_ATTACHED
	}
	if state.Flags()&rpi.TV_STATE_SDTV_NTSC == rpi.TV_STATE_SDTV_NTSC {
		flags |= gopi.DISPLAY_FLAG_NTSC | gopi.DISPLAY_FLAG_SDTV
	}
	if state.Flags()&rpi.TV_STATE_SDTV_PAL == rpi.TV_STATE_SDTV_PAL {
		flags |= gopi.DISPLAY_FLAG_PAL | gopi.DISPLAY_FLAG_SDTV
	}
	if state.Flags()&rpi.TV_STATE_LCD_ATTACHED_DEFAULT == rpi.TV_STATE_LCD_ATTACHED_DEFAULT {
		flags |= gopi.DISPLAY_FLAG_LCD
	}
	return flags
}

func (this *display) Size() (uint32, uint32) {
	return rpi.BCMGetDisplaySize(uint16(this.Id()))
}

func (this *display) PixelsPerInch() uint32 {
	return 0
}

func (this *display) Vendor() string {
	return this.TVDisplayInfo.Vendor()
}

func (this *display) Serial() string {
	return fmt.Sprint(this.TVDisplayInfo.Serial())
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *display) String() string {
	str := "<display"
	str += " id=" + fmt.Sprint(this.Id())
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if vendor := this.Vendor(); vendor != "" {
		str += " vendor=" + strconv.Quote(vendor)
	}
	if serial := this.Serial(); serial != "" {
		str += " serial=" + strconv.Quote(serial)
	}
	if flags := this.Flags(); flags != gopi.DISPLAY_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	if x, y := this.Size(); x > 0 && y > 0 {
		str += " size={" + fmt.Sprint(x, ",", y) + "}"
	}
	if ppi := this.PixelsPerInch(); ppi > 0 {
		str += " pixels_per_inch=" + fmt.Sprint(ppi)
	}
	return str + ">"
}
