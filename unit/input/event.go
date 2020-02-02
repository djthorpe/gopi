/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input

import (
	"fmt"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	source   gopi.InputManager
	device   gopi.InputDevice
	event    gopi.InputEventType
	ts       time.Duration

	rel      gopi.Point
	position gopi.Point

	keycode  gopi.KeyCode
	keystate gopi.KeyState
	scancode uint32
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewRelPositionEvent(rel gopi.Point, position gopi.Point, ts time.Duration) *event {
	return &event{
		event:    gopi.INPUT_EVENT_RELPOSITION,
		rel:      rel,
		position: position,
		ts:       ts,
	}
}

func NewAbsPositionEvent(position gopi.Point,code gopi.KeyCode, ts time.Duration) *event {
	return &event{
		event:    gopi.INPUT_EVENT_ABSPOSITION,
		rel:      gopi.ZeroPoint,
		position: position,
		keycode:  code,
		ts:       ts,
	}
}

func NewKeyEvent(action gopi.KeyAction, code gopi.KeyCode, state gopi.KeyState, scancode uint32, ts time.Duration) *event {
	evt := &event{
		keycode:  code,
		keystate: state,
		scancode: scancode,
		ts:       ts,
	}
	switch action {
	case gopi.KEYACTION_KEY_UP:
		evt.event = gopi.INPUT_EVENT_KEYRELEASE
	case gopi.KEYACTION_KEY_DOWN:
		evt.event = gopi.INPUT_EVENT_KEYPRESS
	case gopi.KEYACTION_KEY_REPEAT:
		evt.event = gopi.INPUT_EVENT_KEYREPEAT
	default:
		return nil
	}
	return evt
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Event

func (*event) Name() string {
	return "gopi.InputEvent"
}

func (*event) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *event) Source() gopi.Unit {
	return this.source
}

func (this *event) Value() interface{} {
	return nil
}

func (this *event) Device() gopi.InputDevice {	
	return this.device
}

func (this *event) Type() gopi.InputEventType {
	return this.event
}

func (this *event) KeyCode() gopi.KeyCode {
	return this.keycode
}

func (this *event) KeyState() gopi.KeyState {
	return this.keystate
}

func (this *event) ScanCode() uint32 {
	return this.scancode
}

	// Abs returns absolute input position
	func (this *event) 	Abs() gopi.Point {
		return this.position
	}

	// Rel returns relative input position
	func (this *event) 	Rel() gopi.Point {
		return this.rel
	}


////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<" + this.Name()
	if this.device != nil {
		str += " device=" + fmt.Sprint(this.device)
	}
	switch this.Type() {
	case gopi.INPUT_EVENT_RELPOSITION:
		str += " rel=" + fmt.Sprint(this.rel)
	case gopi.INPUT_EVENT_ABSPOSITION:
		str += " position=" + fmt.Sprint(this.position)
		if this.keycode != gopi.KEYCODE_NONE {
			str += " keycode=" + fmt.Sprint(this.KeyCode())
		}
		if this.keystate != 0 {
			str += " keystate=" + fmt.Sprint(this.KeyState())
		}
	case gopi.INPUT_EVENT_KEYPRESS, gopi.INPUT_EVENT_KEYRELEASE, gopi.INPUT_EVENT_KEYREPEAT:
		str += " action=" + fmt.Sprint(this.Type())
		if this.keycode != gopi.KEYCODE_NONE {
			str += " keycode=" + fmt.Sprint(this.KeyCode())
		}
		if this.keystate != 0 {
			str += " keystate=" + fmt.Sprint(this.KeyState())
		}
		if this.scancode > 0 {
			str += fmt.Sprintf(" scancode=0x%08X", this.scancode)
		}
	}
	return str + ">"
}
