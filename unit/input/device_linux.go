// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"encoding/binary"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Device struct {
	Id        uint
	Exclusive bool
}

type device struct {
	id         uint
	dev        *os.File
	name       string
	cap        []linux.EVType
	exclusive  bool
	deviceType gopi.InputDeviceType

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

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Device) Name() string { return "gopi.InputDevice" }

func (config Device) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(device)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// INIT gopi.InputDevice

func (this *device) Init(config Device) error {
	if dev, err := linux.EVOpenDevice(config.Id); err != nil {
		return err
	} else {
		this.id = config.Id
		this.dev = dev			
	}
	
	if name, err := linux.EVGetName(this.dev.Fd()); err != nil {
		this.dev.Close()
		this.dev = nil
		return  err
	} else {
		this.name = strings.TrimSpace(name)
	}

	if cap, err := linux.EVGetSupportedEventTypes(this.dev.Fd()); err != nil {
		this.dev.Close()
		this.dev = nil
		return err
	} else {
		this.cap = cap
	}

	if config.Exclusive {
		if err := linux.EVSetGrabState(this.dev.Fd(), true); err != nil {
			this.dev.Close()
			this.dev = nil
			return err
		} else {
			this.exclusive = config.Exclusive
		}
	}

	// Return success
	return nil
}

func (this *device) Close() error {
	// Don't call close if already called
	if this.dev == nil {
		return nil
	}

	// Lock closing
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Remove exclusive access
	if this.exclusive {
		if err := linux.EVSetGrabState(this.Fd(), false); err != nil {
			return err
		}
	}

	// Close device
	if err := this.dev.Close(); err != nil {
		return err
	}

	// Release resources
	this.dev = nil

	// Success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY gopi.InputDevice

func (this *device) String() string {
	if this.dev == nil {
		return "<gopi.inputdevice>"
	} else {
		return "<gopi.inputdevice id=" + fmt.Sprint(this.id) +
		" fd=" + fmt.Sprint(this.Fd())  +
		" name=" + strconv.Quote(this.name) +
			" type=" + fmt.Sprint(this.Type()) +
			" capabilities=" + fmt.Sprint(this.cap) +
			">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES gopi.InputDevice

func (this *device) Name() string {
	return this.name
}

func (this *device) Id() uint {
	return this.id
}

func (this *device) Fd() uintptr {
	if this.dev == nil {
		return 0
	} else {
		return this.dev.Fd()
	}
}


func (this *device) State() gopi.KeyState {
	if this.dev == nil {
		return gopi.KEYSTATE_NONE
	} else {
		return this.keyState
	}
}

func (this *device) Type() gopi.InputDeviceType {
	// Return cached type
	if this.deviceType != 0 {
		return this.deviceType
	}

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Determine device type. We don't know if joysticks are
	// currently supported, however, so will need to find a
	// joystick tester later
	switch {
	case evSupportsEventType(this.cap, linux.EV_KEY, linux.EV_LED, linux.EV_REP):
		this.deviceType = gopi.INPUT_TYPE_KEYBOARD
	case evSupportsEventType(this.cap, linux.EV_KEY, linux.EV_REL):
		this.deviceType = gopi.INPUT_TYPE_MOUSE
	case evSupportsEventType(this.cap, linux.EV_KEY, linux.EV_ABS, linux.EV_MSC):
		this.deviceType = gopi.INPUT_TYPE_JOYSTICK
	case evSupportsEventType(this.cap, linux.EV_KEY, linux.EV_ABS):
		this.deviceType = gopi.INPUT_TYPE_TOUCHSCREEN
	}

	// Return device type
	return this.deviceType
}

func (this *device) Position() gopi.Point {
	return this.position
}

func (this *device) SetPosition(point gopi.Point) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.position = point
}

////////////////////////////////////////////////////////////////////////////////
// MATCH SPECIFIC NAME/FLAGS

func (this *device) Matches(name string, flags gopi.InputDeviceType) bool {
	if this.dev == nil {
		return false
	}
	if name == "" && (flags == gopi.INPUT_TYPE_ANY || flags == gopi.INPUT_TYPE_NONE) {
		return true
	}
	if name == this.name && (flags == gopi.INPUT_TYPE_ANY || flags == gopi.INPUT_TYPE_NONE) {
		return true
	}
	if name == "" && (flags&this.Type() > 0) {
		return true
	}
	// No match
	return false
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *device) read(source gopi.Unit) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var evt linux.EVEvent
	if err := binary.Read(this.dev, binary.LittleEndian, &evt); err != nil {
		this.Log.Error(err)
	} else {
		this.Log.Debug(evt)
	}
}

/*	
		switch evt.Type {
		case linux.EV_KEY:
			evDecodeKey(evt, device)
		case linux.EV_SYN:
			evDecodeSyn(evt, device)
		case linux.EV_ABS:
		case linux.EV_REL:
			evDecodeRel(evt, device)
		case linux.EV_MSC:
			evDecodeMisc(evt, device)
		}
	}
*/


////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func evSupportsEventType(cap []linux.EVType, req ...linux.EVType) bool {
	// Create a map of capabilities to easily lookup
	capmap := make(map[linux.EVType]bool, len(cap))
	for _, r := range cap {
		capmap[r] = true
	}
	// Now iterate through requested capabilities, will return false
	// it any capability is not present
	for _, r := range req {
		if _, exists := capmap[r]; exists == false {
			return false
		}
	}
	return true
}

