/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// InputDeviceType is one or more of keyboard, mouse, touchscreen, etc
	InputDeviceType uint8

	// Key Code
	KeyCode uint16

	// Key State
	KeyState uint16
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

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// InputManager manages all input devices
type InputManager interface {

	// Open an input device
	OpenDevice(bus uint, exclusive bool) (InputDevice, error)

	// Close an input device
	CloseDevice(InputDevice) error

	// Implements gopi.Unit
	Unit
}

// InputDevice represents a keyboard, mouse, touchscreen, etc.
type InputDevice interface {
	Name() string
	Type() InputDeviceType
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
