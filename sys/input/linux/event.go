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
type event struct {
	device       *device
	timestamp    time.Duration
	device_type  gopi.InputDeviceType
	event_type   gopi.InputEventType
	position     gopi.Point
	rel_position gopi.Point
	key_code     gopi.KeyCode
	scan_code    uint32
}

////////////////////////////////////////////////////////////////////////////////
// gopi.InputEvent INTERFACE

func (this *event) Name() string {
	return "InputEvent"
}

func (this *event) Source() gopi.Driver {
	return this.device
}

func (this *event) Timestamp() time.Duration {
	return this.timestamp
}

func (this *event) DeviceType() gopi.InputDeviceType {
	return this.device_type
}

func (this *event) EventType() gopi.InputEventType {
	return this.event_type
}

func (this *event) Keycode() gopi.KeyCode {
	return this.key_code
}

func (this *event) Scancode() uint32 {
	return this.scan_code
}

func (this *event) Position() gopi.Point {
	return this.position
}

func (this *event) Relative() gopi.Point {
	return this.rel_position
}

func (this *event) Slot() uint {
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

/*
	INPUT_EVENT_TOUCHPRESS    InputEventType = 0x0006
	INPUT_EVENT_TOUCHRELEASE  InputEventType = 0x0007
	INPUT_EVENT_TOUCHPOSITION InputEventType = 0x0008
*/

func (this *event) String() string {
	switch this.event_type {
	case gopi.INPUT_EVENT_RELPOSITION:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v relative=%v position=%v ts=%v }", this.event_type, this.device_type, this.rel_position, this.position, this.timestamp)
	case gopi.INPUT_EVENT_ABSPOSITION:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v position=%v ts=%v }", this.event_type, this.device_type, this.position, this.timestamp)
	case gopi.INPUT_EVENT_KEYPRESS, gopi.INPUT_EVENT_KEYRELEASE, gopi.INPUT_EVENT_KEYREPEAT:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v key_code=%v scan_code=%v ts=%v }", this.event_type, this.device_type, this.key_code, this.scan_code, this.timestamp)
	default:
		return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v ts=%v }", this.event_type, this.device_type, this.timestamp)
	}
}
