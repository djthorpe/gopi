// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type State struct {

	// Key state
	keyCode   gopi.KeyCode
	scanCode  uint32
	keyAction gopi.KeyAction
	keyState  gopi.KeyState

	// Position state
	rel      gopi.Point
	position gopi.Point
	last     gopi.Point
	slot     uint32
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *State) Reset() {
	this.keyCode = gopi.KEYCODE_NONE
	this.scanCode = 0
	this.keyAction = gopi.KEYACTION_NONE
	this.keyState = gopi.KEYSTATE_NONE
	this.rel = gopi.ZeroPoint
	this.slot = 0xFFFFFFFF
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *State) KeyState() gopi.KeyState {
	return this.keyState
}

func (this *State) Position() gopi.Point {
	return this.position
}

func (this *State) SetPosition(pt gopi.Point) {
	this.position = pt
}

////////////////////////////////////////////////////////////////////////////////
// DECODE EVEvent and emit events

func (this *State) Decode(evt linux.EVEvent) *event {

	switch evt.Type {
	case linux.EV_KEY:
		evDecodeKey(evt, this)
	case linux.EV_REL:
		evDecodeRel(evt, this)
	case linux.EV_ABS:
		evDecodeAbs(evt, this)
	case linux.EV_MSC:
		evDecodeMsc(evt, this)
	case linux.EV_SYN:
		if evt := evDecodeSyn(evt, this); evt != nil {
			return evt
		}
	}

	// By default don't emit and event
	return nil
}

func (this *State) DecodeKey(evt linux.EVEvent) {
	// Interpret key code and key action (up, down and repeat)
	code, action := gopi.KeyCode(evt.Code), gopi.KeyAction(evt.Value)

	// Alter key state if a modified key was pressed and also
	// handle sticky state keys CAPS, NUM and SCROLL locks
	state := gopi.KEYSTATE_NONE
	switch code {
	case gopi.KEYCODE_CAPSLOCK:
		if action == gopi.KEYACTION_KEY_DOWN {
			this.keyState ^= gopi.KEYSTATE_CAPSLOCK
			// TODO: Change LED's
		}
	case gopi.KEYCODE_NUMLOCK:
		if action == gopi.KEYACTION_KEY_DOWN {
			this.keyState ^= gopi.KEYSTATE_NUMLOCK
			// TODO: Change LED's
		}
	case gopi.KEYCODE_SCROLLLOCK:
		if action == gopi.KEYACTION_KEY_DOWN {
			this.keyState ^= gopi.KEYSTATE_SCROLLLOCK
			// TODO: Change LED's
		}
	case gopi.KEYCODE_LEFTSHIFT:
		state = gopi.KEYSTATE_LEFTSHIFT
	case gopi.KEYCODE_RIGHTSHIFT:
		state = gopi.KEYSTATE_RIGHTSHIFT
	case gopi.KEYCODE_LEFTCTRL:
		state = gopi.KEYSTATE_LEFTCTRL
	case gopi.KEYCODE_RIGHTCTRL:
		state = gopi.KEYSTATE_RIGHTCTRL
	case gopi.KEYCODE_LEFTALT:
		state = gopi.KEYSTATE_LEFTALT
	case gopi.KEYCODE_RIGHTALT:
		state = gopi.KEYSTATE_RIGHTALT
	case gopi.KEYCODE_LEFTMETA:
		state = gopi.KEYSTATE_LEFTMETA
	case gopi.KEYCODE_RIGHTMETA:
		state = gopi.KEYSTATE_RIGHTMETA
	}

	// Set device code and action
	this.keyCode = code
	this.keyAction = action

	// Set device state from key action
	if state != gopi.KEYSTATE_NONE {
		if action == gopi.KEYACTION_KEY_DOWN || action == gopi.KEYACTION_KEY_REPEAT {
			this.keyState |= state
		} else if action == gopi.KEYACTION_KEY_UP {
			this.keyState ^= state
		}
	}
}

func (this *State) DecodeAbs(evt linux.EVEvent) {
	switch evt.Code {
	case linux.EV_CODE_X:
		this.position.X = float32(int32(evt.Value))
	case linux.EV_CODE_Y:
		this.position.Y = float32(int32(evt.Value))
	case linux.EV_CODE_SLOT:
		this.slot = evt.Value
	case linux.EV_CODE_SLOT_ID, linux.EV_CODE_SLOT_X, linux.EV_CODE_SLOT_Y:
		this.Log.Debug("Ignoring multi-touch event:", evt)
	default:
		this.Log.Debug("Ignoring event:", evt)
	}
}

func (this *State) DecodeRel(evt linux.EVEvent) {
	switch evt.Code {
	case linux.EV_CODE_X:
		this.rel.X = float32(int32(evt.Value))
	case linux.EV_CODE_Y:
		this.rel.Y = float32(int32(evt.Value))
	}
}

func (this *State) DecodeMsc(evt linux.EVEvent) {
	switch evt.Code {
	case linux.EV_CODE_SCANCODE:
		this.scanCode = evt.Value
	}
}

func (this *State) DecodeSyn(evt linux.EVEvent) {
	ts := time.Duration(evt.Second)*time.Second + time.Duration(evt.Microsecond)*time.Microsecond
	switch {
	case this.rel.Equals(gopi.ZeroPoint) == false:
		this.position.X += this.rel.X
		this.position.Y += this.rel.Y
		return NewRelPositionEvent(this.rel, this.position, ts)
	case this.keyAction == gopi.KEYACTION_KEY_DOWN || this.keyAction == gopi.KEYACTION_KEY_UP || this.keyAction == gopi.KEYACTION_KEY_REPEAT:
		return NewKeyEvent(this.keyAction, this.keyCode, this.keyState, this.scanCode, ts)
	default:
		return nil
	}
}
