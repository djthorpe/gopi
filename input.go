/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// InputDeviceType is one or more of keyboard, mouse, touchscreen, etc
	InputDeviceType uint8

	// Key Code
	KeyCode uint16

	// Key Code
	KeyAction uint32

	// Key State
	KeyState uint16

	// Type of input event
	InputEventType uint16
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Types of input device
const (
	INPUT_TYPE_NONE        InputDeviceType = 0x00
	INPUT_TYPE_KEYBOARD    InputDeviceType = 0x01
	INPUT_TYPE_MOUSE       InputDeviceType = 0x02
	INPUT_TYPE_TOUCHSCREEN InputDeviceType = 0x04
	INPUT_TYPE_JOYSTICK    InputDeviceType = 0x08
	INPUT_TYPE_REMOTE      InputDeviceType = 0x10
	INPUT_TYPE_ANY         InputDeviceType = 0xFF
)

// Key actions
const (
	KEYACTION_KEY_UP KeyAction = iota
	KEYACTION_KEY_DOWN
	KEYACTION_KEY_REPEAT
	KEYACTION_NONE
)

// Key state
const (
	KEYSTATE_NONE       KeyState = 0x0000
	KEYSTATE_SCROLLLOCK KeyState = 0x0001 // Scroll Lock
	KEYSTATE_NUMLOCK    KeyState = 0x0002 // Num Lock
	KEYSTATE_CAPSLOCK   KeyState = 0x0004 // Caps Lock
	KEYSTATE_LEFTSHIFT  KeyState = 0x0010 // Left Shift
	KEYSTATE_RIGHTSHIFT KeyState = 0x0020 // Right Shift
	KEYSTATE_SHIFT      KeyState = 0x0030 // Either Shift
	KEYSTATE_LEFTALT    KeyState = 0x0040 // Left Alt
	KEYSTATE_RIGHTALT   KeyState = 0x0080 // Right Alt
	KEYSTATE_ALT        KeyState = 0x00C0 // Either Alt
	KEYSTATE_LEFTMETA   KeyState = 0x0100 // Left Meta/Command
	KEYSTATE_RIGHTMETA  KeyState = 0x0200 // Right Meta/Command
	KEYSTATE_META       KeyState = 0x0300 // Either Meta/Command
	KEYSTATE_LEFTCTRL   KeyState = 0x0400 // Left Control
	KEYSTATE_RIGHTCTRL  KeyState = 0x0800 // Right Control
	KEYSTATE_CTRL       KeyState = 0x0C00 // Either Control

	KEYSTATE_MASK KeyState = 0x0CFF // Bitmask
	KEYSTATE_MIN  KeyState = KEYSTATE_SCROLLLOCK
	KEYSTATE_MAX  KeyState = KEYSTATE_CTRL // Maximum
)

