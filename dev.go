package gopi

import (
	"context"
	"image"
	"net"
	"net/url"
	"strings"
	"time"
)

/*
	This file contains interface defininitons for example devices:

	* Argon One case for Raspberry Pi (GPIO, Infrared, Fan Control)
	* eInk Paper Displays (GPIO, SPI, Bitmaps)
	* Google Chromecast control (mDNS, RPC, Protocol Buffers)
	* Rotel Amplifer control (via RS232)
	* IKEA Tradfri Zigbee Gateway

	Ultimately these should be split out into separate repos...
*/

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
// IKEA TRADFRI GATEWAY

// TradfriManager communicates with an Ikea Tradfri gateway
type TradfriManager interface {
	// Connect to a gateway with gateway id, hostname and port
	Connect(string, string, uint16) error

	// Disconnect from a gateway
	Disconnect() error

	// Return all devices
	Devices(context.Context) ([]TradfriDevice, error)

	// Observe for device changes
	ObserveDevice(context.Context, TradfriDevice) error

	// Properties
	Addr() net.Addr  // Return IP Address for Gateway
	Id() string      // Return ID for authentication to gateway
	Version() string // Return version of gateway
}

// TradfriDevice is a device connected to the gateway
type TradfriDevice interface {
	Name() string
	Id() uint
	Type() uint
	Created() time.Time
	Updated() time.Time
	Active() bool
	Vendor() string
	Product() string
	Version() string
}

////////////////////////////////////////////////////////////////////////////////
// ROTEL AMPLIFIER (RS232) CONTROL

// RotelManager controls a connected Rotel Amplifier
type RotelManager interface {
	Publisher

	// Get model number
	Model() string

	// Get properties
	Power() bool
	Source() string
	Volume() uint
	Freq() string
	Bass() int
	Treble() int
	Muted() bool
	Bypass() bool
	Balance() (string, uint)
	Speakers() []string
	Dimmer() uint

	// Set properties
	SetPower(bool) error           // SetPower sets amplifier to standby or on
	SetSource(string) error        // SetSource sets input source
	SetVolume(uint) error          // SetVolume sets the volume between 1 and 96 inclusive
	SetMute(bool) error            // SetMute mutes and unmutes
	SetBypass(bool) error          // SetBypass sets preamp bypass
	SetTreble(int) error           // SetTreble sets treble -10 <> +10
	SetBass(int) error             // SetBass sets bass -10 <> +10
	SetBalance(string, uint) error // L,R between 0 and 15
	SetDimmer(uint) error          // SetDimmer display between 0 and 6 (0 is brightest)

	// Actions
	Play() error
	Stop() error
	Pause() error
	NextTrack() error
	PrevTrack() error
}

// RotelService defines an RPC service connected to the Rotel Amplifer
type RotelService interface {
	Service
}

// RotelEvent is emitted on change of amplifier state
type RotelEvent interface {
	Event
}

// RotelStub is an RPC client which connects to the RPC service
type RotelStub interface {
	ServiceStub

	// Set Properties
	SetPower(context.Context, bool) error // SetPower to on or standby
	SetSource(context.Context, string) error
	SetVolume(context.Context, uint) error
	SetMute(context.Context, bool) error
	SetBypass(context.Context, bool) error
	SetTreble(context.Context, int) error
	SetBass(context.Context, int) error
	SetBalance(context.Context, string, uint) error
	SetDimmer(context.Context, uint) error

	// Actions
	Play(context.Context) error
	Stop(context.Context) error
	Pause(context.Context) error
	NextTrack(context.Context) error
	PrevTrack(context.Context) error

	// Stream change events
	Stream(context.Context, chan<- RotelEvent) error
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
	SeekAbs(Cast, time.Duration) error

	// Seek within media stream relative to current position
	SeekRel(Cast, time.Duration) error

	// Send stop signal, terminating any playing media
	Stop(Cast) error

	// Stream URL to Chromecast supports HTTP and HTTPS protocols,
	// and the stream can be automatically started if the third
	// argument is set to true. Requires application load before
	// calling, to set the transport, or else returns an OutOfOrder
	// error
	LoadURL(Cast, *url.URL, bool) error

	// Returns current volume state (level and muted)
	Volume(Cast) (float32, bool, error)

	// Returns app state
	App(Cast) (CastApp, error)
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

// CastApp represents an application running on the Chromecast
type CastApp interface {
	// Id returns the identifier for the application
	Id() string

	// Name returns the name of the application
	Name() string

	// Status is the current status of the application
	Status() string
}

type CastEvent interface {
	Event

	Flags() CastFlag
	Cast() Cast
	App() CastApp
	Volume() (float32, bool)
}

type CastService interface {
	Service
}

type CastStub interface {
	ServiceStub

	// ListCasts returns all discovered Chromecast devices
	ListCasts(ctx context.Context) ([]Cast, error)

	// Return Volume
	Volume(ctx context.Context, castId string) (float32, bool, error)

	// SetVolume sets the Chromecast sound volume
	SetVolume(ctx context.Context, castId string, value float32) error

	// SetMute mutes and unmutes the sound
	SetMute(ctx context.Context, castId string, value bool) error

	// Return App
	App(ctx context.Context, castId string) (CastApp, error)

	// SetApp loads an application into the Chromecast
	SetApp(ctx context.Context, castId, appId string) error

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

	// SeekAbs within playing audio or video relative to the start of the
	// playing media
	SeekAbs(ctx context.Context, castId string, value time.Duration) error

	// SeekRel within playing audio or video relative to the current position
	SeekRel(ctx context.Context, castId string, value time.Duration) error

	// Stream emits events from Chromecasts, filtered
	// by the id of the chromecast  until context is cancelled. Where
	// the id filter is empty, all connected chromecast events are emitted
	Stream(context.Context, string, chan<- CastEvent) error
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