func evDecodeKey(evt linux.EVEvent, device *device) {
	// Interpret key code and key action (up, down and repeat)
	code, action := gopi.KeyCode(evt.Code), gopi.KeyAction(evt.Value)

	// Alter key state if a modified key was pressed and also
	// handle sticky state keys CAPS, NUM and SCROLL locks
	state := gopi.KEYSTATE_NONE
	switch code {
	case gopi.KEYCODE_CAPSLOCK:
		if action == gopi.KEYACTION_KEY_DOWN {
			device.keyState ^= gopi.KEYSTATE_CAPSLOCK
			// TODO: Change LED's
		}
	case gopi.KEYCODE_NUMLOCK:
		if action == gopi.KEYACTION_KEY_DOWN {
			device.keyState ^= gopi.KEYSTATE_NUMLOCK
			// TODO: Change LED's
		}
	case gopi.KEYCODE_SCROLLLOCK:
		if action == gopi.KEYACTION_KEY_DOWN {
			device.keyState ^= gopi.KEYSTATE_SCROLLLOCK
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
	device.keyCode = code
	device.keyAction = action

	// Set device state from key action
	if state != gopi.KEYSTATE_NONE {
		if action == gopi.KEYACTION_KEY_DOWN || action == gopi.KEYACTION_KEY_REPEAT {
			device.keyState |= state
		} else if action == gopi.KEYACTION_KEY_UP {
			device.keyState ^= state
		}
	}
}

func evDecodeAbs(evt linux.EVEvent, device *device) {
	switch evt.Code {
	case linux.EV_CODE_X:
		device.position.X = float32(int32(evt.Value))
	case linux.EV_CODE_Y:
		device.position.Y = float32(int32(evt.Value))
	case linux.EV_CODE_SLOT:
		device.slot = evt.Value
	default:
		fmt.Println("Ignoring", evt)
	}
}

func evDecodeRel(evt linux.EVEvent, device *device) {
	switch evt.Code {
	case linux.EV_CODE_X:
		device.rel.X = float32(int32(evt.Value))
	case linux.EV_CODE_Y:
		device.rel.Y = float32(int32(evt.Value))
	}
}

func evDecodeMisc(evt linux.EVEvent, device *device) {
	switch evt.Code {
	case linux.EV_CODE_SCANCODE:
		device.scanCode = evt.Value
	}
}

func evDecodeSyn(evt linux.EVEvent, device *device) {
	ts := time.Duration(evt.Second)*time.Second + time.Duration(evt.Microsecond)*time.Microsecond
	switch {
	case device.rel.Equals(gopi.ZeroPoint) == false:
		device.position.X += device.rel.X
		device.position.Y += device.rel.Y
		// Emit RELATIVE POSITION
		fmt.Printf("<gopi.InputEvent rel=%v abs=%v ts=%v>\n", device.rel, device.position, ts)
	case device.keyCode != gopi.KEYCODE_NONE:
		// Emit KEY EVENT
		fmt.Printf("<gopi.InputEvent action=%v code=%v state=%v scan_code=0x%08X ts=%v>\n", device.keyAction, device.keyCode, device.keyState, device.scanCode, ts)
	}
	// Reset key action and rel position
	device.keyCode = gopi.KEYCODE_NONE
	device.rel = gopi.ZeroPoint
}
