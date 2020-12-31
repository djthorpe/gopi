package gopi

import (
	"context"
	"image"
	"net/url"
	"strings"
	"time"
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

// Flags define changes to a device
type CastFlag uint

// CastManager returns all discovered Chromecast devices,
// and allows you to connect and disconnect
type CastManager interface {
	// Return list of discovered Google Chromecast Devices
	Devices(context.Context) ([]Cast, error)

	// Connect to the control channel for a device
	Connect(Cast) error

	// Disconnect from the device
	Disconnect(Cast) error

	// LaunchAppWithId launches application with Id on a cast device.
	LaunchAppWithId(Cast, string) error

	// SetVolume sets the volume for a device, the value is between 0.0
	// and 1.0
	SetVolume(Cast, float32) error

	// SetMuted sets the volume as muted or unmuted. Does not affect the
	// volume level
	SetMuted(Cast, bool) error

	// SetPlay sets media playback state to either PLAY or STOP
	SetPlay(Cast, bool) error

	// SetPause sets media state to PLAY or PAUSE. Will not affect
	// state if currently STOP
	SetPause(Cast, bool) error

	// Seek within media stream relative to start of stream
	Seek(Cast, time.Duration) error

	// Send stop signal, terminating any playing media
	Stop(Cast) error

	// Stream URL to Chromecast supports HTTP and HTTPS protocols,
	// and the stream can be automatically started if the third
	// argument is set to true. Requires application load before
	// calling, to set the transport, or else returns an OutOfOrder
	// error
	LoadURL(Cast, *url.URL, bool) error
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

type CastEvent interface {
	Event

	Cast() Cast
	Flags() CastFlag
}

type CastService interface {
	Service
}

type CastStub interface {
	ServiceStub

	// ListCasts returns all discovered Chromecast devices
	ListCasts(ctx context.Context) ([]Cast, error)

	// SetApp loads an application into the Chromecast
	SetApp(ctx context.Context, castId, appId string) error

	// SetVolume sets the Chromecast sound volume
	SetVolume(ctx context.Context, castId string, value float32) error

	// SetMute mutes and unmutes the sound
	SetMute(ctx context.Context, castId string, value bool) error

	// LoadURL loads a video, audio or image onto the Chromecast,
	// assuming an application has already been loaded
	LoadURL(ctx context.Context, castId string, url *url.URL) error

	// Stop stops currently playing media if a media session is ongoing
	// or else resets the Chromecast to the backdrop if no media session
	Stop(ctx context.Context, castId string) error

	// Play resumes playback after paused media
	Play(ctx context.Context, castId string) error

	// Pause the media session
	Pause(ctx context.Context, castId string) error

	// Seek within playing audio or video relative to the start of the
	// playing media
	Seek(ctx context.Context, castId string, value time.Duration) error
}

// TYPES
const (
	CAST_FLAG_CONNECT CastFlag = (1 << iota)
	CAST_FLAG_VOLUME
	CAST_FLAG_MUTE
	CAST_FLAG_MEDIA
	CAST_FLAG_APP
	CAST_FLAG_DISCONNECT
	CAST_FLAG_NONE CastFlag = 0
	CAST_FLAG_MIN           = CAST_FLAG_CONNECT
	CAST_FLAG_MAX           = CAST_FLAG_DISCONNECT
)

const (
	CAST_APPID_DEFAULT      = "CC1AD845"
	CAST_APPID_MUTABLEMEDIA = "5C292C3E"
	CAST_APPID_BACKDROP     = "E8C28D3C"
)

// STRINGIFY
func (f CastFlag) String() string {
	if f == CAST_FLAG_NONE {
		return f.FlagString()
	}
	str := ""
	for v := CAST_FLAG_MIN; v <= CAST_FLAG_MAX; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.Trim(str, "|")
}

func (f CastFlag) FlagString() string {
	switch f {
	case CAST_FLAG_NONE:
		return "CAST_FLAG_NONE"
	case CAST_FLAG_CONNECT:
		return "CAST_FLAG_CONNECT"
	case CAST_FLAG_VOLUME:
		return "CAST_FLAG_VOLUME"
	case CAST_FLAG_MUTE:
		return "CAST_FLAG_MUTE"
	case CAST_FLAG_MEDIA:
		return "CAST_FLAG_MEDIA"
	case CAST_FLAG_APP:
		return "CAST_FLAG_APP"
	case CAST_FLAG_DISCONNECT:
		return "CAST_FLAG_DISCONNECT"
	default:
		return "[?? Invalid CastFlag value]"
	}
}
