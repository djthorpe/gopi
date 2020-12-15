package input

import (
	"fmt"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// INPUT EVENT INTERFACE

type event struct {
	*Event
}

func (e *event) Name() string {
	return e.GetName()
}

func (e *event) Type() gopi.InputType {
	return gopi.InputType(e.GetType())
}

func (e *event) Key() gopi.KeyCode {
	return gopi.KeyCode(e.GetKey())
}

func (e *event) Device() (gopi.InputDeviceType, uint32) {
	return gopi.InputDevice(e.GetDevice()), e.GetScancode()
}

func (e *event) String() string {
	str := "<event.input"
	if n := e.Name(); n != "" {
		str += " name=" + strconv.Quote(n)
	}
	if t := e.Type(); t != gopi.INPUT_EVENT_NONE {
		str += " type=" + fmt.Sprint(t)
	}
	if k := e.Key(); k != gopi.KEYCODE_NONE {
		str += " key=" + fmt.Sprint(k)
	}
	if d, s := e.Device(); d != gopi.INPUT_DEVICE_NONE {
		str += " device=" + fmt.Sprint(d)
		str += " scancode=" + fmt.Sprintf("0x%08X", s)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - INPUT EVENT

func protoFromInputEvent(evt gopi.InputEvent) *Event {
	if evt == nil {
		return &Event{}
	}
	// TODO: Timestamp
	device, scancode := evt.Device()
	return &Event{
		Name:     evt.Name(),
		Key:      Event_KeyCode(evt.Key()),
		Type:     Event_InputType(evt.Type()),
		Device:   Event_DeviceType(device),
		Scancode: uint32(scancode),
	}
}

func protoToInputEvent(pb *Event) gopi.InputEvent {
	if pb == nil {
		return nil
	} else if pb.Name == "" || pb.Type == 0 {
		return nil
	} else {
		return &event{pb}
	}
}
