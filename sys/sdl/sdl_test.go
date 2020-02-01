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
	}
}
