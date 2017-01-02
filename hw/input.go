/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package hw // import "github.com/djthorpe/gopi/hw"

import (
	"time"
	"fmt"
)

import (
	gopi "github.com/djthorpe/gopi"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type InputDriver interface {
	// Enforces general driver
	gopi.Driver

	// Open Devices by name, type and bus
	OpenDevicesByName(name string, flags InputDeviceType, bus InputDeviceBus) ([]InputDevice, error)

	// Close Device
	CloseDevice(device InputDevice) error

	// Return a list of open devices
	GetOpenDevices() []InputDevice

	// Watch for events for an amount of time
	Watch(delta time.Duration,callback InputEventCallback) error
}

type InputDevice interface {
	// Enforces general driver
	gopi.Driver

	// Get the name of the input device
	GetName() string

	// Get the type of device
	GetType() InputDeviceType

	// Get the bus interface
	GetBus() InputDeviceBus

	// Get current cursor position (for mouse, joystick and touchscreen devices)
	GetPosition() khronos.EGLPoint

	// Set current cursor position
	SetPosition(khronos.EGLPoint)

	// Get key states (caps lock, shift, scroll lock, num lock, etc)
	GetKeyState() InputKeyState

	// Set key state (or states) to on or off. Will return error
	// for key states which are not modifiable
	SetKeyState(flags InputKeyState,state bool) error

	// Returns true if device matches conditions
	Matches(alias string, device_type InputDeviceType, device_bus InputDeviceBus) bool
}

type InputEvent struct {
	// Timestamp of event
	Timestamp time.Duration

	// Type of device which has created the event
	DeviceType InputDeviceType

	// Event type
	EventType InputEventType

	// Key or mouse button press or release
	Keycode InputKeyCode

	// Key scancode
	Scancode uint32

	// Absolute cursor position
	Position khronos.EGLPoint

	// Relative change in position
	Relative khronos.EGLPoint

	// Multi-touch slot identifier
	Slot uint
}

// Device type (keyboard, mouse, touchscreen, etc)
type InputDeviceType uint8

// Event type (button press, button release, etc)
type InputEventType uint16

// Key Code
type InputKeyCode uint16

// Key State
type InputKeyState uint16

// Bus type (USB, Bluetooth, etc)
type InputDeviceBus uint16

// Callback function
type InputEventCallback func(event *InputEvent, device InputDevice)

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
	INPUT_EVENT_NONE        InputEventType = 0x0000

	// Mouse and/or keyboard key/button press events
	INPUT_EVENT_KEYPRESS    InputEventType = 0x0001
	INPUT_EVENT_KEYRELEASE  InputEventType = 0x0002
	INPUT_EVENT_KEYREPEAT   InputEventType = 0x0003

	// Mouse and/or touchscreen move events
	INPUT_EVENT_ABSPOSITION InputEventType = 0x0004
	INPUT_EVENT_RELPOSITION InputEventType = 0x0005

	// Multi-touch events
	INPUT_EVENT_TOUCHPRESS  InputEventType = 0x0006
	INPUT_EVENT_TOUCHRELEASE InputEventType = 0x0007
	INPUT_EVENT_TOUCHPOSITION InputEventType = 0x0008
)

