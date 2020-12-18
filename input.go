package gopi

import "strings"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type InputType uint
type InputDeviceType uint16

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// InputManager provides information on registered devices
type InputManager interface {
	Devices() []InputDevice

	RegisterDevice(InputDevice) error
}

// InputDevice provides information about an input device
type InputDevice interface {
	Name() string
	Type() InputDeviceType
}

// InputEvent is a key press, mouse move, etc.
type InputEvent interface {
	Event
	Key() KeyCode                      // Translated keycode
	Type() InputType                   // Event type (key press, repeat, etc)
	Device() (InputDeviceType, uint32) // Device information
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	INPUT_DEVICE_NONE        InputDeviceType = 0x0000
	INPUT_DEVICE_KEYBOARD    InputDeviceType = 0x0001
	INPUT_DEVICE_MOUSE       InputDeviceType = 0x0002
	INPUT_DEVICE_TOUCHSCREEN InputDeviceType = 0x0004
	INPUT_DEVICE_JOYSTICK    InputDeviceType = 0x0008
	INPUT_DEVICE_REMOTE      InputDeviceType = 0x0010 // IR Remote
	INPUT_DEVICE_SONY_12     InputDeviceType = 0x0020 // 12-bit Sony IR codes
	INPUT_DEVICE_SONY_15     InputDeviceType = 0x0040 // 15-bit Sony IR codes
	INPUT_DEVICE_SONY_20     InputDeviceType = 0x0080 // 20-bit Sony IR codes
	INPUT_DEVICE_RC5_14      InputDeviceType = 0x0100 // 14-bit RC5 IR codes
	INPUT_DEVICE_NEC_32      InputDeviceType = 0x0200 // 32-bit NEC IR codes
	INPUT_DEVICE_APPLETV_32  InputDeviceType = 0x0400 // 32-bit Apple TV IR codes
	INPUT_DEVICE_ANY         InputDeviceType = 0xFFFF
	INPUT_DEVICE_MIN                         = INPUT_DEVICE_KEYBOARD
	INPUT_DEVICE_MAX                         = INPUT_DEVICE_APPLETV_32
)

const (
	INPUT_EVENT_NONE InputType = 0x0000

	// Mouse and/or keyboard key/button press events
	INPUT_EVENT_KEYPRESS   InputType = 0x0001
	INPUT_EVENT_KEYRELEASE InputType = 0x0002
	INPUT_EVENT_KEYREPEAT  InputType = 0x0003

	// Mouse and/or touchscreen move events
	INPUT_EVENT_ABSPOSITION InputType = 0x0004
	INPUT_EVENT_RELPOSITION InputType = 0x0005

	// Multi-touch events
	INPUT_EVENT_TOUCHPRESS    InputType = 0x0006
	INPUT_EVENT_TOUCHRELEASE  InputType = 0x0007
	INPUT_EVENT_TOUCHPOSITION InputType = 0x0008
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e InputType) String() string {
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
		return "[?? Invalid InputEvent value]"
	}
}

func (d InputDeviceType) FlagString() string {
	switch d {
	case INPUT_DEVICE_NONE:
		return "INPUT_DEVICE_NONE"
	case INPUT_DEVICE_KEYBOARD:
		return "INPUT_DEVICE_KEYBOARD"
	case INPUT_DEVICE_MOUSE:
		return "INPUT_DEVICE_MOUSE"
	case INPUT_DEVICE_TOUCHSCREEN:
		return "INPUT_DEVICE_TOUCHSCREEN"
	case INPUT_DEVICE_JOYSTICK:
		return "INPUT_DEVICE_JOYSTICK"
	case INPUT_DEVICE_REMOTE:
		return "INPUT_DEVICE_REMOTE"
	case INPUT_DEVICE_SONY_12:
		return "INPUT_DEVICE_SONY_12"
	case INPUT_DEVICE_SONY_15:
		return "INPUT_DEVICE_SONY_15"
	case INPUT_DEVICE_SONY_20:
		return "INPUT_DEVICE_SONY_20"
	case INPUT_DEVICE_RC5_14:
		return "INPUT_DEVICE_RC5_14"
	case INPUT_DEVICE_NEC_32:
		return "INPUT_DEVICE_NEC_32"
	case INPUT_DEVICE_APPLETV_32:
		return "INPUT_DEVICE_APPLETV_32"
	case INPUT_DEVICE_ANY:
		return "INPUT_DEVICE_ANY"
	default:
		return "[?? Invalid InputDeviceType value]"
	}
}

func (d InputDeviceType) String() string {
	if d == INPUT_DEVICE_NONE || d == INPUT_DEVICE_ANY {
		return d.FlagString()
	}
	str := ""
	for v := INPUT_DEVICE_MIN; v <= INPUT_DEVICE_MAX; v <<= 1 {
		if v&d == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}
