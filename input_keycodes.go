package gopi

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	KEYCODE_NONE             KeyCode = 0x0000
	KEYCODE_ESC              KeyCode = 0x0001
	KEYCODE_1                KeyCode = 0x0002
	KEYCODE_2                KeyCode = 0x0003
	KEYCODE_3                KeyCode = 0x0004
	KEYCODE_4                KeyCode = 0x0005
	KEYCODE_5                KeyCode = 0x0006
	KEYCODE_6                KeyCode = 0x0007
	KEYCODE_7                KeyCode = 0x0008
	KEYCODE_8                KeyCode = 0x0009
	KEYCODE_9                KeyCode = 0x000A
	KEYCODE_0                KeyCode = 0x000B
	KEYCODE_MINUS            KeyCode = 0x000C
	KEYCODE_EQUAL            KeyCode = 0x000D
	KEYCODE_BACKSPACE        KeyCode = 0x000E
	KEYCODE_TAB              KeyCode = 0x000F
	KEYCODE_Q                KeyCode = 0x0010
	KEYCODE_W                KeyCode = 0x0011
	KEYCODE_E                KeyCode = 0x0012
	KEYCODE_R                KeyCode = 0x0013
	KEYCODE_T                KeyCode = 0x0014
	KEYCODE_Y                KeyCode = 0x0015
	KEYCODE_U                KeyCode = 0x0016
	KEYCODE_I                KeyCode = 0x0017
	KEYCODE_O                KeyCode = 0x0018
	KEYCODE_P                KeyCode = 0x0019
	KEYCODE_LEFTBRACE        KeyCode = 0x001A
	KEYCODE_RIGHTBRACE       KeyCode = 0x001B
	KEYCODE_ENTER            KeyCode = 0x001C
	KEYCODE_LEFTCTRL         KeyCode = 0x001D
	KEYCODE_A                KeyCode = 0x001E
	KEYCODE_S                KeyCode = 0x001F
	KEYCODE_D                KeyCode = 0x0020
	KEYCODE_F                KeyCode = 0x0021
	KEYCODE_G                KeyCode = 0x0022
	KEYCODE_H                KeyCode = 0x0023
	KEYCODE_J                KeyCode = 0x0024
	KEYCODE_K                KeyCode = 0x0025
	KEYCODE_L                KeyCode = 0x0026
	KEYCODE_SEMICOLON        KeyCode = 0x0027
	KEYCODE_APOSTROPHE       KeyCode = 0x0028
	KEYCODE_GRAVE            KeyCode = 0x0029
	KEYCODE_LEFTSHIFT        KeyCode = 0x002A
	KEYCODE_BACKSLASH        KeyCode = 0x002B
	KEYCODE_Z                KeyCode = 0x002C
	KEYCODE_X                KeyCode = 0x002D
	KEYCODE_C                KeyCode = 0x002E
	KEYCODE_V                KeyCode = 0x002F
	KEYCODE_B                KeyCode = 0x0030
	KEYCODE_N                KeyCode = 0x0031
	KEYCODE_M                KeyCode = 0x0032
	KEYCODE_COMMA            KeyCode = 0x0033
	KEYCODE_DOT              KeyCode = 0x0034
	KEYCODE_SLASH            KeyCode = 0x0035
	KEYCODE_RIGHTSHIFT       KeyCode = 0x0036
	KEYCODE_KPASTERISK       KeyCode = 0x0037
	KEYCODE_LEFTALT          KeyCode = 0x0038
	KEYCODE_SPACE            KeyCode = 0x0039
	KEYCODE_CAPSLOCK         KeyCode = 0x003A
	KEYCODE_F1               KeyCode = 0x003B
	KEYCODE_F2               KeyCode = 0x003C
	KEYCODE_F3               KeyCode = 0x003D
	KEYCODE_F4               KeyCode = 0x003E
	KEYCODE_F5               KeyCode = 0x003F
	KEYCODE_F6               KeyCode = 0x0040
	KEYCODE_F7               KeyCode = 0x0041
	KEYCODE_F8               KeyCode = 0x0042
	KEYCODE_F9               KeyCode = 0x0043
	KEYCODE_F10              KeyCode = 0x0044
	KEYCODE_NUMLOCK          KeyCode = 0x0045
	KEYCODE_SCROLLLOCK       KeyCode = 0x0046
	KEYCODE_KP7              KeyCode = 0x0047
	KEYCODE_KP8              KeyCode = 0x0048
	KEYCODE_KP9              KeyCode = 0x0049
	KEYCODE_KPMINUS          KeyCode = 0x004A
	KEYCODE_KP4              KeyCode = 0x004B
	KEYCODE_KP5              KeyCode = 0x004C
	KEYCODE_KP6              KeyCode = 0x004D
	KEYCODE_KPPLUS           KeyCode = 0x004E
	KEYCODE_KP1              KeyCode = 0x004F
	KEYCODE_KP2              KeyCode = 0x0050
	KEYCODE_KP3              KeyCode = 0x0051
	KEYCODE_KP0              KeyCode = 0x0052
	KEYCODE_KPDOT            KeyCode = 0x0053
	KEYCODE_F11              KeyCode = 0x0057
	KEYCODE_F12              KeyCode = 0x0058
	KEYCODE_KPENTER          KeyCode = 0x0060
	KEYCODE_RIGHTCTRL        KeyCode = 0x0061
	KEYCODE_KPSLASH          KeyCode = 0x0062
	KEYCODE_SYSRQ            KeyCode = 0x0063
	KEYCODE_RIGHTALT         KeyCode = 0x0064
	KEYCODE_LINEFEED         KeyCode = 0x0065
	KEYCODE_HOME             KeyCode = 0x0066
	KEYCODE_UP               KeyCode = 0x0067
	KEYCODE_PAGEUP           KeyCode = 0x0068
	KEYCODE_LEFT             KeyCode = 0x0069
	KEYCODE_RIGHT            KeyCode = 0x006A
	KEYCODE_END              KeyCode = 0x006B
	KEYCODE_DOWN             KeyCode = 0x006C
	KEYCODE_PAGEDOWN         KeyCode = 0x006D
	KEYCODE_INSERT           KeyCode = 0x006E
	KEYCODE_DELETE           KeyCode = 0x006F
	KEYCODE_MACRO            KeyCode = 0x0070
	KEYCODE_MUTE             KeyCode = 0x0071
	KEYCODE_VOLUMEDOWN       KeyCode = 0x0072
	KEYCODE_VOLUMEUP         KeyCode = 0x0073
	KEYCODE_POWER            KeyCode = 0x0074
	KEYCODE_KPEQUAL          KeyCode = 0x0075
	KEYCODE_KPPLUSMINUS      KeyCode = 0x0076
	KEYCODE_KPCOMMA          KeyCode = 0x0079
	KEYCODE_LEFTMETA         KeyCode = 0x007D
	KEYCODE_RIGHTMETA        KeyCode = 0x007E
	KEYCODE_SLEEP            KeyCode = 0x008E
	KEYCODE_WAKEUP           KeyCode = 0x008F
	KEYCODE_KPLEFTPAREN      KeyCode = 0x00B3
	KEYCODE_KPRIGHTPAREN     KeyCode = 0x00B4
	KEYCODE_F13              KeyCode = 0x00B7
	KEYCODE_F14              KeyCode = 0x00B8
	KEYCODE_F15              KeyCode = 0x00B9
	KEYCODE_F16              KeyCode = 0x00BA
	KEYCODE_F17              KeyCode = 0x00BB
	KEYCODE_F18              KeyCode = 0x00BC
	KEYCODE_F19              KeyCode = 0x00BD
	KEYCODE_F20              KeyCode = 0x00BE
	KEYCODE_F21              KeyCode = 0x00BF
	KEYCODE_F22              KeyCode = 0x00C0
	KEYCODE_F23              KeyCode = 0x00C1
	KEYCODE_F24              KeyCode = 0x00C2
	KEYCODE_CLOSE            KeyCode = 0x00CE
	KEYCODE_PLAY             KeyCode = 0x00CF
	KEYCODE_PRINT            KeyCode = 0x00D2
	KEYCODE_SEARCH           KeyCode = 0x00D9
	KEYCODE_CANCEL           KeyCode = 0x00DF
	KEYCODE_BRIGHTNESS_DOWN  KeyCode = 0x00E0
	KEYCODE_BRIGHTNESS_UP    KeyCode = 0x00E1
	KEYCODE_BRIGHTNESS_CYCLE KeyCode = 0x00F3
	KEYCODE_BRIGHTNESS_AUTO  KeyCode = 0x00F4
	KEYCODE_MAX              KeyCode = 0x02FF
)

