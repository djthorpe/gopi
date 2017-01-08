/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"encoding/binary"
	"io"
	"time"
)

import (
	hw "github.com/djthorpe/gopi/hw"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE TYPES

type evType uint16
type evKeyCode uint16
type evKeyAction uint32

type evEvent struct {
	Second      uint32
	Microsecond uint32
	Type        evType
	Code        evKeyCode
	Value       uint32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Event types
// See https://www.kernel.org/doc/Documentation/input/event-codes.txt
const (
	EV_SYN       evType = 0x0000 // Used as markers to separate events
	EV_KEY       evType = 0x0001 // Used to describe state changes of keyboards, buttons
	EV_REL       evType = 0x0002 // Used to describe relative axis value changes
	EV_ABS       evType = 0x0003 // Used to describe absolute axis value changes
	EV_MSC       evType = 0x0004 // Miscellaneous uses that didn't fit anywhere else
	EV_SW        evType = 0x0005 // Used to describe binary state input switches
	EV_LED       evType = 0x0011 // Used to turn LEDs on devices on and off
	EV_SND       evType = 0x0012 // Sound output, such as buzzers
	EV_REP       evType = 0x0014 // Enables autorepeat of keys in the input core
	EV_FF        evType = 0x0015 // Sends force-feedback effects to a device
	EV_PWR       evType = 0x0016 // Power management events
	EV_FF_STATUS evType = 0x0017 // Device reporting of force-feedback effects back to the host
	EV_MAX       evType = 0x001F
)

const (
	EV_CODE_X        evKeyCode = 0x0000
	EV_CODE_Y        evKeyCode = 0x0001
	EV_CODE_SCANCODE evKeyCode = 0x0004 // Keyboard scan code
	EV_CODE_SLOT     evKeyCode = 0x002F // Slot for multi touch positon
	EV_CODE_SLOT_X   evKeyCode = 0x0035 // X for multi touch position
	EV_CODE_SLOT_Y   evKeyCode = 0x0036 // Y for multi touch position
	EV_CODE_SLOT_ID  evKeyCode = 0x0039 // Unique ID for multi touch position
)

const (
	EV_VALUE_KEY_NONE   evKeyAction = 0x00000000
	EV_VALUE_KEY_UP     evKeyAction = 0x00000000
	EV_VALUE_KEY_DOWN   evKeyAction = 0x00000001
	EV_VALUE_KEY_REPEAT evKeyAction = 0x00000002
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t evType) String() string {
	switch t {
	case EV_SYN:
		return "EV_SYN"
	case EV_KEY:
		return "EV_KEY"
	case EV_REL:
		return "EV_REL"
	case EV_ABS:
		return "EV_ABS"
	case EV_MSC:
		return "EV_MSC"
	case EV_SW:
		return "EV_SW"
	case EV_LED:
		return "EV_LED"
	case EV_SND:
		return "EV_SND"
	case EV_REP:
		return "EV_REP"
	case EV_FF:
		return "EV_FF"
	case EV_PWR:
		return "EV_PWR"
	case EV_FF_STATUS:
		return "EV_FF_STATUS"
	default:
		return "[?? Unknown evType value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// WATCH

func (this *InputDriver) Watch(delta time.Duration, callback hw.InputEventCallback) error {
	if err := this.poll.Watch(delta, func(fd int, flags PollMode) {
		// Obtain device
		device, exists := this.devices[fd]
		if exists == false {
			return
		}
		// Read raw event data
		var event evEvent
		err := binary.Read(device.handle, binary.LittleEndian, &event)
		if err == io.EOF {
			return
		}
		if err != nil {
			this.log.Error("<linux.Input>Wtch Error: %v", err)
			return
		}
		// Process the event data, callback
		if emit_event := this.evDecode(&event, device); emit_event != nil {
			callback(emit_event, device)
		}
	}); err != nil {
		return this.log.Error("<linux.Input>Watch Error: %v", err)
	}

	// success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// DECODE

// This function takes a raw event in and sets the device internal state, and
// may return an abstract input event, or nil if there is no event to emit.
func (this *InputDriver) evDecode(raw_event *evEvent, device *evDevice) *hw.InputEvent {
	this.log.Debug2("<linux.InputDriver>evDecode{ type=%v code=%v value=0x%08X ts=%v,%v }", raw_event.Type, raw_event.Code, raw_event.Value, raw_event.Second, raw_event.Microsecond)
	switch raw_event.Type {
	case EV_SYN:
		return this.evDecodeSyn(raw_event, device)
	case EV_KEY:
		this.evDecodeKey(raw_event, device)
	case EV_ABS:
		return this.evDecodeAbs(raw_event, device)
	case EV_REL:
		this.evDecodeRel(raw_event, device)
	case EV_MSC:
		this.evDecodeMsc(raw_event, device)
	default:
		this.log.Warn("<linux.Input>Watch device=%v event=%v Ignoring event type", device, raw_event.Type)
	}

	// Don't return an event
	return nil
}

// Decode the EV_SYN syncronization raw event.
func (this *InputDriver) evDecodeSyn(raw_event *evEvent, device *evDevice) *hw.InputEvent {
	event := hw.InputEvent{}
	event.Timestamp = time.Duration(time.Duration(raw_event.Second)*time.Second + time.Duration(raw_event.Microsecond)*time.Microsecond)
	event.DeviceType = device.device_type
	event.Position = device.position

	// Mouse and keyboard movements
	if device.rel_position.Equals(khronos.EGLZeroPoint) == false {
		event.EventType = hw.INPUT_EVENT_RELPOSITION
		event.Relative = device.rel_position
		device.rel_position = khronos.EGLZeroPoint
		device.last_position = device.position
	} else if device.position.Equals(device.last_position) == false {
		event.EventType = hw.INPUT_EVENT_ABSPOSITION
		device.last_position = device.position
	} else if device.key_action == EV_VALUE_KEY_UP {
		event.EventType = hw.INPUT_EVENT_KEYRELEASE
		event.Keycode = hw.InputKeyCode(device.key_code)
		event.Scancode = device.scan_code
		device.key_action = EV_VALUE_KEY_NONE
	} else if device.key_action == EV_VALUE_KEY_DOWN {
		event.EventType = hw.INPUT_EVENT_KEYPRESS
		event.Keycode = hw.InputKeyCode(device.key_code)
		event.Scancode = device.scan_code
		device.key_action = EV_VALUE_KEY_NONE
	} else if device.key_action == EV_VALUE_KEY_REPEAT {
		event.EventType = hw.INPUT_EVENT_KEYREPEAT
		event.Keycode = hw.InputKeyCode(device.key_code)
		event.Scancode = device.scan_code
		device.key_action = EV_VALUE_KEY_NONE
	} else {
		return nil
	}

	// Check for multi-touch positional changes
	slot := device.slots[device.slot]
	if slot.active {
		this.log.Debug("SLOT=%v POSITION=%v", device.slot, slot.position)
	}

	return &event
}

func (this *InputDriver) evDecodeKey(raw_event *evEvent, device *evDevice) {
	device.key_code = evKeyCode(raw_event.Code)
	device.key_action = evKeyAction(raw_event.Value)

	// Set the device state from the key action. For the locks (Caps, Scroll
	// and Num) we also reflect the change with the LED and "flip" the state
	// from the current state.
	key_state := hw.INPUT_KEYSTATE_NONE

	switch hw.InputKeyCode(device.key_code) {
	case hw.INPUT_KEY_CAPSLOCK:
		// Flip CAPS LOCK state and set LED
		if device.key_action == EV_VALUE_KEY_DOWN {
			device.state ^= hw.INPUT_KEYSTATE_CAPSLOCK
			evSetLEDState(device.handle, EV_LED_CAPSL, device.state&hw.INPUT_KEYSTATE_CAPSLOCK != hw.INPUT_KEYSTATE_NONE)
		}
	case hw.INPUT_KEY_NUMLOCK:
		// Flip NUM LOCK state and set LED
		if device.key_action == EV_VALUE_KEY_DOWN {
			device.state ^= hw.INPUT_KEYSTATE_NUMLOCK
			evSetLEDState(device.handle, EV_LED_NUML, device.state&hw.INPUT_KEYSTATE_NUMLOCK != hw.INPUT_KEYSTATE_NONE)
		}
	case hw.INPUT_KEY_SCROLLLOCK:
		// Flip SCROLL LOCK state and set LED
		if device.key_action == EV_VALUE_KEY_DOWN {
			device.state ^= hw.INPUT_KEYSTATE_SCROLLLOCK
			evSetLEDState(device.handle, EV_LED_SCROLLL, device.state&hw.INPUT_KEYSTATE_SCROLLLOCK != hw.INPUT_KEYSTATE_NONE)
		}
	case hw.INPUT_KEY_LEFTSHIFT:
		key_state = hw.INPUT_KEYSTATE_LEFTSHIFT
	case hw.INPUT_KEY_RIGHTSHIFT:
		key_state = hw.INPUT_KEYSTATE_RIGHTSHIFT
	case hw.INPUT_KEY_LEFTCTRL:
		key_state = hw.INPUT_KEYSTATE_LEFTCTRL
	case hw.INPUT_KEY_RIGHTCTRL:
		key_state = hw.INPUT_KEYSTATE_RIGHTCTRL
	case hw.INPUT_KEY_LEFTALT:
		key_state = hw.INPUT_KEYSTATE_LEFTALT
	case hw.INPUT_KEY_RIGHTALT:
		key_state = hw.INPUT_KEYSTATE_RIGHTALT
	case hw.INPUT_KEY_LEFTMETA:
		key_state = hw.INPUT_KEYSTATE_LEFTMETA
	case hw.INPUT_KEY_RIGHTMETA:
		key_state = hw.INPUT_KEYSTATE_RIGHTMETA
	}

	// Set state from key action
	if key_state != hw.INPUT_KEYSTATE_NONE {
		if device.key_action == EV_VALUE_KEY_DOWN || device.key_action == EV_VALUE_KEY_REPEAT {
			device.state |= key_state
		} else if device.key_action == EV_VALUE_KEY_UP {
			device.state &= (hw.INPUT_KEYSTATE_MAX ^ key_state)
		}
	}
}

func (this *InputDriver) evDecodeAbs(raw_event *evEvent, device *evDevice) *hw.InputEvent {
	if raw_event.Code == EV_CODE_X {
		device.position.X = int(raw_event.Value)
	} else if raw_event.Code == EV_CODE_Y {
		device.position.Y = int(raw_event.Value)
	} else if raw_event.Code == EV_CODE_SLOT {
		device.slot = raw_event.Value
	} else if raw_event.Code == EV_CODE_SLOT_ID || raw_event.Code == EV_CODE_SLOT_X || raw_event.Code == EV_CODE_SLOT_Y {
		switch {
		case device.slot < uint32(0) || device.slot >= INPUT_MAX_MULTITOUCH_SLOTS:
			this.log.Warn("<linux.InputDevice> Ignoring out-of-range slot %v", device.slot)
		case raw_event.Code == EV_CODE_SLOT_ID:
			return this.evDecodeAbsTouch(raw_event, device)
		case raw_event.Code == EV_CODE_SLOT_X:
			device.slots[device.slot].position.X = int(raw_event.Value)
			device.slots[device.slot].active = true
		case raw_event.Code == EV_CODE_SLOT_Y:
			device.slots[device.slot].position.Y = int(raw_event.Value)
			device.slots[device.slot].active = true
		}
	} else {
		this.log.Debug2("%v Ignoring code %v", raw_event.Type, raw_event.Code)
	}
	return nil
}

func (this *InputDriver) evDecodeAbsTouch(raw_event *evEvent, device *evDevice) *hw.InputEvent {
	event := hw.InputEvent{}
	event.Timestamp = time.Duration(time.Duration(raw_event.Second)*time.Second + time.Duration(raw_event.Microsecond)*time.Microsecond)
	event.DeviceType = device.device_type

	// Decode the slot_id, if -1 then this is the release for a slot
	slot_id := int16(raw_event.Value)
	if slot_id >= 0 {
		device.slots[device.slot].active = true
		device.slots[device.slot].id = slot_id
		event.EventType = hw.INPUT_EVENT_TOUCHPRESS
	} else {
		device.slots[device.slot].active = false
		event.EventType = hw.INPUT_EVENT_TOUCHRELEASE
	}

	// Populate the slot and keycode
	event.Slot = uint(device.slot)
	event.Keycode = hw.INPUT_BTN_TOUCH

	// Return the event to emit
	return &event
}

func (this *InputDriver) evDecodeRel(raw_event *evEvent, device *evDevice) {
	switch raw_event.Code {
	case EV_CODE_X:
		device.position.X = device.position.X + int(raw_event.Value)
		device.rel_position.X = int(raw_event.Value)
		return
	case EV_CODE_Y:
		device.position.Y = device.position.Y + int(raw_event.Value)
		device.rel_position.Y = int(raw_event.Value)
		return
	}
	this.log.Debug2("%v Ignoring code %v", raw_event.Type, raw_event.Code)
}

func (this *InputDriver) evDecodeMsc(raw_event *evEvent, device *evDevice) {
	if raw_event.Code == EV_CODE_SCANCODE {
		device.scan_code = raw_event.Value
	} else {
		this.log.Debug2("%v Ignoring code=%v value=%v", raw_event.Type, raw_event.Code, raw_event.Value)
	}
}
