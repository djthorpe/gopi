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
	SetPower(bool) error     // SetPower sets amplifier to standby or on
	SetSource(string) error  // SetSource sets input source
	SetVolume(uint) error    // SetVolume sets the volume between 1 and 96 inclusive
	SetMute(bool) error      // SetMute mutes and unmutes
	SetBypass(bool) error    // SetBypass sets preamp bypass
	SetTreble(int) error     // SetTreble sets treble -10 <> +10
	SetBass(int) error       // SetBass sets bass -10 <> +10
	SetBalance(string) error // SetBalance L,R or 0
	SetDimmer(uint) error    // SetDimmer display between 0 and 6 (0 is brightest)

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

	Flags() RotelFlag // Flags returns the state that has changed
}

// RotelFlag provides flags on state changes
type RotelFlag uint16

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
	SetBalance(context.Context, string) error // L, R or 0
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

const (
	ROTEL_FLAG_POWER RotelFlag = (1 << iota)
	ROTEL_FLAG_VOLUME
	ROTEL_FLAG_MUTE
	ROTEL_FLAG_BASS
	ROTEL_FLAG_TREBLE
	ROTEL_FLAG_BALANCE
	ROTEL_FLAG_SOURCE
	ROTEL_FLAG_FREQ
	ROTEL_FLAG_BYPASS
	ROTEL_FLAG_SPEAKER
	ROTEL_FLAG_DIMMER
	ROTEL_FLAG_NONE RotelFlag = 0
	ROTEL_FLAG_MIN            = ROTEL_FLAG_POWER
	ROTEL_FLAG_MAX            = ROTEL_FLAG_DIMMER
)

