package gopi

////////////////////////////////////////////////////////////////////////////////
// TYPES

type InputEvent uint

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type KeyEvent interface {
	Event

	Code() KeyCode     // Translated keycode
	Scan() uint32      // Scancode
	Event() InputEvent // Event
	Repeat() bool      // Event is a repeat
	IRAddress() uint16 // IR device address
	IRCommand() uint16 // IR device command
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	INPUT_EVENT_NONE InputEvent = 0x0000

	// Mouse and/or keyboard key/button press events
	INPUT_EVENT_KEYPRESS   InputEvent = 0x0001
	INPUT_EVENT_KEYRELEASE InputEvent = 0x0002
	INPUT_EVENT_KEYREPEAT  InputEvent = 0x0003

	// Mouse and/or touchscreen move events
	INPUT_EVENT_ABSPOSITION InputEvent = 0x0004
	INPUT_EVENT_RELPOSITION InputEvent = 0x0005

	// Multi-touch events
	INPUT_EVENT_TOUCHPRESS    InputEvent = 0x0006
	INPUT_EVENT_TOUCHRELEASE  InputEvent = 0x0007
	INPUT_EVENT_TOUCHPOSITION InputEvent = 0x0008
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e InputEvent) String() string {
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
