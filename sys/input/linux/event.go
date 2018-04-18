// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Input event
type input_event struct {
	device       *device
	timestamp    time.Duration
	device_type  gopi.InputDeviceType
	event_type   gopi.InputEventType
	position     gopi.Point
	rel_position gopi.Point
	key_code     gopi.KeyCode
	scan_code    uint32
	slot         uint
}

////////////////////////////////////////////////////////////////////////////////
// gopi.InputEvent INTERFACE

func (this *input_event) Name() string {
	return "InputEvent"
}

func (this *input_event) Source() gopi.Driver {
	return this.device
}

func (this *input_event) Timestamp() time.Duration {
	return this.timestamp
}

func (this *input_event) DeviceType() gopi.InputDeviceType {
	return this.device_type
}

func (this *input_event) EventType() gopi.InputEventType {
	return this.event_type
}

func (this *input_event) Keycode() gopi.KeyCode {
	return this.key_code
}

func (this *input_event) Scancode() uint32 {
	return this.scan_code
}

func (this *input_event) Position() gopi.Point {
	return this.position
}

func (this *input_event) Relative() gopi.Point {
	return this.rel_position
}

func (this *input_event) Slot() uint {
	return this.slot
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *input_event) String() string {
	switch this.event_type {
	case gopi.INPUT_EVENT_RELPOSITION:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v relative=%v position=%v ts=%v }", this.event_type, this.device_type, this.rel_position, this.position, this.timestamp)
	case gopi.INPUT_EVENT_ABSPOSITION:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v position=%v ts=%v }", this.event_type, this.device_type, this.position, this.timestamp)
	case gopi.INPUT_EVENT_KEYPRESS, gopi.INPUT_EVENT_KEYRELEASE, gopi.INPUT_EVENT_KEYREPEAT:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v key_code=%v scan_code=%v ts=%v }", this.event_type, this.device_type, this.key_code, this.scan_code, this.timestamp)
	case gopi.INPUT_EVENT_TOUCHPRESS, gopi.INPUT_EVENT_TOUCHRELEASE:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v key_code=%v slot=%v position=%v ts=%v }", this.event_type, this.device_type, this.key_code, this.position, this.slot, this.timestamp)
	case gopi.INPUT_EVENT_TOUCHPOSITION:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v slot=%v position=%v ts=%v }", this.event_type, this.device_type, this.position, this.slot, this.timestamp)
	default:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v ts=%v }", this.event_type, this.device_type, this.timestamp)
	}
}
