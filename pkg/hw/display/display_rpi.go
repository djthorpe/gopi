// +build rpi

package display

import (
	"fmt"
	"strconv"
	"strings"

	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
	multierror "github.com/hashicorp/go-multierror"
)

type display struct {
	id     uint16
	handle rpi.DXDisplayHandle
	info   rpi.DXDisplayModeInfo
}

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func (this *display) new(id uint16) error {
	if handle, err := rpi.DXDisplayOpen(rpi.DXDisplayId(id)); err != nil {
		return err
	} else if info, err := rpi.DXDisplayGetInfo(handle); err != nil {
		rpi.DXDisplayClose(handle)
		return err
	} else {
		this.id = id
		this.handle = handle
		this.info = info
	}

	// Return success
	return nil
}

func (this *display) close() error {
	var result error
	if this.handle != rpi.DX_NO_HANDLE {
		if err := rpi.DXDisplayClose(this.handle); err != nil {
			result = multierror.Append(result, err)
		}
		this.handle = rpi.DX_NO_HANDLE
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *display) Id() uint16 {
	return this.id
}

func (this *display) Name() string {
	name := fmt.Sprint(rpi.DXDisplayId(this.id))
	if strings.HasPrefix(name, "DX_DISPLAYID_") {
		return strings.TrimPrefix(name, "DX_DISPLAYID_")
	} else {
		return ""
	}
}

func (this *display) Size() (uint32, uint32) {
	return this.info.Size.W, this.info.Size.H
}

func (this *display) PixelsPerInch() uint32 {
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *display) String() string {
	str := "<display"
	str += " id=" + fmt.Sprint(this.id)
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if x, y := this.Size(); x > 0 && y > 0 {
		str += " size={" + fmt.Sprint(x, ",", y) + "}"
	}
	if ppi := this.PixelsPerInch(); ppi > 0 {
		str += " pixels_per_inch=" + fmt.Sprint(ppi)
	}
	return str + ">"
}
