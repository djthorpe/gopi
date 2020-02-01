// +build sdl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package sdl

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: sdl2
#include <SDL2/SDL.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	SDLSubsystemType uint32
	SDLErr           string
	SDLRect          C.SDL_Rect
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SDL_INIT_NONE       SDLSubsystemType = 0x00000000
	SDL_INIT_TIMER      SDLSubsystemType = 0x00000001
	SDL_INIT_AUDIO      SDLSubsystemType = 0x00000010
	SDL_INIT_VIDEO      SDLSubsystemType = 0x00000020
	SDL_INIT_CDROM      SDLSubsystemType = 0x00000100
	SDL_INIT_JOYSTICK   SDLSubsystemType = 0x00000200
	SDL_INIT_EVERYTHING SDLSubsystemType = 0x0000FFFF
	SDL_INIT_MIN                         = SDL_INIT_TIMER
	SDL_INIT_MAX                         = 0x00008000
)

////////////////////////////////////////////////////////////////////////////////
// INIT AND SHUTDOWN

func SDLInit(flags SDLSubsystemType) error {
	if C.SDL_Init(C.Uint32(flags)) != 0 {
		return SDLError()
	} else {
		return nil
	}
}

func SDLInitSubsystem(flags SDLSubsystemType) error {
	if C.SDL_InitSubSystem(C.Uint32(flags)) != 0 {
		return SDLError()
	} else {
		return nil
	}
}

func SDLQuit() {
	C.SDL_Quit()
}

func SDLQuitSubsystem(flags SDLSubsystemType) {
	C.SDL_QuitSubSystem(C.Uint32(flags))
}

func SDLWasInit(flags SDLSubsystemType) SDLSubsystemType {
	return SDLSubsystemType(C.SDL_WasInit(C.Uint32(flags)))
}

////////////////////////////////////////////////////////////////////////////////
// ERRORS

func SDLError() error {
	if str := C.SDL_GetError(); str == nil {
		return nil
	} else {
		return SDLErr(C.GoString(str))
	}
}

func SDLClearError() {
	C.SDL_ClearError()
}

func (this SDLErr) Error() string {
	return string(this)
}

////////////////////////////////////////////////////////////////////////////////
// DISPLAYS

func SDLGetNumVideoDisplays() (uint, error) {
	if num := C.SDL_GetNumVideoDisplays(); num < 0 {
		return 0, SDLError()
	} else {
		return uint(num), nil
	}
}

func SDLGetDisplayName(displayId uint) (string, error) {
	if cstr := C.SDL_GetDisplayName(C.int(displayId)); cstr == nil {
		return "", SDLError()
	} else {
		return C.GoString(cstr), nil
	}
}

func SDLGetDisplayDPI(displayId uint) (float32, float32, float32, error) {
	var vdpi, hdpi, ddpi C.float
	if err := C.SDL_GetDisplayDPI(C.int(displayId), &ddpi, &hdpi, &vdpi); err != 0 {
		return 0, 0, 0, SDLError()
	} else {
		return float32(ddpi), float32(hdpi), float32(vdpi), nil
	}
}

func SDLGetDisplayBounds(displayId uint) (SDLRect, error) {
	var rect C.SDL_Rect
	if err := C.SDL_GetDisplayBounds(C.int(displayId), &rect); err != 0 {
		return SDLRect(rect), SDLError()
	} else {
		return SDLRect(rect), nil
	}
}

func SDLGetDisplayUsableBounds(displayId uint) (SDLRect, error) {
	var rect C.SDL_Rect
	if err := C.SDL_GetDisplayUsableBounds(C.int(displayId), &rect); err != 0 {
		return SDLRect(rect), SDLError()
	} else {
		return SDLRect(rect), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r SDLRect) String() string {
	return fmt.Sprintf("<SDLRect origin={ %d,%d } size={ %d,%d }>", r.x, r.y, r.w, r.h)
}

func (s SDLSubsystemType) String() string {
	str := ""
	if s == SDL_INIT_NONE {
		return s.StringFlag()
	}
	for v := SDL_INIT_MIN; v < SDL_INIT_MAX; v <<= 1 {
		if s&v == v {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (s SDLSubsystemType) StringFlag() string {
	switch s {
	case SDL_INIT_NONE:
		return "SDL_INIT_NONE"
	case SDL_INIT_TIMER:
		return "SDL_INIT_TIMER"
	case SDL_INIT_AUDIO:
		return "SDL_INIT_AUDIO"
	case SDL_INIT_VIDEO:
		return "SDL_INIT_VIDEO"
	case SDL_INIT_CDROM:
		return "SDL_INIT_CDROM"
	case SDL_INIT_JOYSTICK:
		return "SDL_INIT_JOYSTICK"
	case SDL_INIT_EVERYTHING:
		return "SDL_INIT_EVERYTHING"
	default:
		return fmt.Sprintf("SDL_INIT_%08X", uint32(s))
	}
}
