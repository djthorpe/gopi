/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package macsdl

// Macintosh framework version

// #cgo LDFLAGS: -framework SDL2
// #include <SDL2/SDL.h>
import "C"

import (
	"errors"
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCT

// Subsystem initialization flags
type SDLInitFlags uint32

// Hardware configuration
type Hardware struct {
	Init SDLInitFlags
}

type hardware struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Which subsystems to start, the default is everything
	SDL_INIT_NONE           SDLInitFlags = iota
	SDL_INIT_EVERYTHING     SDLInitFlags = C.SDL_INIT_EVERYTHING
	SDL_INIT_TIMER          SDLInitFlags = C.SDL_INIT_TIMER
	SDL_INIT_AUDIO          SDLInitFlags = C.SDL_INIT_AUDIO
	SDL_INIT_VIDEO          SDLInitFlags = C.SDL_INIT_VIDEO
	SDL_INIT_JOYSTICK       SDLInitFlags = C.SDL_INIT_JOYSTICK
	SDL_INIT_HAPTIC         SDLInitFlags = C.SDL_INIT_HAPTIC
	SDL_INIT_GAMECONTROLLER SDLInitFlags = C.SDL_INIT_GAMECONTROLLER
	SDL_INIT_EVENTS         SDLInitFlags = C.SDL_INIT_EVENTS
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware and display
	gopi.RegisterModule(gopi.Module{Name: "sdl/hardware", Type: gopi.MODULE_TYPE_HARDWARE, New: newHardware})
}

func newHardware(config *gopi.AppConfig, logger gopi.Logger) (gopi.Driver, error) {
	var err gopi.Error
	if driver, ok := gopi.Open2(Hardware{}, logger, &err).(gopi.HardwareDriver2); !ok {
		return nil, err
	} else {
		return driver, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - HARDWARE

// Open
func (config Hardware) Open(logger gopi.Logger) (gopi.Driver, error) {
	this := new(hardware)
	this.log = logger

	if config.Init == SDL_INIT_NONE {
		config.Init = SDL_INIT_EVERYTHING
	}

	this.log.Debug("sdl.Hardware.Open(Init=%v)", config.Init)

	if C.SDL_Init(C.Uint32(config.Init)) != 0 {
		return nil, getError()
	}

	return this, nil
}

// GetName returns the name of the hardware
func (this *hardware) Name() string {
	return string(C.GoString(C.SDL_GetPlatform()))
}

// Return serial number
func (this *hardware) SerialNumber() string {
	return "NOT IMPLEMENTED"
}

// Revision returns the SDL revision string
func (this *hardware) Revision() string {
	return (string)(C.GoString(C.SDL_GetRevision()))
}

func (this *hardware) NumberOfDisplays() uint {
	n := int(C.SDL_GetNumVideoDisplays())
	if n < 0 {
		return 0
	}
	return uint(n)
}

func (this *hardware) Close() error {
	this.log.Debug("sdl.Hardware.Close()")
	C.SDL_Quit()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *hardware) String() string {
	return fmt.Sprintf("sys.sdl.Hardware{ name=\"%v\" revision=\"%v\" number_of_displays=%v }", this.Name(), this.Revision(), this.NumberOfDisplays())
}

func (v SDLInitFlags) String() string {
	switch v {
	case SDL_INIT_NONE:
		return "SDL_INIT_NONE"
	case SDL_INIT_EVERYTHING:
		return "SDL_INIT_EVERYTHING"
	case SDL_INIT_TIMER:
		return "SDL_INIT_TIMER"
	case SDL_INIT_AUDIO:
		return "SDL_INIT_AUDIO"
	case SDL_INIT_VIDEO:
		return "SDL_INIT_VIDEO"
	case SDL_INIT_JOYSTICK:
		return "SDL_INIT_JOYSTICK"
	case SDL_INIT_HAPTIC:
		return "SDL_INIT_HAPTIC"
	case SDL_INIT_GAMECONTROLLER:
		return "SDL_INIT_GAMECONTROLLER"
	case SDL_INIT_EVENTS:
		return "SDL_INIT_EVENTS"
	default:
		return "[?? Invalid SDLInitFlags value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

func getError() error {
	if err := C.SDL_GetError(); err != nil {
		gostr := C.GoString(err)
		if len(gostr) > 0 {
			return errors.New(gostr)
		}
	}
	return nil
}