const (
	KEYCODE_BTN0      KeyCode = 0x0100
	KEYCODE_BTN1      KeyCode = 0x0101
	KEYCODE_BTN2      KeyCode = 0x0102
	KEYCODE_BTN3      KeyCode = 0x0103
	KEYCODE_BTN4      KeyCode = 0x0104
	KEYCODE_BTN5      KeyCode = 0x0105
	KEYCODE_BTN6      KeyCode = 0x0106
	KEYCODE_BTN7      KeyCode = 0x0107
	KEYCODE_BTN8      KeyCode = 0x0108
	KEYCODE_BTN9      KeyCode = 0x0109
	KEYCODE_BTNLEFT   KeyCode = 0x0110
	KEYCODE_BTNRIGHT  KeyCode = 0x0111
	KEYCODE_BTNMIDDLE KeyCode = 0x0112
	KEYCODE_BTNSIDE   KeyCode = 0x0113
	KEYCODE_BTNEXTRA  KeyCode = 0x0114
	KEYCODE_BTNTOUCH  KeyCode = 0x014A
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (k KeyCode) String() string {
	switch k {
	case KEYCODE_NONE:
		return "KEYCODE_NONE"
	case KEYCODE_ESC:
		return "KEYCODE_ESC"
	case KEYCODE_1:
		return "KEYCODE_1"
	case KEYCODE_2:
		return "KEYCODE_2"
	case KEYCODE_3:
		return "KEYCODE_3"
	case KEYCODE_4:
		return "KEYCODE_4"
	case KEYCODE_5:
		return "KEYCODE_5"
	case KEYCODE_6:
		return "KEYCODE_6"
	case KEYCODE_7:
		return "KEYCODE_7"
	case KEYCODE_8:
		return "KEYCODE_8"
	case KEYCODE_9:
		return "KEYCODE_9"
	case KEYCODE_0:
		return "KEYCODE_0"
	case KEYCODE_MINUS:
		return "KEYCODE_MINUS"
	case KEYCODE_EQUAL:
		return "KEYCODE_EQUAL"
	case KEYCODE_BACKSPACE:
		return "KEYCODE_BACKSPACE"
	case KEYCODE_TAB:
		return "KEYCODE_TAB"
	case KEYCODE_Q:
		return "KEYCODE_Q"
	case KEYCODE_W:
		return "KEYCODE_W"
	case KEYCODE_E:
		return "KEYCODE_E"
	case KEYCODE_R:
		return "KEYCODE_R"
	case KEYCODE_T:
		return "KEYCODE_T"
	case KEYCODE_Y:
		return "KEYCODE_Y"
	case KEYCODE_U:
		return "KEYCODE_U"
	case KEYCODE_I:
		return "KEYCODE_I"
	case KEYCODE_O:
		return "KEYCODE_O"
	case KEYCODE_P:
		return "KEYCODE_P"
	case KEYCODE_LEFTBRACE:
		return "KEYCODE_LEFTBRACE"
	case KEYCODE_RIGHTBRACE:
		return "KEYCODE_RIGHTBRACE"
	case KEYCODE_ENTER:
		return "KEYCODE_ENTER"
	case KEYCODE_LEFTCTRL:
		return "KEYCODE_LEFTCTRL"
	case KEYCODE_A:
		return "KEYCODE_A"
	case KEYCODE_S:
		return "KEYCODE_S"
	case KEYCODE_D:
		return "KEYCODE_D"
	case KEYCODE_F:
		return "KEYCODE_F"
	case KEYCODE_G:
		return "KEYCODE_G"
	case KEYCODE_H:
		return "KEYCODE_H"
	case KEYCODE_J:
		return "KEYCODE_J"
	case KEYCODE_K:
		return "KEYCODE_K"
	case KEYCODE_L:
		return "KEYCODE_L"
	case KEYCODE_SEMICOLON:
		return "KEYCODE_SEMICOLON"
	case KEYCODE_APOSTROPHE:
		return "KEYCODE_APOSTROPHE"
	case KEYCODE_GRAVE:
		return "KEYCODE_GRAVE"
	case KEYCODE_LEFTSHIFT:
		return "KEYCODE_LEFTSHIFT"
	case KEYCODE_BACKSLASH:
		return "KEYCODE_BACKSLASH"
	case KEYCODE_Z:
		return "KEYCODE_Z"
	case KEYCODE_X:
		return "KEYCODE_X"
	case KEYCODE_C:
		return "KEYCODE_C"
	case KEYCODE_V:
		return "KEYCODE_V"
	case KEYCODE_B:
		return "KEYCODE_B"
	case KEYCODE_N:
		return "KEYCODE_N"
	case KEYCODE_M:
		return "KEYCODE_M"
	case KEYCODE_COMMA:
		return "KEYCODE_COMMA"
	case KEYCODE_DOT:
		return "KEYCODE_DOT"
	case KEYCODE_SLASH:
		return "KEYCODE_SLASH"
	case KEYCODE_RIGHTSHIFT:
		return "KEYCODE_RIGHTSHIFT"
	case KEYCODE_KPASTERISK:
		return "KEYCODE_KPASTERISK"
	case KEYCODE_LEFTALT:
		return "KEYCODE_LEFTALT"
	case KEYCODE_SPACE:
		return "KEYCODE_SPACE"
	case KEYCODE_CAPSLOCK:
		return "KEYCODE_CAPSLOCK"
	case KEYCODE_F1:
		return "KEYCODE_F1"
	case KEYCODE_F2:
		return "KEYCODE_F2"
	case KEYCODE_F3:
		return "KEYCODE_F3"
	case KEYCODE_F4:
		return "KEYCODE_F4"
	case KEYCODE_F5:
		return "KEYCODE_F5"
	case KEYCODE_F6:
		return "KEYCODE_F6"
	case KEYCODE_F7:
		return "KEYCODE_F7"
	case KEYCODE_F8:
		return "KEYCODE_F8"
	case KEYCODE_F9:
		return "KEYCODE_F9"
	case KEYCODE_F10:
		return "KEYCODE_F10"
	case KEYCODE_NUMLOCK:
		return "KEYCODE_NUMLOCK"
	case KEYCODE_SCROLLLOCK:
		return "KEYCODE_SCROLLLOCK"
	case KEYCODE_KP7:
		return "KEYCODE_KP7"
	case KEYCODE_KP8:
		return "KEYCODE_KP8"
	case KEYCODE_KP9:
		return "KEYCODE_KP9"
	case KEYCODE_KPMINUS:
		return "KEYCODE_KPMINUS"
	case KEYCODE_KP4:
		return "KEYCODE_KP4"
	case KEYCODE_KP5:
		return "KEYCODE_KP5"
	case KEYCODE_KP6:
		return "KEYCODE_KP6"
	case KEYCODE_KPPLUS:
		return "KEYCODE_KPPLUS"
	case KEYCODE_KP1:
		return "KEYCODE_KP1"
	case KEYCODE_KP2:
		return "KEYCODE_KP2"
	case KEYCODE_KP3:
		return "KEYCODE_KP3"
	case KEYCODE_KP0:
		return "KEYCODE_KP0"
	case KEYCODE_KPDOT:
		return "KEYCODE_KPDOT"
	case KEYCODE_F11:
		return "KEYCODE_F11"
	case KEYCODE_F12:
		return "KEYCODE_F12"
	case KEYCODE_KPENTER:
		return "KEYCODE_KPENTER"
	case KEYCODE_RIGHTCTRL:
		return "KEYCODE_RIGHTCTRL"
	case KEYCODE_KPSLASH:
		return "KEYCODE_KPSLASH"
	case KEYCODE_SYSRQ:
		return "KEYCODE_SYSRQ"
	case KEYCODE_RIGHTALT:
		return "KEYCODE_RIGHTALT"
	case KEYCODE_LINEFEED:
		return "KEYCODE_LINEFEED"
	case KEYCODE_HOME:
		return "KEYCODE_HOME"
	case KEYCODE_UP:
		return "KEYCODE_UP"
	case KEYCODE_PAGEUP:
		return "KEYCODE_PAGEUP"
	case KEYCODE_LEFT:
		return "KEYCODE_LEFT"
	case KEYCODE_RIGHT:
		return "KEYCODE_RIGHT"
	case KEYCODE_END:
		return "KEYCODE_END"
	case KEYCODE_DOWN:
		return "KEYCODE_DOWN"
	case KEYCODE_PAGEDOWN:
		return "KEYCODE_PAGEDOWN"
	case KEYCODE_INSERT:
		return "KEYCODE_INSERT"
	case KEYCODE_DELETE:
		return "KEYCODE_DELETE"
	case KEYCODE_MACRO:
		return "KEYCODE_MACRO"
	case KEYCODE_MUTE:
		return "KEYCODE_MUTE"
	case KEYCODE_VOLUMEDOWN:
		return "KEYCODE_VOLUMEDOWN"
	case KEYCODE_VOLUMEUP:
		return "KEYCODE_VOLUMEUP"
	case KEYCODE_POWER:
		return "KEYCODE_POWER"
	case KEYCODE_KPEQUAL:
		return "KEYCODE_KPEQUAL"
	case KEYCODE_KPPLUSMINUS:
		return "KEYCODE_KPPLUSMINUS"
	case KEYCODE_KPCOMMA:
		return "KEYCODE_KPCOMMA"
	case KEYCODE_LEFTMETA:
		return "KEYCODE_LEFTMETA"
	case KEYCODE_RIGHTMETA:
		return "KEYCODE_RIGHTMETA"
	case KEYCODE_KPLEFTPAREN:
		return "KEYCODE_KPLEFTPAREN"
	case KEYCODE_KPRIGHTPAREN:
		return "KEYCODE_KPRIGHTPAREN"
	case KEYCODE_F13:
		return "KEYCODE_F13"
	case KEYCODE_F14:
		return "KEYCODE_F14"
	case KEYCODE_F15:
		return "KEYCODE_F15"
	case KEYCODE_F16:
		return "KEYCODE_F16"
	case KEYCODE_F17:
		return "KEYCODE_F17"
	case KEYCODE_F18:
		return "KEYCODE_F18"
	case KEYCODE_F19:
		return "KEYCODE_F19"
	case KEYCODE_F20:
		return "KEYCODE_F20"
	case KEYCODE_F21:
		return "KEYCODE_F21"
	case KEYCODE_F22:
		return "KEYCODE_F22"
	case KEYCODE_F23:
		return "KEYCODE_F23"
	case KEYCODE_F24:
		return "KEYCODE_F24"
	case KEYCODE_CLOSE:
		return "KEYCODE_CLOSE"
	case KEYCODE_PLAY:
		return "KEYCODE_PLAY"
	case KEYCODE_PRINT:
		return "KEYCODE_PRINT"
	case KEYCODE_CANCEL:
		return "KEYCODE_CANCEL"
	case KEYCODE_BTN0:
		return "KEYCODE_BTN0"
	case KEYCODE_BTN1:
		return "KEYCODE_BTN1"
	case KEYCODE_BTN2:
		return "KEYCODE_BTN2"
	case KEYCODE_BTN3:
		return "KEYCODE_BTN3"
	case KEYCODE_BTN4:
		return "KEYCODE_BTN4"
	case KEYCODE_BTN5:
		return "KEYCODE_BTN5"
	case KEYCODE_BTN6:
		return "KEYCODE_BTN6"
	case KEYCODE_BTN7:
		return "KEYCODE_BTN7"
	case KEYCODE_BTN8:
		return "KEYCODE_BTN8"
	case KEYCODE_BTN9:
		return "KEYCODE_BTN9"
	case KEYCODE_BTNLEFT:
		return "KEYCODE_BTNLEFT"
	case KEYCODE_BTNRIGHT:
		return "KEYCODE_BTNRIGHT"
	case KEYCODE_BTNMIDDLE:
		return "KEYCODE_BTNMIDDLE"
	case KEYCODE_BTNSIDE:
		return "KEYCODE_BTNSIDE"
	case KEYCODE_BTNEXTRA:
		return "KEYCODE_BTNEXTRA"
	case KEYCODE_BTNTOUCH:
		return "KEYCODE_BTNTOUCH"
	default:
		return fmt.Sprintf("KEYCODE_0x%04X", uint16(k))
	}
}
