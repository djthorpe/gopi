package lirc

import (
	"fmt"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	name  string
	mode  gopi.LIRCMode
	value uint32
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func NewEvent(name string, mode gopi.LIRCMode, value uint32) gopi.LIRCEvent {
	this := new(event)
	this.name = name
	this.mode = mode
	this.value = value
	return this
}

func (this *event) Name() string {
	return this.name
}

func (this *event) Type() gopi.LIRCType {
	return gopi.LIRCType(this.value & 0xFF000000)
}

func (this *event) Mode() gopi.LIRCMode {
	return this.mode
}

func (this *event) Value() interface{} {
	return this.value & 0x00FFFFFF
}

func (this *event) String() string {
	str := "<lircevent"
	if name := this.name; name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if this.mode != gopi.LIRC_MODE_NONE {
		str += " mode=" + fmt.Sprint(this.Mode())
	}
	if this.value != 0 {
		str += " type=" + fmt.Sprint(this.Type())
		str += " value=" + fmt.Sprint(this.Value())
	}
	return str + ">"
}
