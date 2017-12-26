/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package mock

import (
	// Frameworks
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Device struct {
	Name     string
	Type     gopi.InputDeviceType
	Bus      gopi.InputDeviceBus
	Position gopi.Point
}

type device struct {
	log         gopi.Logger
	subscribers []chan gopi.Event
	name        string
	typ         gopi.InputDeviceType
	bus         gopi.InputDeviceBus
	position    gopi.Point
	keystate    gopi.KeyState
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Device) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.mock.InputDevice.Open{ name=%v type=%v bus=%v position=%v }", config.Name, config.Type, config.Bus, config.Position)

	if config.Name == "" {
		return nil, gopi.ErrBadParameter
	}
	if config.Type == gopi.INPUT_TYPE_ANY || config.Type == gopi.INPUT_TYPE_NONE {
		return nil, gopi.ErrBadParameter
	}
	if config.Bus == gopi.INPUT_BUS_ANY {
		return nil, gopi.ErrBadParameter
	}

	this := new(device)
	this.log = logger
	this.name = config.Name
	this.typ = config.Type
	this.bus = config.Bus
	this.position = config.Position
	this.keystate = gopi.KEYSTATE_NONE
	this.subscribers = make([]chan gopi.Event, 0)

	// Success
	return this, nil
}

// Close
func (this *device) Close() error {
	this.log.Debug("sys.mock.InputDevice.Close{ }")
	if this.subscribers == nil {
		this.log.Warn("sys.mock.InputDevice.Close: Called Close() more than once")
		return nil
	}
	this.subscribers = nil
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INPUT DEVICE INTERFACE

// Name of the input device
func (this *device) Name() string {
	return this.name
}

// Type of device
func (this *device) Type() gopi.InputDeviceType {
	return this.typ
}

// Bus interface
func (this *device) Bus() gopi.InputDeviceBus {
	return this.bus
}

// Position of cursor (for mouse, joystick and touchscreen devices)
func (this *device) Position() gopi.Point {
	return this.position
}

// Set absolute current cursor position
func (this *device) SetPosition(position gopi.Point) {
	this.position = position
}

// Get key states (caps lock, shift, scroll lock, num lock, etc)
func (this *device) KeyState() gopi.KeyState {
	return this.keystate
}

// Set key state (or states) to on or off. Will return error
// for key states which are not modifiable
func (this *device) SetKeyState(flags gopi.KeyState, state bool) error {
	// iterate through the key states
	for f := gopi.KEYSTATE_MIN; f <= gopi.KEYSTATE_MAX; f <<= 1 {
		fmt.Println(f)
	}
	return nil
}

// Returns true if device matches all conditions
func (this *device) Matches(alias string, device_type gopi.InputDeviceType, device_bus gopi.InputDeviceBus) bool {
	if alias != "" && alias != this.name {
		// If alias is set, it needs to match the name
		return false
	}
	if device_type != gopi.INPUT_TYPE_ANY && device_type != this.typ {
		// Match device type
		return false
	}
	if device_bus != gopi.INPUT_BUS_ANY && device_bus != this.bus {
		// Match device bus
		return false
	}
	// All alias, type and bus match
	return true
}

////////////////////////////////////////////////////////////////////////////////
// PUBLISHER INTERFACE

func (this *device) Subscribe() chan gopi.Event {
	subscriber := make(chan gopi.Event)
	this.subscribers = append(this.subscribers, subscriber)
	return subscriber
}

// Unsubscribe from events emitted
func (this *device) Unsubscribe(subscriber chan gopi.Event) {
	for i, s := range this.subscribers {
		if subscriber == s {
			this.subscribers[i] = nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	if this.typ&(gopi.INPUT_TYPE_JOYSTICK|gopi.INPUT_TYPE_MOUSE|gopi.INPUT_TYPE_TOUCHSCREEN) != 0 {
		return fmt.Sprintf("sys.mock.InputDevice{ name=%v type=%v bus=%v keystate=%v position=%v }", this.name, this.typ, this.bus, this.keystate, this.position)
	} else {
		return fmt.Sprintf("sys.mock.InputDevice{ name=%v type=%v bus=%v keystate=%v }", this.name, this.typ, this.bus, this.keystate)
	}
}
