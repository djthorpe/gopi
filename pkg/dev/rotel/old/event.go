/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	// Frameworks
	"fmt"

	gopi "github.com/djthorpe/gopi"
	rotel "github.com/djthorpe/rotel"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

type evt struct {
	source gopi.Driver
	typ    rotel.EventType
	state  rotel.RotelState
}

////////////////////////////////////////////////////////////////////////////////
// EMIT EVENTS

func (this *driver) evtPower(value rotel.Power) {
	if this.state.Power != value {
		this.state.Power = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_POWER,
			state:  this.state,
		})
	}
}

func (this *driver) evtSource(value rotel.Source) {
	if this.state.Source != value {
		this.state.Source = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_SOURCE,
			state:  this.state,
		})
	}
}

func (this *driver) evtVolume(value rotel.Volume) {
	if this.state.Volume != value {
		this.state.Volume = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_VOLUME,
			state:  this.state,
		})
	}
}

func (this *driver) evtFreq(value string) {
	if this.state.Freq != value {
		this.state.Freq = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_FREQ,
			state:  this.state,
		})
	}
}

func (this *driver) evtMute(value rotel.Mute) {
	if this.state.Mute != value {
		this.state.Mute = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_MUTE,
			state:  this.state,
		})
	}
}

func (this *driver) evtBypass(value rotel.Bypass) {
	if this.state.Bypass != value {
		this.state.Bypass = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_BYPASS,
			state:  this.state,
		})
	}
}

func (this *driver) evtBass(value rotel.Tone) {
	if this.state.Bass != value {
		this.state.Bass = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_BASS,
			state:  this.state,
		})
	}
}

func (this *driver) evtTreble(value rotel.Tone) {
	if this.state.Treble != value {
		this.state.Treble = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_TREBLE,
			state:  this.state,
		})
	}
}

func (this *driver) evtBalance(value rotel.Balance) {
	if this.state.Balance != value {
		this.state.Balance = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_BALANCE,
			state:  this.state,
		})
	}
}

func (this *driver) evtDimmer(value rotel.Dimmer) {
	if this.state.Dimmer != value {
		this.state.Dimmer = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_DIMMER,
			state:  this.state,
		})
	}
}

func (this *driver) evtSpeaker(value rotel.Speaker) {
	if this.state.Speaker != value {
		this.state.Speaker = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_SPEAKER,
			state:  this.state,
		})
	}
}

func (this *driver) evtUpdate(value rotel.Update) {
	if this.state.Update != value {
		this.state.Update = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_UPDATE,
			state:  this.state,
		})
	}
}

////////////////////////////////////////////////////////////////////////////////
// EVENT IMPLEMENTATION

func (this *evt) Name() string {
	return "RotelEvent"
}

func (this *evt) Source() gopi.Driver {
	return this.source
}

func (this *evt) Type() rotel.EventType {
	return this.typ
}

func (this *evt) State() rotel.RotelState {
	return this.state
}

func (this *evt) String() string {
	switch this.typ {
	case rotel.EVENT_TYPE_POWER:
		return fmt.Sprintf("<rotel.Event>{ type=%v power=%v }", this.typ, this.state.Power)
	case rotel.EVENT_TYPE_SOURCE:
		return fmt.Sprintf("<rotel.Event>{ type=%v source=%v }", this.typ, this.state.Source)
	case rotel.EVENT_TYPE_VOLUME:
		return fmt.Sprintf("<rotel.Event>{ type=%v volume=%v }", this.typ, this.state.Volume)
	case rotel.EVENT_TYPE_MUTE:
		return fmt.Sprintf("<rotel.Event>{ type=%v mute=%v }", this.typ, this.state.Mute)
	case rotel.EVENT_TYPE_FREQ:
		return fmt.Sprintf("<rotel.Event>{ type=%v freq=%v }", this.typ, this.state.Freq)
	case rotel.EVENT_TYPE_BASS:
		return fmt.Sprintf("<rotel.Event>{ type=%v bass=%v }", this.typ, this.state.Bass)
	case rotel.EVENT_TYPE_TREBLE:
		return fmt.Sprintf("<rotel.Event>{ type=%v treble=%v }", this.typ, this.state.Treble)
	case rotel.EVENT_TYPE_BYPASS:
		return fmt.Sprintf("<rotel.Event>{ type=%v bypass=%v }", this.typ, this.state.Bypass)
	case rotel.EVENT_TYPE_BALANCE:
		return fmt.Sprintf("<rotel.Event>{ type=%v balance=%v }", this.typ, this.state.Balance)
	case rotel.EVENT_TYPE_SPEAKER:
		return fmt.Sprintf("<rotel.Event>{ type=%v speaker=%v }", this.typ, this.state.Speaker)
	case rotel.EVENT_TYPE_DIMMER:
		return fmt.Sprintf("<rotel.Event>{ type=%v dimmer=%v }", this.typ, this.state.Dimmer)
	case rotel.EVENT_TYPE_UPDATE:
		return fmt.Sprintf("<rotel.Event>{ type=%v update=%v }", this.typ, this.state.Update)
	default:
		return fmt.Sprintf("<rotel.Event>{ type=%v }", this.typ)
	}
}
