package gopi

import (
	"context"
	"image"
)

////////////////////////////////////////////////////////////////////////////////
// ARGON ONE CASE

// Ref: https://github.com/Argon40Tech/Argon-ONE-i2c-Codes
type ArgonOnePowerMode uint

// ArgonOne interfaces with the ArgonOne case for the
// Raspberry Pi
type ArgonOne interface {
	// Set fan duty cycle (0-100)
	SetFan(uint8) error

	// Set Power Mode
	SetPower(ArgonOnePowerMode) error
}

// CONSTANTS
const (
	ARGONONE_POWER_DEFAULT ArgonOnePowerMode = iota
	ARGONONE_POWER_ALWAYSON
	ARGONONE_POWER_UART
)

// STRINGIFY
func (v ArgonOnePowerMode) String() string {
	switch v {
	case ARGONONE_POWER_DEFAULT:
		return "ARGONONE_POWER_DEFAULT"
	case ARGONONE_POWER_ALWAYSON:
		return "ARGONONE_POWER_ALWAYSON"
	case ARGONONE_POWER_UART:
		return "ARGONONE_POWER_UART"
	default:
		return "[?? Invalid ArgonOnePowerMode value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// WAVESHARE E-PAPER DISPLAY (EPD)

type EPD interface {
	// Return screen size
	Size() Size

	// Clear display
	Clear(context.Context) error

	// Draw image on the display
	Draw(context.Context, image.Image) error

	// Size image and draw on the display. A size value of
	// 1.0 is equivalent to calling Draw
	DrawSized(context.Context, float64, image.Image) error

	// Sleep display
	Sleep() error
}

////////////////////////////////////////////////////////////////////////////////
// GOOGLE CHROMECAST

// CastManager returns all discovered Chromecast devices,
// and allows you to connect and disconnect
type CastManager interface {
	// Return list of discovered Google Chromecast Devices
	Devices(context.Context) ([]Cast, error)

	// Connect to the control channel for a device
	Connect(Cast) error

	// Disconnect from the device
	Disconnect(Cast) error
}

// Cast represents a Google Chromecast device
type Cast interface {
	// Id returns the identifier for a chromecast
	Id() string

	// Name returns the readable name for a chromecast
	Name() string

	// Model returns the reported model information
	Model() string

	// Service returns the currently running service
	Service() string

	// State returns 0 if backdrop, else returns 1
	State() uint
}
