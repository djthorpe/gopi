package rotel

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	gopi.RotelFlag
	*State
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEvent(flag gopi.RotelFlag, state *State) gopi.RotelEvent {
	return &event{flag, state}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *event) Name() string {
	return this.State.Model()
}

func (this *event) Flags() gopi.RotelFlag {
	return this.RotelFlag
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<rotel.event"
	str += fmt.Sprintf(" name=%q", this.Name())
	if this.RotelFlag != gopi.ROTEL_FLAG_NONE {
		str += fmt.Sprint(" flags=", this.RotelFlag)
	}
	if this.State != nil {
		str += fmt.Sprint(" ", this.State)
	}
	return str + ">"
}
