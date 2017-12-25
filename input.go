/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"context"
	"time"
)

// Input driver allows you to open and close input devices
// and subscribe to events emitted by devices
type Input interface {
	Driver
	Publisher

	// Open Devices by name, type and bus
	OpenDevicesByName(name string, flags InputDeviceType, bus InputDeviceBus) ([]InputDevice, error)

	// Close Device
	CloseDevice(device InputDevice) error

	// Return a list of open devices
	GetOpenDevices() []InputDevice

	// Watch for events with context
	Watch(ctx context.Context) error
}

type InputDevice interface {
	Driver
	Publisher

	// Name of the input device
	Name() string

	// Type of device
	Type() InputDeviceType

	// Bus interface
	Bus() InputDeviceBus

	// Position of cursor (for mouse, joystick and touchscreen devices)
	Position() Point

	// Set absolute current cursor position
	SetPosition(Point)

	// Get key states (caps lock, shift, scroll lock, num lock, etc)
	KeyState() KeyState

	// Set key state (or states) to on or off. Will return error
	// for key states which are not modifiable
	SetKeyState(flags KeyState, state bool) error

	// Returns true if device matches conditions
	Matches(string, InputDeviceType, InputDeviceBus) bool
}

type InputEvent interface {
	Event

	// Timestamp of event
	Timestamp() time.Duration

	// Type of device which has created the event
	DeviceType() InputDeviceType

	// Event type
	EventType() InputEventType

	// Key or mouse button press or release
	Keycode() KeyCode

	// Key scancode
	Scancode() uint32

	// Absolute cursor position
	Position() Point

	// Relative change in position
	Relative() Point

	// Multi-touch slot identifier
	Slot() uint
}

// Device type (keyboard, mouse, touchscreen, etc)
type InputDeviceType uint8

// Event type (button press, button release, etc)
type InputEventType uint16

// Key Code
type KeyCode uint16

// Key State
type KeyState uint16

// Bus type (USB, Bluetooth, etc)
type InputDeviceBus uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Types of input device
const (
	INPUT_TYPE_NONE        InputDeviceType = 0x00
	INPUT_TYPE_KEYBOARD    InputDeviceType = 0x01
	INPUT_TYPE_MOUSE       InputDeviceType = 0x02
	INPUT_TYPE_TOUCHSCREEN InputDeviceType = 0x04
	INPUT_TYPE_JOYSTICK    InputDeviceType = 0x08
	INPUT_TYPE_ANY         InputDeviceType = 0xFF
)

// Types of input connection
const (
	INPUT_BUS_NONE      InputDeviceBus = 0x0000
	INPUT_BUS_PCI       InputDeviceBus = 0x0001
	INPUT_BUS_ISAPNP    InputDeviceBus = 0x0002
	INPUT_BUS_USB       InputDeviceBus = 0x0003
	INPUT_BUS_HIL       InputDeviceBus = 0x0004
	INPUT_BUS_BLUETOOTH InputDeviceBus = 0x0005
	INPUT_BUS_VIRTUAL   InputDeviceBus = 0x0006
	INPUT_BUS_ISA       InputDeviceBus = 0x0010
	INPUT_BUS_I8042     InputDeviceBus = 0x0011
	INPUT_BUS_XTKBD     InputDeviceBus = 0x0012
	INPUT_BUS_RS232     InputDeviceBus = 0x0013
	INPUT_BUS_GAMEPORT  InputDeviceBus = 0x0014
	INPUT_BUS_PARPORT   InputDeviceBus = 0x0015
	INPUT_BUS_AMIGA     InputDeviceBus = 0x0016
	INPUT_BUS_ADB       InputDeviceBus = 0x0017
	INPUT_BUS_I2C       InputDeviceBus = 0x0018
	INPUT_BUS_HOST      InputDeviceBus = 0x0019
	INPUT_BUS_GSC       InputDeviceBus = 0x001A
	INPUT_BUS_ATARI     InputDeviceBus = 0x001B
	INPUT_BUS_SPI       InputDeviceBus = 0x001C
	INPUT_BUS_ANY       InputDeviceBus = 0xFFFF
)

// Input events
const (
	INPUT_EVENT_NONE InputEventType = 0x0000

	// Mouse and/or keyboard key/button press events
	INPUT_EVENT_KEYPRESS   InputEventType = 0x0001
	INPUT_EVENT_KEYRELEASE InputEventType = 0x0002
	INPUT_EVENT_KEYREPEAT  InputEventType = 0x0003

	// Mouse and/or touchscreen move events
	INPUT_EVENT_ABSPOSITION InputEventType = 0x0004
	INPUT_EVENT_RELPOSITION InputEventType = 0x0005

	// Multi-touch events
	INPUT_EVENT_TOUCHPRESS    InputEventType = 0x0006
	INPUT_EVENT_TOUCHRELEASE  InputEventType = 0x0007
	INPUT_EVENT_TOUCHPOSITION InputEventType = 0x0008
)

// Input key state
const (
	KEYSTATE_NONE       KeyState = 0x0000
	KEYSTATE_MIN        KeyState = KEYSTATE_SCROLLLOCK
	KEYSTATE_SCROLLLOCK KeyState = 0x0001        // Scroll Lock
	KEYSTATE_NUMLOCK    KeyState = 0x0002        // Num Lock
	KEYSTATE_CAPSLOCK   KeyState = 0x0004        // Caps Lock
	KEYSTATE_LEFTSHIFT  KeyState = 0x0010        // Left Shift
	KEYSTATE_RIGHTSHIFT KeyState = 0x0020        // Right Shift
	KEYSTATE_SHIFT      KeyState = 0x0030        // Either Shift
	KEYSTATE_LEFTALT    KeyState = 0x0040        // Left Alt
	KEYSTATE_RIGHTALT   KeyState = 0x0080        // Right Alt
	KEYSTATE_ALT        KeyState = 0x00C0        // Either Alt
	KEYSTATE_LEFTMETA   KeyState = 0x0100        // Left Meta/Command
	KEYSTATE_RIGHTMETA  KeyState = 0x0200        // Right Meta/Command
	KEYSTATE_META       KeyState = 0x0300        // Either Meta/Command
	KEYSTATE_LEFTCTRL   KeyState = 0x0400        // Left Control
	KEYSTATE_RIGHTCTRL  KeyState = 0x0800        // Right Control
	KEYSTATE_CTRL       KeyState = 0x0C00        // Either Control
	KEYSTATE_MAX        KeyState = KEYSTATE_CTRL // Maximum
	KEYSTATE_MASK       KeyState = 0x0CFF        // Bitmask
)
