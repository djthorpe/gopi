// +build sdl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package sdl_test

import (
	"testing"

	// Frameworks
	sdl "github.com/djthorpe/gopi/v2/sys/sdl"
)

func Test_SDL_000(t *testing.T) {
	if err := sdl.SDLInit(sdl.SDL_INIT_EVERYTHING); err != nil {
		t.Error(err)
	} else {
		sdl.SDLQuit()
	}
}

func Test_SDL_001(t *testing.T) {
	err := sdl.SDLInit(sdl.SDL_INIT_EVERYTHING)
	if err != nil {
		t.Error(err)
	}
	defer sdl.SDLQuit()
	subsystems := sdl.SDLWasInit(sdl.SDL_INIT_NONE)
	t.Log("SDLWasInit returns", subsystems)
}

func Test_SDL_002(t *testing.T) {
	err := sdl.SDLInit(sdl.SDL_INIT_EVERYTHING)
	if err != nil {
		t.Error(err)
	}
	defer sdl.SDLQuit()
	if numdisplays, err := sdl.SDLGetNumVideoDisplays(); err != nil {
		t.Error(err)
	} else {
		for i := uint(0); i < numdisplays; i++ {
			if name, err := sdl.SDLGetDisplayName(i); err != nil {
				t.Error(err)
			} else {
				t.Log(i, name)
			}
		}
	}
}

func Test_SDL_003(t *testing.T) {
	err := sdl.SDLInit(sdl.SDL_INIT_EVERYTHING)
	if err != nil {
		t.Error(err)
	}
	defer sdl.SDLQuit()
	if numdisplays, err := sdl.SDLGetNumVideoDisplays(); err != nil {
		t.Error(err)
	} else {
		for i := uint(0); i < numdisplays; i++ {
			if ddpi, hdpi, vdpi, err := sdl.SDLGetDisplayDPI(i); err != nil {
				t.Error(err)
			} else {
				t.Logf("Display %v: ddpi=%.1f hdpi=%.1f vdpi=%.1f", i, ddpi, hdpi, vdpi)
			}
		}
	}
}

func Test_SDL_004(t *testing.T) {
	err := sdl.SDLInit(sdl.SDL_INIT_EVERYTHING)
	if err != nil {
		t.Error(err)
	}
	defer sdl.SDLQuit()
	if numdisplays, err := sdl.SDLGetNumVideoDisplays(); err != nil {
		t.Error(err)
	} else {
		for i := uint(0); i < numdisplays; i++ {
			if bounds, err := sdl.SDLGetDisplayBounds(i); err != nil {
				t.Error(err)
			} else if usableBounds, err := sdl.SDLGetDisplayUsableBounds(i); err != nil {
				t.Error(err)
			} else {
				t.Logf("Display %v: bounds=%v usable_bounds=%v", i, bounds, usableBounds)
			}
		}
	}
}
