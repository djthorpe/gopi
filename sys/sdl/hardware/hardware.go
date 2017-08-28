package hardware /* import "github.com/djthorpe/gopi/sys/sdl/hardware" */

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

// Empty Hardware configuration
type Hardware struct {
}

type hardware struct {
	log gopi.Logger
}

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
	this.log.Debug("sdl.Hardware.Open()")

	if C.SDL_Init(C.Uint32(C.SDL_INIT_EVERYTHING)) != 0 {
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
	return fmt.Sprintf("sys.sdl.Hardware{ name=\"%v\" serial=\"%v\" number_of_displays=%v }", this.Name(), this.SerialNumber(), this.NumberOfDisplays())
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
