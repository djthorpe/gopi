package keycode

import (
	"fmt"
	"strconv"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/lirc/codec"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	name string
	gopi.KeyCode
	*codec.CodecEvent
}

/////////////////////////////////////////////////////////////////////
// INIT

func NewInputEvent(name string, keycode gopi.KeyCode, evt *codec.CodecEvent) *event {
	return &event{name, keycode, evt}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC PROPERTIES

func (this *event) Name() string {
	return this.name
}

func (this *event) Type() gopi.InputType {
	return this.CodecEvent.Type
}

func (this *event) Key() gopi.KeyCode {
	return this.KeyCode
}

func (this *event) Device() (gopi.InputDevice, uint32) {
	return this.CodecEvent.Device | gopi.INPUT_DEVICE_REMOTE, this.CodecEvent.Code
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<event.input"
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if key := this.Key(); key != 0 {
		str += " key=" + fmt.Sprint(key)
	}
	str += " " + fmt.Sprint(this.CodecEvent)
	return str + ">"
}
