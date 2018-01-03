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
	device     *device
	timestamp  time.Duration
	event_type gopi.InputEventType
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
	return this.device.device_type
}

func (this *event) EventType() gopi.InputEventType {
	return this.event_type
}

func (this *event) Keycode() gopi.KeyCode {
	return 0
}

func (this *event) Scancode() uint32 {
	return 0
}

func (this *event) Position() gopi.Point {
	return gopi.Point{}
}

func (this *event) Relative() gopi.Point {
	return gopi.Point{}
}

func (this *event) Slot() uint {
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	return fmt.Sprintf("<sys.input.linux.InputEvent>{ type=%v device=%v ts=%v }", this.event_type, this.device.device_type, this.timestamp)
}
