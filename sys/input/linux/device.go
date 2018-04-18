// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/sys/hw/linux"
	"github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Input device
type InputDevice struct {
	// Filepoller
	FilePoll linux.FilePollInterface

	// Path to device
	Path string

	// Whether to try and get exclusivity
	Exclusive bool
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE TYPES

// Represents an input device such as a keyboard, mouse or touchscreen
type device struct {
	log       gopi.Logger
	path      string
	filepoll  linux.FilePollInterface
	pubsub    *event.PubSub
	exclusive bool

	// Handle to the device
	handle *os.File

	// The Name of the input device
	name string

	// The Physical ID of the input device
	phys string

	// Unique Identifier
	uniq string

	// The type of device, or NONE if not known
	device_type gopi.InputDeviceType

	// The bus which the device is attached to, or NONE if not known
	bus gopi.InputDeviceBus

	// Product and version
	product uint16
	vendor  uint16
	version uint16

	// Capabilities
	capabilities []evType

	// Positions for mice, joystick and touchscreens
	position      gopi.Point
	rel_position  gopi.Point
	last_position gopi.Point

	// Key presses
	key_code   evKeyCode
	key_action evKeyAction
	scan_code  uint32

	// Multi-touch support
	slot  uint32
	slots []slot

	/*
		// the current key state, which is a set of OR'd flags
		state gopi.KeyState

	*/
}

// Represents multi-touch slot information
type slot struct {
	id       int16
	position gopi.Point
	active   bool
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new InputDevice object or return error
func (config InputDevice) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.input.linux.InputDevice.Open>{ path=%v exclusive=%v }", config.Path, config.Exclusive)

	// Check incoming configuration parameters
	if config.FilePoll == nil {
		return nil, gopi.ErrBadParameter
	}
	if config.Path == "" {
		return nil, gopi.ErrBadParameter
	}

	this := new(device)
	this.log = log
	this.path = config.Path
	this.exclusive = config.Exclusive
	this.filepoll = config.FilePoll

	// Open the event stream for reading and writing
	if handle, err := os.OpenFile(config.Path, os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		this.handle = handle
	}

	// Get name of device
	if name, err := evGetName(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	} else {
		this.name = name
	}

	// Get phys & uniq of device. Ignore errors here,
	// since it seems this isn't reported by touchscreen
	this.phys, _ = evGetPhys(this.handle)
	this.uniq, _ = evGetUniq(this.handle)

	// Get information about the device
	if bus, vendor, product, version, err := evGetInfo(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	} else {
		this.bus = gopi.InputDeviceBus(bus)
		this.vendor = vendor
		this.product = product
		this.version = version
	}

	// Get capabilities
	if capabilities, err := evGetSupportedEventTypes(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	} else {
		this.capabilities = capabilities
	}

	// Determine device type. We don't know if joysticks are
	// currently supported, however, so will need to find a
	// joystick tester later
	switch {
	case evSupportsEventType(this.capabilities, EV_KEY, EV_LED, EV_REP):
		this.device_type = gopi.INPUT_TYPE_KEYBOARD
	case evSupportsEventType(this.capabilities, EV_KEY, EV_REL):
		this.device_type = gopi.INPUT_TYPE_MOUSE
	case evSupportsEventType(this.capabilities, EV_KEY, EV_ABS, EV_MSC):
		this.device_type = gopi.INPUT_TYPE_JOYSTICK
	case evSupportsEventType(this.capabilities, EV_KEY, EV_ABS):
		this.device_type = gopi.INPUT_TYPE_TOUCHSCREEN
	}

	// Start watching
	if err := this.filepoll.Watch(this.handle, linux.FILEPOLL_MODE_READ, this.evReceive); err != nil {
		this.handle.Close()
		return nil, err
	}

	// Obtain exclusive use of device
	if this.exclusive {
		if err := evSetGrabState(this.handle, true); err != nil {
			this.handle.Close()
			return nil, err
		}
	}

	// Set multi-touch slot array to track slots
	this.slot = 0
	this.slots = make([]slot, INPUT_MAX_MULTITOUCH_SLOTS)

	// PubSub
	this.pubsub = event.NewPubSub(0)

	// Success
	return this, nil
}

// Close InputDevice
func (this *device) Close() error {
	this.log.Debug("<sys.input.linux.InputDevice.Close>{ }")

	// remove exclusive access
	if this.exclusive {
		if err := evSetGrabState(this.handle, false); err != nil {
			this.log.Warn("<linux.InputDevice>Close Error: %v", err)
		}
		this.exclusive = false
	}

	// Unwatch device
	if err := this.filepoll.Unwatch(this.handle); err != nil {
		this.log.Warn("Unwatch: %v", err)
	}

	// close file handle
	if err := this.handle.Close(); err != nil {
		return err
	} else {
		this.handle = nil
	}

	// Close subscriber channels
	this.pubsub.Close()

	// Blank out
	this.filepoll = nil
	this.pubsub = nil
	this.handle = nil

	// return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MATCH DEVICE

// Return true if the device matches an alias, type and bus
func (this *device) Matches(alias string, flags gopi.InputDeviceType, bus gopi.InputDeviceBus) bool {
	this.log.Debug2("<sys.input.linux.InputDevice.Matches>{ alias=%v flags=%v bus=%v }", alias, flags, bus)
	// Check the device type. We use NONE or ANY to match any device
	// type. The input argument can be OR'd in order to match more than one
	// device type.
	if flags == gopi.INPUT_TYPE_NONE {
		flags = gopi.INPUT_TYPE_ANY
	}
	if flags != gopi.INPUT_TYPE_ANY {
		if this.device_type&flags == 0 {
			return false
		}
	}
	// Check device bus. Only one type of bus can
	// be selected at any one time, or NONE or ANY
	// will select any bus
	if bus == gopi.INPUT_BUS_NONE {
		bus = gopi.INPUT_BUS_ANY
	}
	if bus != gopi.INPUT_BUS_ANY {
		if this.bus != bus {
			return false
		}
	}
	// check alias against name, uniq or phys
	// if empty then return true
	if alias == "" {
		return true
	}
	if alias == this.uniq {
		return true
	}
	if alias == this.phys {
		return true
	}
	if alias == this.name {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE IMPLEMENTATION

// Return name of the device
func (this *device) Name() string {
	return this.name
}

// Return information on what we think the device is (mouse, keyboard, touchscreen)
func (this *device) Type() gopi.InputDeviceType {
	return this.device_type
}

// Return the bus we think the device is connected to
func (this *device) Bus() gopi.InputDeviceBus {
	return this.bus
}

// Return absolute cursor position
func (this *device) Position() gopi.Point {
	return this.position
}

////////////////////////////////////////////////////////////////////////////////
// PUBLISH AND SUBSCRIBE INTERFACE IMPLEMTATION

// Subscribe to events emitted. Returns unique subscriber
// identifier and channel on which events are emitted
func (this *device) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

// Unsubscribe from events emitted
func (this *device) Unsubscribe(subscriber <-chan gopi.Event) {
	this.pubsub.Unsubscribe(subscriber)
}

// Emit events
func (this *device) Emit(evt gopi.Event) {
	this.log.Info("Emit: %v", evt)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	return fmt.Sprintf("<sys.input.linux.InputDevice>{ name=\"%s\" phys=\"%v\" uniq=\"%v\" type=%v bus=%v position=%v product=0x%04X vendor=0x%04X version=0x%04X capabilities=%v exclusive=%v fd=%v path=%v }", this.name, this.phys, this.uniq, this.device_type, this.bus, this.position, this.product, this.vendor, this.version, this.capabilities, this.exclusive, this.handle.Fd(), this.path)
}
