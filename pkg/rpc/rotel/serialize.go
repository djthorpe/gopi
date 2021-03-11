package rotel

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	rotel "github.com/djthorpe/gopi/v3/pkg/dev/rotel"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	*Event
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - EVENTS

func toProtoNull() *Event {
	return &Event{}
}

func toProtoEvent(evt gopi.RotelEvent) *Event {
	return &Event{}
}

func fromProtoEvent(evt *Event) gopi.RotelEvent {
	if evt == nil {
		return nil
	} else {
		return &event{evt}
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - STATE

func toProtoState(dev gopi.RotelManager) *State {
	dev_, ok := dev.(*rotel.Manager)
	if ok == false {
		return nil
	}
	return &State{
		Model:  dev_.Model(),
		Power:  dev_.Power(),
		Source: dev_.Source(),
		Volume: uint32(dev_.Volume()),
		Muted:  dev_.Muted(),
		Bypass: dev_.Bypass(),
		Bass:   int32(dev_.Bass()),
		Treble: int32(dev_.Treble()),
	}
}

/////////////////////////////////////////////////////////////////////
// EVENT IMPLEMENTATION

func (*event) Name() string {
	return "rotel"
}

func (this *event) String() string {
	str := "<event"
	str += fmt.Sprintf(" name=%q", this.Name())
	str += fmt.Sprintf(" ", this.Event)
	return str + ">"
}
