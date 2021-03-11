package rotel

import (
	"fmt"
	"strings"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	RotelFlag
	*State
}

type RotelFlag uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FLAG_POWER RotelFlag = (1 << iota)
	FLAG_VOLUME
	FLAG_MUTE
	FLAG_BASS
	FLAG_TREBLE
	FLAG_BALANCE
	FLAG_SOURCE
	FLAG_FREQ
	FLAG_BYPASS
	FLAG_SPEAKER
	FLAG_DIMMER
	FLAG_MODEL
	FLAG_NONE RotelFlag = 0
	FLAG_MIN            = FLAG_POWER
	FLAG_MAX            = FLAG_MODEL
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEvent(state *State, flag RotelFlag) gopi.Event {
	return &event{flag, state}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (*event) Name() string {
	return "rotel"
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<event"
	str += fmt.Sprintf(" name=%q", this.Name())
	if this.RotelFlag != FLAG_NONE {
		str += fmt.Sprint(" flags=", this.RotelFlag)
	}
	if this.State != nil {
		str += fmt.Sprint(" ", this.State)
	}
	return str + ">"
}

func (f RotelFlag) String() string {
	if f == FLAG_NONE {
		return f.FlagString()
	}
	str := ""
	for v := FLAG_MIN; v <= FLAG_MAX; v <<= 1 {
		if v&f == v {
			str += "|" + v.FlagString()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (f RotelFlag) FlagString() string {
	switch f {
	case FLAG_NONE:
		return "FLAG_NONE"
	case FLAG_POWER:
		return "FLAG_POWER"
	case FLAG_VOLUME:
		return "FLAG_VOLUME"
	case FLAG_MUTE:
		return "FLAG_MUTE"
	case FLAG_BASS:
		return "FLAG_BASS"
	case FLAG_TREBLE:
		return "FLAG_TREBLE"
	case FLAG_BALANCE:
		return "FLAG_BALANCE"
	case FLAG_SOURCE:
		return "FLAG_SOURCE"
	case FLAG_FREQ:
		return "FLAG_FREQ"
	case FLAG_BYPASS:
		return "FLAG_BYPASS"
	case FLAG_SPEAKER:
		return "FLAG_SPEAKER"
	case FLAG_DIMMER:
		return "FLAG_DIMMER"
	case FLAG_MODEL:
		return "FLAG_MODEL"
	default:
		return "[?? Invalid RotelFlag value]"
	}
}
