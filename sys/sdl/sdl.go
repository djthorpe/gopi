// +build sdl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package sdl

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
	C.SDL_QuitSubSystem(flags)
}

func SDLWasInit(flags SDLSubsystemType) SDLSubsystemType {
	return SDLSubsystemType(C.SDL_WasInit(flags))
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
