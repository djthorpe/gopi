/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package hw // import "github.com/djthorpe/gopi/hw"

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	INPUT_KEY_NONE         InputKeyCode = 0x0000
	INPUT_KEY_ESC          InputKeyCode = 0x0001
	INPUT_KEY_1            InputKeyCode = 0x0002
	INPUT_KEY_2            InputKeyCode = 0x0003
	INPUT_KEY_3            InputKeyCode = 0x0004
	INPUT_KEY_4            InputKeyCode = 0x0005
	INPUT_KEY_5            InputKeyCode = 0x0006
	INPUT_KEY_6            InputKeyCode = 0x0007
	INPUT_KEY_7            InputKeyCode = 0x0008
	INPUT_KEY_8            InputKeyCode = 0x0009
	INPUT_KEY_9            InputKeyCode = 0x000A
	INPUT_KEY_0            InputKeyCode = 0x000B
	INPUT_KEY_MINUS        InputKeyCode = 0x000C
	INPUT_KEY_EQUAL        InputKeyCode = 0x000D
	INPUT_KEY_BACKSPACE    InputKeyCode = 0x000E
	INPUT_KEY_TAB          InputKeyCode = 0x000F
	INPUT_KEY_Q            InputKeyCode = 0x0010
	INPUT_KEY_W            InputKeyCode = 0x0011
	INPUT_KEY_E            InputKeyCode = 0x0012
	INPUT_KEY_R            InputKeyCode = 0x0013
	INPUT_KEY_T            InputKeyCode = 0x0014
	INPUT_KEY_Y            InputKeyCode = 0x0015
	INPUT_KEY_U            InputKeyCode = 0x0016
	INPUT_KEY_I            InputKeyCode = 0x0017
	INPUT_KEY_O            InputKeyCode = 0x0018
	INPUT_KEY_P            InputKeyCode = 0x0019
	INPUT_KEY_LEFTBRACE    InputKeyCode = 0x001A
	INPUT_KEY_RIGHTBRACE   InputKeyCode = 0x001B
	INPUT_KEY_ENTER        InputKeyCode = 0x001C
	INPUT_KEY_LEFTCTRL     InputKeyCode = 0x001D
	INPUT_KEY_A            InputKeyCode = 0x001E
	INPUT_KEY_S            InputKeyCode = 0x001F
	INPUT_KEY_D            InputKeyCode = 0x0020
	INPUT_KEY_F            InputKeyCode = 0x0021
	INPUT_KEY_G            InputKeyCode = 0x0022
	INPUT_KEY_H            InputKeyCode = 0x0023
	INPUT_KEY_J            InputKeyCode = 0x0024
	INPUT_KEY_K            InputKeyCode = 0x0025
	INPUT_KEY_L            InputKeyCode = 0x0026
	INPUT_KEY_SEMICOLON    InputKeyCode = 0x0027
	INPUT_KEY_APOSTROPHE   InputKeyCode = 0x0028
	INPUT_KEY_GRAVE        InputKeyCode = 0x0029
	INPUT_KEY_LEFTSHIFT    InputKeyCode = 0x002A
	INPUT_KEY_BACKSLASH    InputKeyCode = 0x002B
	INPUT_KEY_Z            InputKeyCode = 0x002C
	INPUT_KEY_X            InputKeyCode = 0x002D
	INPUT_KEY_C            InputKeyCode = 0x002E
	INPUT_KEY_V            InputKeyCode = 0x002F
	INPUT_KEY_B            InputKeyCode = 0x0030
	INPUT_KEY_N            InputKeyCode = 0x0031
	INPUT_KEY_M            InputKeyCode = 0x0032
	INPUT_KEY_COMMA        InputKeyCode = 0x0033
	INPUT_KEY_DOT          InputKeyCode = 0x0034
	INPUT_KEY_SLASH        InputKeyCode = 0x0035
	INPUT_KEY_RIGHTSHIFT   InputKeyCode = 0x0036
	INPUT_KEY_KPASTERISK   InputKeyCode = 0x0037
	INPUT_KEY_LEFTALT      InputKeyCode = 0x0038
	INPUT_KEY_SPACE        InputKeyCode = 0x0039
	INPUT_KEY_CAPSLOCK     InputKeyCode = 0x003A
	INPUT_KEY_F1           InputKeyCode = 0x003B
	INPUT_KEY_F2           InputKeyCode = 0x003C
	INPUT_KEY_F3           InputKeyCode = 0x003D
	INPUT_KEY_F4           InputKeyCode = 0x003E
	INPUT_KEY_F5           InputKeyCode = 0x003F
	INPUT_KEY_F6           InputKeyCode = 0x0040
	INPUT_KEY_F7           InputKeyCode = 0x0041
	INPUT_KEY_F8           InputKeyCode = 0x0042
	INPUT_KEY_F9           InputKeyCode = 0x0043
	INPUT_KEY_F10          InputKeyCode = 0x0044
	INPUT_KEY_NUMLOCK      InputKeyCode = 0x0045
	INPUT_KEY_SCROLLLOCK   InputKeyCode = 0x0046
	INPUT_KEY_KP7          InputKeyCode = 0x0047
	INPUT_KEY_KP8          InputKeyCode = 0x0048
	INPUT_KEY_KP9          InputKeyCode = 0x0049
	INPUT_KEY_KPMINUS      InputKeyCode = 0x004A
	INPUT_KEY_KP4          InputKeyCode = 0x004B
	INPUT_KEY_KP5          InputKeyCode = 0x004C
	INPUT_KEY_KP6          InputKeyCode = 0x004D
	INPUT_KEY_KPPLUS       InputKeyCode = 0x004E
	INPUT_KEY_KP1          InputKeyCode = 0x004F
	INPUT_KEY_KP2          InputKeyCode = 0x0050
	INPUT_KEY_KP3          InputKeyCode = 0x0051
	INPUT_KEY_KP0          InputKeyCode = 0x0052
	INPUT_KEY_KPDOT        InputKeyCode = 0x0053
	INPUT_KEY_F11          InputKeyCode = 0x0057
	INPUT_KEY_F12          InputKeyCode = 0x0058
	INPUT_KEY_KPENTER      InputKeyCode = 0x0060
	INPUT_KEY_RIGHTCTRL    InputKeyCode = 0x0061
	INPUT_KEY_KPSLASH      InputKeyCode = 0x0062
	INPUT_KEY_SYSRQ        InputKeyCode = 0x0063
	INPUT_KEY_RIGHTALT     InputKeyCode = 0x0064
	INPUT_KEY_LINEFEED     InputKeyCode = 0x0065
	INPUT_KEY_HOME         InputKeyCode = 0x0066
	INPUT_KEY_UP           InputKeyCode = 0x0067
	INPUT_KEY_PAGEUP       InputKeyCode = 0x0068
	INPUT_KEY_LEFT         InputKeyCode = 0x0069
	INPUT_KEY_RIGHT        InputKeyCode = 0x006A
	INPUT_KEY_END          InputKeyCode = 0x006B
	INPUT_KEY_DOWN         InputKeyCode = 0x006C
	INPUT_KEY_PAGEDOWN     InputKeyCode = 0x006D
	INPUT_KEY_INSERT       InputKeyCode = 0x006E
	INPUT_KEY_DELETE       InputKeyCode = 0x006F
	INPUT_KEY_MACRO        InputKeyCode = 0x0070
	INPUT_KEY_MUTE         InputKeyCode = 0x0071
	INPUT_KEY_VOLUMEDOWN   InputKeyCode = 0x0072
	INPUT_KEY_VOLUMEUP     InputKeyCode = 0x0073
	INPUT_KEY_POWER        InputKeyCode = 0x0074
	INPUT_KEY_KPEQUAL      InputKeyCode = 0x0075
	INPUT_KEY_KPPLUSMINUS  InputKeyCode = 0x0076
	INPUT_KEY_KPCOMMA      InputKeyCode = 0x0079
	INPUT_KEY_LEFTMETA     InputKeyCode = 0x007D
	INPUT_KEY_RIGHTMETA    InputKeyCode = 0x007E
	INPUT_KEY_KPLEFTPAREN  InputKeyCode = 0x00B3
	INPUT_KEY_KPRIGHTPAREN InputKeyCode = 0x00B4
	INPUT_KEY_F13          InputKeyCode = 0x00B7
	INPUT_KEY_F14          InputKeyCode = 0x00B8
	INPUT_KEY_F15          InputKeyCode = 0x00B9
	INPUT_KEY_F16          InputKeyCode = 0x00BA
	INPUT_KEY_F17          InputKeyCode = 0x00BB
	INPUT_KEY_F18          InputKeyCode = 0x00BC
	INPUT_KEY_F19          InputKeyCode = 0x00BD
	INPUT_KEY_F20          InputKeyCode = 0x00BE
	INPUT_KEY_F21          InputKeyCode = 0x00BF
	INPUT_KEY_F22          InputKeyCode = 0x00C0
	INPUT_KEY_F23          InputKeyCode = 0x00C1
	INPUT_KEY_F24          InputKeyCode = 0x00C2
	INPUT_KEY_CLOSE        InputKeyCode = 0x00CE
	INPUT_KEY_PLAY         InputKeyCode = 0x00CF
	INPUT_KEY_PRINT        InputKeyCode = 0x00D2
	INPUT_KEY_CANCEL       InputKeyCode = 0x00DF
	INPUT_KEY_MAX          InputKeyCode = 0x02FF
)

