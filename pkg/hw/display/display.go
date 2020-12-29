// +build rpi

package display

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/djthorpe/gopi/v3"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type display struct {
	sync.Mutex
	rpi.DXDisplayId
	rpi.TVDisplayInfo
	rpi.DXDisplayHandle
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDisplay(id rpi.DXDisplayId) (*display, error) {
	this := new(display)

	if id == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewDisplay")
	} else {
		this.DXDisplayId = id
	}

	if info, err := rpi.VCHI_TVGetDisplayInfo(id); err != nil {
		return nil, err
	} else {
		this.TVDisplayInfo = info
	}

	return this, nil
}

func (this *display) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error
	if this.DXDisplayHandle != 0 {
		if err := rpi.DXDisplayClose(this.DXDisplayHandle); err != nil {
			result = multierror.Append(result, err)
		}

	}

	this.DXDisplayHandle = 0

	return result
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
	if info, err := rpi.DXDisplayGetInfo(this.DXDisplayHandle); err != nil {
		fmt.Println(err)
		return 0, 0
	} else {
		return info.Size.W, info.Size.H
	}
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
