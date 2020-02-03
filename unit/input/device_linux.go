// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

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
	Bus       gopi.Bus
}

type device struct {
	id         uint
	dev        *os.File
	name       string
	cap        []linux.EVType
	exclusive  bool
	deviceType gopi.InputDeviceType
	bus        gopi.Bus

	State
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
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.bus = config.Bus
	}

	if dev, err := linux.EVOpenDevice(config.Id); err != nil {
		return err
	} else {
		this.id = config.Id
		this.dev = dev
	}

	if name, err := linux.EVGetName(this.dev.Fd()); err != nil {
		this.dev.Close()
		this.dev = nil
		return err
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

	// Reset state
	this.State.Reset()
	this.State.log = this.Log

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
	this.cap = nil
	this.bus = nil

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
			" fd=" + fmt.Sprint(this.Fd()) +
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

func (this *device) KeyState() gopi.KeyState {
	if this.dev == nil {
		return gopi.KEYSTATE_NONE
	} else {
		return this.State.KeyState()
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
	return this.State.Position()
}

func (this *device) SetPosition(point gopi.Point) {
	this.State.SetPosition(point)
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

func (this *device) read(source gopi.InputManager) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var evt linux.EVEvent
	if err := binary.Read(this.dev, binary.LittleEndian, &evt); err != nil {
		this.Log.Error(err)
		return
	}

	if evt := this.State.Decode(evt); evt != nil {
		evt.device = this
		evt.source = source
		this.bus.Emit(evt)
	}
}

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