func (f RotelFlag) String() string {
	if f == ROTEL_FLAG_NONE {
		return f.FlagString()
	}
	str := ""
	for v := ROTEL_FLAG_MIN; v <= ROTEL_FLAG_MAX; v <<= 1 {
		if v&f == v {
			str += "|" + v.FlagString()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (f RotelFlag) FlagString() string {
	switch f {
	case ROTEL_FLAG_NONE:
		return "ROTEL_FLAG_NONE"
	case ROTEL_FLAG_POWER:
		return "ROTEL_FLAG_POWER"
	case ROTEL_FLAG_VOLUME:
		return "ROTEL_FLAG_VOLUME"
	case ROTEL_FLAG_MUTE:
		return "ROTEL_FLAG_MUTE"
	case ROTEL_FLAG_BASS:
		return "ROTEL_FLAG_BASS"
	case ROTEL_FLAG_TREBLE:
		return "ROTEL_FLAG_TREBLE"
	case ROTEL_FLAG_BALANCE:
		return "ROTEL_FLAG_BALANCE"
	case ROTEL_FLAG_SOURCE:
		return "ROTEL_FLAG_SOURCE"
	case ROTEL_FLAG_FREQ:
		return "ROTEL_FLAG_FREQ"
	case ROTEL_FLAG_BYPASS:
		return "ROTEL_FLAG_BYPASS"
	case ROTEL_FLAG_SPEAKER:
		return "ROTEL_FLAG_SPEAKER"
	case ROTEL_FLAG_DIMMER:
		return "ROTEL_FLAG_DIMMER"
	default:
		return "[?? Invalid RotelFlag value]"
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

	// Return a chromecast by id or name, returns nil if not found
	Get(string) Cast

	// Connect to the control channel for a device
	Connect(context.Context, Cast) error

	// Disconnect from the device
	Disconnect(Cast) error

	// SetVolume sets the volume for a device, the value is between 0.0 and 1.0
	SetVolume(context.Context, Cast, float32) error

	// SetMuted sets the volume as muted or unmuted. Does not affect the
	// volume level
	SetMuted(context.Context, Cast, bool) error

	// LaunchAppWithId launches application with Id on a cast device.
	LaunchAppWithId(context.Context, Cast, string) error

	// ConnectMedia starts a media session
	ConnectMedia(context.Context, Cast) error

	// DisconnectMedia ends a media session
	DisconnectMedia(context.Context, Cast) error

	// LoadMedia loads a video, audio, webpage or image onto the Chromecast,
	// assuming an application has already been loaded. Autoplay parameter
	// starts media playback immediately
	LoadMedia(context.Context, Cast, *url.URL, bool) error

	// SetPlay sets media playback state to either PLAY or STOP
	//SetPlay(context.Context, Cast, bool) error

	// SetPause sets media state to PLAY or PAUSE. Will not affect
	// state if currently STOP
	//SetPause(context.Context, Cast, bool) error

	/*

		// Seek within media stream relative to start of stream
		SeekAbs(Cast, time.Duration) error

		// Seek within media stream relative to current position
		SeekRel(Cast, time.Duration) error

		// Send stop signal, terminating any playing media
		Stop(Cast) error
	*/
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

	// Returns current volume state (level,0->1 and muted)
	// and will return 0,false if not known (not connected to
	// cast device)
	Volume() (float32, bool)
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
	//	App() CastApp
	//	Volume() (float32, bool)
}

type CastService interface {
	Service
}

type CastStub interface {
	ServiceStub

	// List returns all discovered Chromecast devices within
	// a certain time
	List(context.Context, time.Duration) ([]Cast, error)

	// Stream emits chromecast events on a channel
	Stream(context.Context, chan<- CastEvent) error

	// Connect to chromecast
	Connect(context.Context, string) (Cast, error)

	// Connect to chromecast
	Disconnect(context.Context, string) error

	// Initiate Media Session
	ConnectMedia(context.Context, string) (Cast, error)

	// Terminate Media Session
	DisconnectMedia(context.Context, string) (Cast, error)

	// SetVolume sets the sound volume
	SetVolume(ctx context.Context, key string, level float32) (Cast, error)

	// SetMuted mutes and unmutes the sound
	SetMuted(ctx context.Context, key string, value bool) (Cast, error)

	// LaunchAppWithId loads an application into the Chromecast
	LaunchAppWithId(context.Context, string, string) (Cast, error)

	// LoadMedia to chromecast from URL and with autoplay flag
	LoadMedia(context.Context, string, *url.URL, bool) (Cast, error)

	// Play resumes playback after paused media
	//Play(context.Context, string) error

	// Pause the media session
	//Pause(context.Context, string) error

	/*

		// Stop stops currently playing media if a media session is ongoing
		// or else resets the Chromecast to the backdrop if no media session
		Stop(ctx context.Context, castId string) error

		// SeekAbs within playing audio or video relative to the start of the
		// playing media
		SeekAbs(ctx context.Context, castId string, value time.Duration) error

		// SeekRel within playing audio or video relative to the current position
		SeekRel(ctx context.Context, castId string, value time.Duration) error
	*/
}

// TYPES
const (
	CAST_FLAG_CONNECT CastFlag = (1 << iota)
	CAST_FLAG_DISCOVERY
	CAST_FLAG_NAME
	CAST_FLAG_APP
	CAST_FLAG_VOLUME
	CAST_FLAG_MUTE
	CAST_FLAG_MEDIA
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
	case CAST_FLAG_DISCOVERY:
		return "CAST_FLAG_DISCOVERY"
	case CAST_FLAG_VOLUME:
		return "CAST_FLAG_VOLUME"
	case CAST_FLAG_MUTE:
		return "CAST_FLAG_MUTE"
	case CAST_FLAG_MEDIA:
		return "CAST_FLAG_MEDIA"
	case CAST_FLAG_APP:
		return "CAST_FLAG_APP"
	case CAST_FLAG_NAME:
		return "CAST_FLAG_NAME"
	case CAST_FLAG_DISCONNECT:
		return "CAST_FLAG_DISCONNECT"
	default:
		return "[?? Invalid CastFlag value]"
	}
}
