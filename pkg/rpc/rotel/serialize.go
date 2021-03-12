package rotel

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	rotel "github.com/djthorpe/gopi/v3/pkg/dev/rotel"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	pb *Event
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - EVENTS

func toProtoNull() *Event {
	return &Event{}
}

func toProtoEvent(evt gopi.RotelEvent, dev gopi.RotelManager) *Event {
	return &Event{
		Name:  evt.Name(),
		Flags: Event_Flag(evt.Flags()),
		State: toProtoState(dev),
	}
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
		Power:   dev_.Power(),
		Source:  dev_.Source(),
		Volume:  uint32(dev_.Volume()),
		Muted:   dev_.Muted(),
		Bypass:  dev_.Bypass(),
		Bass:    int32(dev_.Bass()),
		Treble:  int32(dev_.Treble()),
		Balance: toProtoBalance(dev_.Balance()),
		Dimmer:  uint32(dev_.Dimmer()),
	}
}

func toProtoBalance(location string, value uint) *Balance {
	return &Balance{
		Location: location,
		Value:    uint32(value),
	}
}

/////////////////////////////////////////////////////////////////////
// EVENT IMPLEMENTATION

func (this *event) Name() string {
	return this.pb.Name
}

func (this *event) Flags() gopi.RotelFlag {
	return gopi.RotelFlag(this.pb.Flags)
}

func (this *event) String() string {
	str := "<event"
	if name := this.Name(); name != "" {
		str += fmt.Sprintf(" name=%q", this.Name())
	}
	if flags := this.Flags(); flags != gopi.ROTEL_FLAG_NONE {
		str += fmt.Sprint(" flags=", this.Flags())
	}
	if state := this.pb.State; state != nil {
		str += fmt.Sprint(" state=", state)
	}
	return str + ">"
}
