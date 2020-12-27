package display

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

type event struct {
	display gopi.Display
	flags   gopi.DisplayFlag
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEvent(display gopi.Display, flags gopi.DisplayFlag) *event {
	return &event{display, flags}
}

func (this *event) Name() string {
	return this.display.Name()
}

func (this *event) Value() interface{} {
	return this.flags
}

func (this *event) String() string {
	str := "<display.event"
	if this.display != nil {
		str += " display=" + fmt.Sprint(this.display)
	}
	if this.flags != 0 {
		str += " flags=" + fmt.Sprint(this.flags)
	}
	return str + ">"
}
