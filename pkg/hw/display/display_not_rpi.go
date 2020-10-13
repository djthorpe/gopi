// +build !rpi

package display

import (
	"fmt"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
)

type display struct {
	gopi.Unit

	id uint16
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *display) Id() uint16 {
	return this.id
}

func (this *display) Name() string {
	return ""
}

func (this *display) Size() (uint32, uint32) {
	return 0, 0
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