// Input key state
const (
	INPUT_KEYSTATE_NONE     InputKeyState = 0x0000
	INPUT_KEYSTATE_CAPS     InputKeyState = 0x0002 // Caps Lock
	INPUT_KEYSTATE_SCROLL   InputKeyState = 0x0004 // Scroll Lock
	INPUT_KEYSTATE_SHIFT    InputKeyState = 0x0008 // Shift
	INPUT_KEYSTATE_ALT      InputKeyState = 0x0010 // Alt
	INPUT_KEYSTATE_CMD      InputKeyState = 0x0020 // Command
	INPUT_KEYSTATE_NUM      InputKeyState = 0x0040 // Num Lock
	INPUT_KEYSTATE_MAX      InputKeyState = INPUT_KEYSTATE_NUM
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY FUNCTIONS


func (e InputEvent) String() string {
	switch(e.EventType) {
	case INPUT_EVENT_KEYPRESS, INPUT_EVENT_KEYRELEASE:
		return fmt.Sprintf("<linux.InputEvent>{ type=%v device=%v keycode=0x%04X scancode=0x%08X ts=%v }",e.EventType,e.DeviceType,uint16(e.Keycode),e.Scancode,e.Timestamp)
	case INPUT_EVENT_KEYREPEAT:
		return fmt.Sprintf("<linux.InputEvent>{ type=%v device=%v keycode=0x%04X ts=%v }",e.EventType,e.DeviceType,uint16(e.Keycode),e.Timestamp)
	case INPUT_EVENT_ABSPOSITION:
		return fmt.Sprintf("<linux.InputEvent>{ type=%v device=%v position=%v ts=%v }",e.EventType,e.DeviceType,e.Position,e.Timestamp)
	case INPUT_EVENT_RELPOSITION:
		return fmt.Sprintf("<linux.InputEvent>{ type=%v device=%v position=%v relative=%v ts=%v }",e.EventType,e.DeviceType,e.Position,e.Relative,e.Timestamp)
	case INPUT_EVENT_TOUCHPRESS, INPUT_EVENT_TOUCHRELEASE:
		return fmt.Sprintf("<linux.InputEvent>{ type=%v device=%v slot=%v keycode=0x%04X ts=%v }",e.EventType,e.DeviceType,e.Slot,uint16(e.Keycode),e.Timestamp)
	case INPUT_EVENT_TOUCHPOSITION:
		return fmt.Sprintf("<linux.InputEvent>{ type=%v device=%v slot=%v position=%v ts=%v }",e.EventType,e.DeviceType,e.Slot,e.Position,e.Timestamp)
	default:
		return fmt.Sprintf("<linux.InputEvent>{ type=%v device=%v keycode=0x%04X position=%v relative=%v ts=%v }",e.EventType,e.DeviceType,uint16(e.Keycode),e.Position,e.Relative,e.Timestamp)
	}
}

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
	default:
		return "[?? Invalid InputDeviceType value]"
	}
}

func (b InputDeviceBus) String() string {
	switch b {
	case INPUT_BUS_NONE:
		return "INPUT_BUS_NONE"
	case INPUT_BUS_PCI:
		return "INPUT_BUS_PCI"
	case INPUT_BUS_ISAPNP:
		return "INPUT_BUS_ISAPNP"
	case INPUT_BUS_USB:
		return "INPUT_BUS_USB"
	case INPUT_BUS_HIL:
		return "INPUT_BUS_HIL"
	case INPUT_BUS_BLUETOOTH:
		return "INPUT_BUS_BLUETOOTH"
	case INPUT_BUS_VIRTUAL:
		return "INPUT_BUS_VIRTUAL"
	case INPUT_BUS_ISA:
		return "INPUT_BUS_ISA"
	case INPUT_BUS_I8042:
		return "INPUT_BUS_I8042"
	case INPUT_BUS_XTKBD:
		return "INPUT_BUS_XTKBD"
	case INPUT_BUS_RS232:
		return "INPUT_BUS_RS232"
	case INPUT_BUS_GAMEPORT:
		return "INPUT_BUS_GAMEPORT"
	case INPUT_BUS_PARPORT:
		return "INPUT_BUS_PARPORT"
	case INPUT_BUS_AMIGA:
		return "INPUT_BUS_AMIGA"
	case INPUT_BUS_ADB:
		return "INPUT_BUS_ADB"
	case INPUT_BUS_I2C:
		return "INPUT_BUS_I2C"
	case INPUT_BUS_HOST:
		return "INPUT_BUS_HOST"
	case INPUT_BUS_GSC:
		return "INPUT_BUS_GSC"
	case INPUT_BUS_ATARI:
		return "INPUT_BUS_ATARI"
	case INPUT_BUS_SPI:
		return "INPUT_BUS_SPI"
	default:
		return "[?? Invalid InputDeviceBus value]"
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