const (
	INPUT_BTN_0      InputKeyCode = 0x0100
	INPUT_BTN_1      InputKeyCode = 0x0101
	INPUT_BTN_2      InputKeyCode = 0x0102
	INPUT_BTN_3      InputKeyCode = 0x0103
	INPUT_BTN_4      InputKeyCode = 0x0104
	INPUT_BTN_5      InputKeyCode = 0x0105
	INPUT_BTN_6      InputKeyCode = 0x0106
	INPUT_BTN_7      InputKeyCode = 0x0107
	INPUT_BTN_8      InputKeyCode = 0x0108
	INPUT_BTN_9      InputKeyCode = 0x0109
	INPUT_BTN_LEFT   InputKeyCode = 0x0110
	INPUT_BTN_RIGHT  InputKeyCode = 0x0111
	INPUT_BTN_MIDDLE InputKeyCode = 0x0112
	INPUT_BTN_SIDE   InputKeyCode = 0x0113
	INPUT_BTN_EXTRA  InputKeyCode = 0x0114
	INPUT_BTN_TOUCH  InputKeyCode = 0x014A
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (k InputKeyCode) String() string {
	switch k {
	case INPUT_KEY_NONE:
		return "INPUT_KEY_NONE"
	case INPUT_KEY_ESC:
		return "INPUT_KEY_ESC"
	case INPUT_KEY_1:
		return "INPUT_KEY_1"
	case INPUT_KEY_2:
		return "INPUT_KEY_2"
	case INPUT_KEY_3:
		return "INPUT_KEY_3"
	case INPUT_KEY_4:
		return "INPUT_KEY_4"
	case INPUT_KEY_5:
		return "INPUT_KEY_5"
	case INPUT_KEY_6:
		return "INPUT_KEY_6"
	case INPUT_KEY_7:
		return "INPUT_KEY_7"
	case INPUT_KEY_8:
		return "INPUT_KEY_8"
	case INPUT_KEY_9:
		return "INPUT_KEY_9"
	case INPUT_KEY_0:
		return "INPUT_KEY_0"
	case INPUT_KEY_MINUS:
		return "INPUT_KEY_MINUS"
	case INPUT_KEY_EQUAL:
		return "INPUT_KEY_EQUAL"
	case INPUT_KEY_BACKSPACE:
		return "INPUT_KEY_BACKSPACE"
	case INPUT_KEY_TAB:
		return "INPUT_KEY_TAB"
	case INPUT_KEY_Q:
		return "INPUT_KEY_Q"
	case INPUT_KEY_W:
		return "INPUT_KEY_W"
	case INPUT_KEY_E:
		return "INPUT_KEY_E"
	case INPUT_KEY_R:
		return "INPUT_KEY_R"
	case INPUT_KEY_T:
		return "INPUT_KEY_T"
	case INPUT_KEY_Y:
		return "INPUT_KEY_Y"
	case INPUT_KEY_U:
		return "INPUT_KEY_U"
	case INPUT_KEY_I:
		return "INPUT_KEY_I"
	case INPUT_KEY_O:
		return "INPUT_KEY_O"
	case INPUT_KEY_P:
		return "INPUT_KEY_P"
	case INPUT_KEY_LEFTBRACE:
		return "INPUT_KEY_LEFTBRACE"
	case INPUT_KEY_RIGHTBRACE:
		return "INPUT_KEY_RIGHTBRACE"
	case INPUT_KEY_ENTER:
		return "INPUT_KEY_ENTER"
	case INPUT_KEY_LEFTCTRL:
		return "INPUT_KEY_LEFTCTRL"
	case INPUT_KEY_A:
		return "INPUT_KEY_A"
	case INPUT_KEY_S:
		return "INPUT_KEY_S"
	case INPUT_KEY_D:
		return "INPUT_KEY_D"
	case INPUT_KEY_F:
		return "INPUT_KEY_F"
	case INPUT_KEY_G:
		return "INPUT_KEY_G"
	case INPUT_KEY_H:
		return "INPUT_KEY_H"
	case INPUT_KEY_J:
		return "INPUT_KEY_J"
	case INPUT_KEY_K:
		return "INPUT_KEY_K"
	case INPUT_KEY_L:
		return "INPUT_KEY_L"
	case INPUT_KEY_SEMICOLON:
		return "INPUT_KEY_SEMICOLON"
	case INPUT_KEY_APOSTROPHE:
		return "INPUT_KEY_APOSTROPHE"
	case INPUT_KEY_GRAVE:
		return "INPUT_KEY_GRAVE"
	case INPUT_KEY_LEFTSHIFT:
		return "INPUT_KEY_LEFTSHIFT"
	case INPUT_KEY_BACKSLASH:
		return "INPUT_KEY_BACKSLASH"
	case INPUT_KEY_Z:
		return "INPUT_KEY_Z"
	case INPUT_KEY_X:
		return "INPUT_KEY_X"
	case INPUT_KEY_C:
		return "INPUT_KEY_C"
	case INPUT_KEY_V:
		return "INPUT_KEY_V"
	case INPUT_KEY_B:
		return "INPUT_KEY_B"
	case INPUT_KEY_N:
		return "INPUT_KEY_N"
	case INPUT_KEY_M:
		return "INPUT_KEY_M"
	case INPUT_KEY_COMMA:
		return "INPUT_KEY_COMMA"
	case INPUT_KEY_DOT:
		return "INPUT_KEY_DOT"
	case INPUT_KEY_SLASH:
		return "INPUT_KEY_SLASH"
	case INPUT_KEY_RIGHTSHIFT:
		return "INPUT_KEY_RIGHTSHIFT"
	case INPUT_KEY_KPASTERISK:
		return "INPUT_KEY_KPASTERISK"
	case INPUT_KEY_LEFTALT:
		return "INPUT_KEY_LEFTALT"
	case INPUT_KEY_SPACE:
		return "INPUT_KEY_SPACE"
	case INPUT_KEY_CAPSLOCK:
		return "INPUT_KEY_CAPSLOCK"
	case INPUT_KEY_F1:
		return "INPUT_KEY_F1"
	case INPUT_KEY_F2:
		return "INPUT_KEY_F2"
	case INPUT_KEY_F3:
		return "INPUT_KEY_F3"
	case INPUT_KEY_F4:
		return "INPUT_KEY_F4"
	case INPUT_KEY_F5:
		return "INPUT_KEY_F5"
	case INPUT_KEY_F6:
		return "INPUT_KEY_F6"
	case INPUT_KEY_F7:
		return "INPUT_KEY_F7"
	case INPUT_KEY_F8:
		return "INPUT_KEY_F8"
	case INPUT_KEY_F9:
		return "INPUT_KEY_F9"
	case INPUT_KEY_F10:
		return "INPUT_KEY_F10"
	case INPUT_KEY_NUMLOCK:
		return "INPUT_KEY_NUMLOCK"
	case INPUT_KEY_SCROLLLOCK:
		return "INPUT_KEY_SCROLLLOCK"
	case INPUT_KEY_KP7:
		return "INPUT_KEY_KP7"
	case INPUT_KEY_KP8:
		return "INPUT_KEY_KP8"
	case INPUT_KEY_KP9:
		return "INPUT_KEY_KP9"
	case INPUT_KEY_KPMINUS:
		return "INPUT_KEY_KPMINUS"
	case INPUT_KEY_KP4:
		return "INPUT_KEY_KP4"
	case INPUT_KEY_KP5:
		return "INPUT_KEY_KP5"
	case INPUT_KEY_KP6:
		return "INPUT_KEY_KP6"
	case INPUT_KEY_KPPLUS:
		return "INPUT_KEY_KPPLUS"
	case INPUT_KEY_KP1:
		return "INPUT_KEY_KP1"
	case INPUT_KEY_KP2:
		return "INPUT_KEY_KP2"
	case INPUT_KEY_KP3:
		return "INPUT_KEY_KP3"
	case INPUT_KEY_KP0:
		return "INPUT_KEY_KP0"
	case INPUT_KEY_KPDOT:
		return "INPUT_KEY_KPDOT"
	case INPUT_KEY_F11:
		return "INPUT_KEY_F11"
	case INPUT_KEY_F12:
		return "INPUT_KEY_F12"
	case INPUT_KEY_KPENTER:
		return "INPUT_KEY_KPENTER"
	case INPUT_KEY_RIGHTCTRL:
		return "INPUT_KEY_RIGHTCTRL"
	case INPUT_KEY_KPSLASH:
		return "INPUT_KEY_KPSLASH"
	case INPUT_KEY_SYSRQ:
		return "INPUT_KEY_SYSRQ"
	case INPUT_KEY_RIGHTALT:
		return "INPUT_KEY_RIGHTALT"
	case INPUT_KEY_LINEFEED:
		return "INPUT_KEY_LINEFEED"
	case INPUT_KEY_HOME:
		return "INPUT_KEY_HOME"
	case INPUT_KEY_UP:
		return "INPUT_KEY_UP"
	case INPUT_KEY_PAGEUP:
		return "INPUT_KEY_PAGEUP"
	case INPUT_KEY_LEFT:
		return "INPUT_KEY_LEFT"
	case INPUT_KEY_RIGHT:
		return "INPUT_KEY_RIGHT"
	case INPUT_KEY_END:
		return "INPUT_KEY_END"
	case INPUT_KEY_DOWN:
		return "INPUT_KEY_DOWN"
	case INPUT_KEY_PAGEDOWN:
		return "INPUT_KEY_PAGEDOWN"
	case INPUT_KEY_INSERT:
		return "INPUT_KEY_INSERT"
	case INPUT_KEY_DELETE:
		return "INPUT_KEY_DELETE"
	case INPUT_KEY_MACRO:
		return "INPUT_KEY_MACRO"
	case INPUT_KEY_MUTE:
		return "INPUT_KEY_MUTE"
	case INPUT_KEY_VOLUMEDOWN:
		return "INPUT_KEY_VOLUMEDOWN"
	case INPUT_KEY_VOLUMEUP:
		return "INPUT_KEY_VOLUMEUP"
	case INPUT_KEY_POWER:
		return "INPUT_KEY_POWER"
	case INPUT_KEY_KPEQUAL:
		return "INPUT_KEY_KPEQUAL"
	case INPUT_KEY_KPPLUSMINUS:
		return "INPUT_KEY_KPPLUSMINUS"
	case INPUT_KEY_KPCOMMA:
		return "INPUT_KEY_KPCOMMA"
	case INPUT_KEY_LEFTMETA:
		return "INPUT_KEY_LEFTMETA"
	case INPUT_KEY_RIGHTMETA:
		return "INPUT_KEY_RIGHTMETA"
	case INPUT_KEY_KPLEFTPAREN:
		return "INPUT_KEY_KPLEFTPAREN"
	case INPUT_KEY_KPRIGHTPAREN:
		return "INPUT_KEY_KPRIGHTPAREN"
	case INPUT_KEY_F13:
		return "INPUT_KEY_F13"
	case INPUT_KEY_F14:
		return "INPUT_KEY_F14"
	case INPUT_KEY_F15:
		return "INPUT_KEY_F15"
	case INPUT_KEY_F16:
		return "INPUT_KEY_F16"
	case INPUT_KEY_F17:
		return "INPUT_KEY_F17"
	case INPUT_KEY_F18:
		return "INPUT_KEY_F18"
	case INPUT_KEY_F19:
		return "INPUT_KEY_F19"
	case INPUT_KEY_F20:
		return "INPUT_KEY_F20"
	case INPUT_KEY_F21:
		return "INPUT_KEY_F21"
	case INPUT_KEY_F22:
		return "INPUT_KEY_F22"
	case INPUT_KEY_F23:
		return "INPUT_KEY_F23"
	case INPUT_KEY_F24:
		return "INPUT_KEY_F24"
	case INPUT_KEY_CLOSE:
		return "INPUT_KEY_CLOSE"
	case INPUT_KEY_PLAY:
		return "INPUT_KEY_PLAY"
	case INPUT_KEY_PRINT:
		return "INPUT_KEY_PRINT"
	case INPUT_KEY_CANCEL:
		return "INPUT_KEY_CANCEL"
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
	case INPUT_BTN_TOUCH:
		return "INPUT_BTN_TOUCH"
	default:
		return fmt.Sprintf("INPUT_KEY_0x%04X", uint16(k))
	}
}
