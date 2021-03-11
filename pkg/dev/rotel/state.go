package rotel

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type State struct {
	model        string
	power        string
	update       string // rs232 update
	volume, mute string
	bass, treble string
	balance      []string
	source       string
	freq         string
	bypass       string
	speaker      string
	dimmer       string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	commands = []struct {
		re *regexp.Regexp
		fn func(this *State, args []string) (RotelFlag, error)
	}{
		{regexp.MustCompile("^model=(\\w+)$"), SetModel},
		{regexp.MustCompile("^power=(on|standby)$"), SetPower},
		{regexp.MustCompile("^volume=(\\d+)$"), SetVolume},
		{regexp.MustCompile("^update_mode=(auto|manual)$"), SetUpdateMode},
		{regexp.MustCompile("^update_mode=(auto|manual)$"), SetUpdateMode},
		{regexp.MustCompile("^bass=([\\+\\-]?\\d+)$"), SetBass},
		{regexp.MustCompile("^treble=([\\+\\-]?\\d+)$"), SetTreble},
		{regexp.MustCompile("^balance=([LR]?)(\\d+)$"), SetBalance},
		{regexp.MustCompile("^mute=(on|off)$"), SetMute},
		{regexp.MustCompile("^source=(\\w+)$"), SetSource},
		{regexp.MustCompile("^freq=(.+)$"), SetFreq},
		{regexp.MustCompile("^bypass=(on|off)$"), SetBypass},
		{regexp.MustCompile("^speaker=(a|b|a_b|off)$"), SetSpeaker},
		{regexp.MustCompile("^dimmer=(\\d+)$"), SetDimmer},
	}
)

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *State) Model() string {
	return this.model
}

func (this *State) Power() bool {
	return this.power == "on"
}

func (this *State) Volume() uint {
	if this.power == "on" {
		if vol, err := strconv.ParseUint(this.volume, 0, 32); err == nil {
			return uint(vol)
		}
	}
	return 0
}

func (this *State) Source() string {
	if this.power == "on" {
		return this.source
	} else {
		return ""
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

// Update returns a query to get state
func (this *State) Update() string {
	switch {
	case this.model == "":
		return "model?"
	case this.power == "":
		return "power?"
	case this.power != "on": // When power is off, don't read other values
		return ""
	case this.update == "":
		return "rs232_update_on!"
	case this.volume == "":
		return "volume?"
	case this.source == "":
		return "source?"
	case this.freq == "":
		return "freq?"
	case this.bypass == "":
		return "bypass?"
	case this.speaker == "":
		return "speaker?"
	case this.mute == "":
		return "mute?"
	case this.bass == "":
		return "bass?"
	case this.treble == "":
		return "treble?"
	case this.balance == nil:
		return "balance?"
	case this.dimmer == "":
		return "dimmer?"
	}

	// By default, no state needs read
	return ""
}

// Set sets state from data coming from amp
func (this *State) Set(param string) (RotelFlag, error) {
	for _, command := range commands {
		if args := command.re.FindStringSubmatch(param); len(args) != 0 {
			return command.fn(this, args[1:])
		}
	}
	// Cannot match command
	return 0, gopi.ErrUnexpectedResponse.WithPrefix(strconv.Quote(param))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *State) String() string {
	str := "<state"
	if model := this.Model(); model != "" {
		str += fmt.Sprintf(" model=%q", this.model)
	}
	str += fmt.Sprint(" power=", this.Power())
	if source := this.Source(); source != "" {
		str += fmt.Sprintf(" source=%q", source)
	}
	if vol := this.Volume(); vol != 0 {
		str += fmt.Sprint(" vol=", vol)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func SetModel(this *State, args []string) (RotelFlag, error) {
	if args[0] == "" {
		return 0, gopi.ErrBadParameter.WithPrefix("SetModel")
	} else if this.model != args[0] {
		this.model = args[0]
		return FLAG_MODEL, nil
	}
	return 0, nil
}

func SetPower(this *State, args []string) (RotelFlag, error) {
	if args[0] == "" {
		return 0, gopi.ErrBadParameter.WithPrefix("SetPower")
	} else if this.power != args[0] {
		this.power = args[0]
		return FLAG_POWER, nil
	}
	return 0, nil
}

func SetUpdateMode(this *State, args []string) (RotelFlag, error) {
	if args[0] == "" {
		return 0, gopi.ErrBadParameter.WithPrefix("SetUpdateMode")
	}
	this.update = args[0]
	return 0, nil
}

func SetVolume(this *State, args []string) (RotelFlag, error) {
	if volume, err := strconv.ParseUint(args[0], 0, 32); err != nil {
		return 0, err
	} else if volume_ := fmt.Sprint(volume); volume_ != this.volume {
		this.volume = volume_
		return FLAG_VOLUME, nil
	}
	return 0, nil
}

func SetBass(this *State, args []string) (RotelFlag, error) {
	if bass, err := strconv.ParseInt(args[0], 0, 32); err != nil {
		return 0, err
	} else if bass_ := fmt.Sprint(bass); bass_ != this.bass {
		this.bass = bass_
		return FLAG_BASS, nil
	}
	return 0, nil
}

func SetTreble(this *State, args []string) (RotelFlag, error) {
	if treble, err := strconv.ParseInt(args[0], 0, 32); err != nil {
		return 0, err
	} else if treble_ := fmt.Sprint(treble); treble_ != this.treble {
		this.treble = treble_
		return FLAG_TREBLE, nil
	}
	return 0, nil
}

func SetBalance(this *State, args []string) (RotelFlag, error) {
	if scalar, err := strconv.ParseUint(args[1], 0, 32); err != nil {
		return 0, err
	} else {
		scalar_ := fmt.Sprint(scalar)
		if this.balance == nil || scalar_ != this.balance[1] || args[0] != this.balance[0] {
			this.balance = []string{args[0], fmt.Sprint(scalar)}
			return FLAG_BALANCE, nil
		}
	}
	return 0, nil
}

func SetMute(this *State, args []string) (RotelFlag, error) {
	if args[0] != this.mute {
		this.mute = args[0]
		return FLAG_MUTE, nil
	}
	return 0, nil
}

func SetSource(this *State, args []string) (RotelFlag, error) {
	if args[0] != this.source {
		this.source = args[0]
		return FLAG_SOURCE, nil
	}
	return 0, nil
}

func SetFreq(this *State, args []string) (RotelFlag, error) {
	if args[0] != this.freq {
		this.freq = args[0]
		return FLAG_FREQ, nil
	}
	return 0, nil
}

func SetBypass(this *State, args []string) (RotelFlag, error) {
	if args[0] != this.bypass {
		this.bypass = args[0]
		return FLAG_BYPASS, nil
	}
	return 0, nil
}

func SetSpeaker(this *State, args []string) (RotelFlag, error) {
	if args[0] != this.speaker {
		this.speaker = args[0]
		return FLAG_SPEAKER, nil
	}
	return 0, nil
}

func SetDimmer(this *State, args []string) (RotelFlag, error) {
	if dimmer, err := strconv.ParseUint(args[0], 0, 32); err != nil {
		return 0, err
	} else if dimmer_ := fmt.Sprint(dimmer); this.dimmer != dimmer_ {
		this.dimmer = dimmer_
		return FLAG_DIMMER, nil
	}
	return 0, nil
}