// Input events
const (
	INPUT_EVENT_NONE          InputEventType = 0x0000
	INPUT_EVENT_KEYPRESS      InputEventType = 0x0001
	INPUT_EVENT_KEYRELEASE    InputEventType = 0x0002
	INPUT_EVENT_KEYREPEAT     InputEventType = 0x0003
	INPUT_EVENT_ABSPOSITION   InputEventType = 0x0004
	INPUT_EVENT_RELPOSITION   InputEventType = 0x0005
	INPUT_EVENT_TOUCHPRESS    InputEventType = 0x0006
	INPUT_EVENT_TOUCHRELEASE  InputEventType = 0x0007
	INPUT_EVENT_TOUCHPOSITION InputEventType = 0x0008
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// InputManager manages all input devices
type InputManager interface {

	// Open Devices by name and/or type and return a list of opened devices
	OpenDevicesByNameType(name string, flags InputDeviceType, exclusive bool) ([]InputDevice, error)

	// Open an input device, and 'grab' exclusively
	OpenDevice(bus uint, exclusive bool) (InputDevice, error)

	// Close an input device
	CloseDevice(InputDevice) error

	// Implements gopi.Unit
	Unit
}

// InputDevice represents a keyboard, mouse, touchscreen, etc.
type InputDevice interface {
	// Name returns the name of the device
	Name() string

	// Id returns a unique ID for the device
	Id() uint

	// Returns the file descriptor for the device, or zero
	Fd() uintptr

	// Type returns the type of input device
	Type() InputDeviceType

	// KeyState indicates keyboard state when a modififer key is pressed or locked
	KeyState() KeyState

	// Position returns the absolute position for the device
	// (if mouse, joystick or touchscreen)
	Position() Point

	// SetPosition sets the absolute position for the device
	SetPosition(Point)

	// Matches returns true if a device has specific capabilities or name
	Matches(name string, flags InputDeviceType) bool
}

type InputEvent interface {
	// Return device which is emitting the event
	Device() InputDevice

	// Type of event
	Type() InputEventType

	// KeyCode returned when a key press event
	KeyCode() KeyCode

	// KeyState returned when a key press event
	KeyState() KeyState

	// ScanCode returned when a key press event
	ScanCode() uint32

	// Abs returns absolute input position
	Abs() Point

	// Rel returns relative input position
	Rel() Point

	// Implements gopi.Event
	Event
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t InputDeviceType) String() string {
	switch t {
	case INPUT_TYPE_NONE:
		return "INPUT_TYPE_NONE"
	case INPUT_TYPE_KEYBOARD:
		return "INPUT_TYPE_KEYBOARD"
	case INPUT_TYPE_MOUSE:
		return "INPUT_TYPE_MOUSE"
	case INPUT_TYPE_TOUCHSCREEN:
		return "INPUT_TYPE_TOUCHSCREEN"
	case INPUT_TYPE_JOYSTICK:
		return "INPUT_TYPE_JOYSTICK"
	case INPUT_TYPE_REMOTE:
		return "INPUT_TYPE_REMOTE"
	default:
		return "[?? Invalid InputDeviceType value]"
	}
}

func (k KeyAction) String() string {
	switch k {
	case KEYACTION_KEY_UP:
		return "KEYACTION_KEY_UP"
	case KEYACTION_KEY_DOWN:
		return "KEYACTION_KEY_DOWN"
	case KEYACTION_KEY_REPEAT:
		return "KEYACTION_KEY_REPEAT"
	default:
		return "[?? Invalid KeyAction value]"
	}
}

func (s KeyState) String() string {
	str := ""
	if s == KEYSTATE_NONE {
		return s.StringFlag()
	}
	for v := KEYSTATE_MIN; v <= KEYSTATE_MAX; v <<= 1 {
		if s&v == s {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (s KeyState) StringFlag() string {
	switch s {
	case KEYSTATE_NONE:
		return "KEYSTATE_NONE"
	case KEYSTATE_SCROLLLOCK:
		return "KEYSTATE_SCROLLLOCK"
	case KEYSTATE_NUMLOCK:
		return "KEYSTATE_NUMLOCK"
	case KEYSTATE_CAPSLOCK:
		return "KEYSTATE_CAPSLOCK"
	case KEYSTATE_LEFTSHIFT:
		return "KEYSTATE_LEFTSHIFT"
	case KEYSTATE_RIGHTSHIFT:
		return "KEYSTATE_RIGHTSHIFT"
	case KEYSTATE_LEFTALT:
		return "KEYSTATE_LEFTALT"
	case KEYSTATE_RIGHTALT:
		return "KEYSTATE_RIGHTALT"
	case KEYSTATE_LEFTMETA:
		return "KEYSTATE_LEFTMETA"
	case KEYSTATE_RIGHTMETA:
		return "KEYSTATE_RIGHTMETA"
	case KEYSTATE_LEFTCTRL:
		return "KEYSTATE_LEFTCTRL"
	case KEYSTATE_RIGHTCTRL:
		return "KEYSTATE_RIGHTCTRL"
	default:
		return "[?? Invalid KeyState value]"
	}
}

func (e InputEventType) String() string {
	switch e {
	case INPUT_EVENT_NONE:
		return "INPUT_EVENT_NONE"
	case INPUT_EVENT_KEYPRESS:
		return "INPUT_EVENT_KEYPRESS"
	case INPUT_EVENT_KEYRELEASE:
		return "INPUT_EVENT_KEYRELEASE"
	case INPUT_EVENT_KEYREPEAT:
		return "INPUT_EVENT_KEYREPEAT"
	case INPUT_EVENT_ABSPOSITION:
		return "INPUT_EVENT_ABSPOSITION"
	case INPUT_EVENT_RELPOSITION:
		return "INPUT_EVENT_RELPOSITION"
	case INPUT_EVENT_TOUCHPRESS:
		return "INPUT_EVENT_TOUCHPRESS"
	case INPUT_EVENT_TOUCHRELEASE:
		return "INPUT_EVENT_TOUCHRELEASE"
	case INPUT_EVENT_TOUCHPOSITION:
		return "INPUT_EVENT_TOUCHPOSITION"
	default:
		return "[?? Invalid InputEventType value]"
	}
}
