/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	INPUT_KEY_ESC       evKeyCode = 0x0001
	INPUT_KEY_1         evKeyCode = 0x0002
	INPUT_KEY_2         evKeyCode = 0x0003
	INPUT_KEY_3         evKeyCode = 0x0004
	INPUT_KEY_4         evKeyCode = 0x0005
	INPUT_KEY_5         evKeyCode = 0x0006
	INPUT_KEY_6         evKeyCode = 0x0007
	INPUT_KEY_7         evKeyCode = 0x0008
	INPUT_KEY_8         evKeyCode = 0x0009
	INPUT_KEY_9         evKeyCode = 0x000A
	INPUT_KEY_0         evKeyCode = 0x000B
	INPUT_KEY_MINUS     evKeyCode = 0x000C
	INPUT_KEY_EQUAL     evKeyCode = 0x000D
	INPUT_KEY_BACKSPACE evKeyCode = 0x000E
	INPUT_KEY_TAB       evKeyCode = 0x000F
)

const (
	INPUT_BTN_MISC   evKeyCode = 0x0100
	INPUT_BTN_0      evKeyCode = 0x0100
	INPUT_BTN_1      evKeyCode = 0x0101
	INPUT_BTN_2      evKeyCode = 0x0102
	INPUT_BTN_3      evKeyCode = 0x0103
	INPUT_BTN_4      evKeyCode = 0x0104
	INPUT_BTN_5      evKeyCode = 0x0105
	INPUT_BTN_6      evKeyCode = 0x0106
	INPUT_BTN_7      evKeyCode = 0x0107
	INPUT_BTN_8      evKeyCode = 0x0108
	INPUT_BTN_9      evKeyCode = 0x0109
	INPUT_BTN_MOUSE  evKeyCode = 0x0110
	INPUT_BTN_LEFT   evKeyCode = 0x0110
	INPUT_BTN_RIGHT  evKeyCode = 0x0111
	INPUT_BTN_MIDDLE evKeyCode = 0x0112
	INPUT_BTN_SIDE   evKeyCode = 0x0113
	INPUT_BTN_EXTRA  evKeyCode = 0x0114
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (k evKeyCode) String() string {
	switch k {
	case INPUT_KEY_ESC:
		return "INPUT_KEY_ESC"
	case INPUT_BTN_0:
		return "INPUT_BTN_0"
	case INPUT_BTN_1:
		return "INPUT_BTN_1"
	case INPUT_BTN_2:
		return "INPUT_BTN_2"
	case INPUT_BTN_3:
		return "INPUT_BTN_3"
	case INPUT_BTN_4:
		return "INPUT_BTN_4"
	case INPUT_BTN_5:
		return "INPUT_BTN_5"
	case INPUT_BTN_6:
		return "INPUT_BTN_6"
	case INPUT_BTN_7:
		return "INPUT_BTN_7"
	case INPUT_BTN_8:
		return "INPUT_BTN_8"
	case INPUT_BTN_9:
		return "INPUT_BTN_9"
	case INPUT_BTN_LEFT:
		return "INPUT_BTN_LEFT"
	case INPUT_BTN_RIGHT:
		return "INPUT_BTN_RIGHT"
	case INPUT_BTN_MIDDLE:
		return "INPUT_BTN_MIDDLE"
	case INPUT_BTN_SIDE:
		return "INPUT_BTN_SIDE"
	case INPUT_BTN_EXTRA:
		return "INPUT_BTN_EXTRA"
	default:
		return fmt.Sprintf("INPUT_KEY_0x%04X", uint16(k))
	}
}
