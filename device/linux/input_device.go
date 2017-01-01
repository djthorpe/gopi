/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"fmt"
	"os"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	khronos "github.com/djthorpe/gopi/khronos"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC TYPES

type InputDevice struct {
	// The device path
	Path string

	// Whether to obtain exclusive access to device
	Exclusive bool
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE TYPES

// Represents an input device such as a keyboard, mouse or touchscreen
type evDevice struct {

	// The Path of the input device
	path string

	// The Name of the input device
	name string

	// The Physical ID of the input device
	phys string

	// Unique Identifier
	uniq string

	// logging object
	log *util.LoggerDevice

	// The type of device, or NONE if not known
	device_type hw.InputDeviceType

	// The bus which the device is attached to, or NONE if not known
	bus hw.InputDeviceBus

	// Product and version
	product uint16
	vendor  uint16
	version uint16

	// Capabilities
	capabilities []evType

	// Handle to the device
	handle *os.File

	// Positions, keys and states
	position      khronos.EGLPoint
	last_position khronos.EGLPoint
	rel_position  khronos.EGLPoint
	key_code      evKeyCode
	scan_code     uint32
	key_action    evKeyAction

	// exclusive access to device
	exclusive     bool

	// the current key state, which is a set of OR'd flags
	state         hw.InputKeyState

	// Multi-touch support
	slot          uint32
	slots         []evDeviceSlot
}

// Represents multi-touch slot information
type evDeviceSlot struct {
	id            int16
	position      khronos.EGLPoint
	active        bool
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	INPUT_MAX_MULTITOUCH_SLOTS uint32 = 32
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new InputDevice object or return error
func (config InputDevice) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	var err error

	log.Debug("<linux.InputDevice>Open path=%v", config.Path)

	this := new(evDevice)
	this.log = log
	this.path = config.Path

	// Open the event stream for reading and writing
	if this.handle, err = os.OpenFile(config.Path, os.O_RDWR, 0); err != nil {
		return nil, err
	}

	// Get name of device
	if this.name, err = evGetName(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	}

	// Get phys & uniq of device. Ignore errors here,
	// since it seems this isn't reported by touchscreen
	this.phys, err = evGetPhys(this.handle)
	this.uniq, err = evGetUniq(this.handle)

	// Get information about the device
	var bus uint16
	if bus, this.vendor, this.product, this.version, err = evGetInfo(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	}
	// Convert to the right type
	this.bus = hw.InputDeviceBus(bus)

	// Get capabilities
	if this.capabilities, err = evGetSupportedEventTypes(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	}

	// Determine device type. We don't know if joysticks are
	// currently supported, however, so will need to find a
	// joystick tester later
	switch {
	case this.evSupportsEventType(EV_KEY, EV_LED, EV_REP):
		this.device_type = hw.INPUT_TYPE_KEYBOARD
	case this.evSupportsEventType(EV_KEY, EV_REL):
		this.device_type = hw.INPUT_TYPE_MOUSE
	case this.evSupportsEventType(EV_KEY, EV_ABS, EV_MSC):
		this.device_type = hw.INPUT_TYPE_JOYSTICK
	case this.evSupportsEventType(EV_KEY, EV_ABS):
		this.device_type = hw.INPUT_TYPE_TOUCHSCREEN
	}

	// Obtain exclusive use of device
	this.exclusive = config.Exclusive
	if this.exclusive {
		if err := evSetGrabState(this.handle,true); err != nil {
			this.handle.Close()
			return nil, err
		}
	}

	// Set multi-touch slot array to track slots
	this.slot = 0
	this.slots = make([]evDeviceSlot,INPUT_MAX_MULTITOUCH_SLOTS)

	// Success
	return this, nil
}

// Close InputDevice
func (this *evDevice) Close() error {
	this.log.Debug("<linux.InputDevice>Close device=%v", this)

	// remove exclusive access
	if this.exclusive {
		if err := evSetGrabState(this.handle,false); err != nil {
			this.log.Warn("<linux.InputDevice>Close Error: %v",err)
		}
		this.exclusive = false
	}

	// close file handle
	err := this.handle.Close()
	if err != nil {
		return err
	}

	// return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE InputDevice IMPLEMENTATION

// Return name of the device
func (this *evDevice) GetName() string {
	return this.name
}

// Return information on what we think the device is (mouse, keyboard, touchscreen)
func (this *evDevice) GetType() hw.InputDeviceType {
	return this.device_type
}

// Return the bus we think the device is connected to
func (this *evDevice) GetBus() hw.InputDeviceBus {
	return this.bus
}

// Return the file descriptor
func (this *evDevice) GetFd() int {
	return int(this.handle.Fd())
}

// Return the path of the device
func (this *evDevice) GetPath() string {
	return this.path
}

// Return cursor position
func (this *evDevice) GetPosition() khronos.EGLPoint {
	return this.position
}

// Set cursor position
func (this *evDevice) SetPosition(position khronos.EGLPoint) {
	this.position = position
	this.last_position = position
}

// Return true if the device matches an alias and/or a device type and/or bus
func (this *evDevice) Matches(alias string, device_type hw.InputDeviceType, device_bus hw.InputDeviceBus) bool {
	// Check the device type. We use NONE or ANY to match any device
	// type. The input argument can be OR'd in order to match more than one
	// device type.
	if device_type == hw.INPUT_TYPE_NONE {
		device_type = hw.INPUT_TYPE_ANY
	}
	if device_type != hw.INPUT_TYPE_ANY {
		if this.device_type&device_type == 0 {
			return false
		}
	}
	// Check device bus. Only one type of bus can
	// be selected at any one time, or NONE or ANY
	// will select any bus
	if device_bus == hw.INPUT_BUS_NONE {
		device_bus = hw.INPUT_BUS_ANY
	}
	if device_bus != hw.INPUT_BUS_ANY {
		if this.bus != device_bus {
			return false
		}
	}
	// check alias, if empty then return true
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

func (this *evDevice) GetKeyState() hw.InputKeyState {
	current_state := this.state
	states, err := evGetLEDState(this.handle)
	if err != nil {
		this.log.Warn("<linux.InputDevice> Error: %v",err)
		return current_state
	}
	if states == nil || len(states) == 0 {
		return current_state
	}
	for _, state := range states {
		switch(state) {
			case EV_LED_NUML:
				current_state |= hw.INPUT_KEYSTATE_NUM
			case EV_LED_CAPSL:
				current_state |= hw.INPUT_KEYSTATE_CAPS
			case EV_LED_SCROLLL:
				current_state |= hw.INPUT_KEYSTATE_SCROLL
		}
	}
	return current_state
}

func (this *evDevice) SetKeyState(flags hw.InputKeyState,state bool) error {
	// Set the current state
	this.state = flags

	// Iterate through the states, and set the CAPS, SCROLL and NUM indicators
	// on the keyboard
	num_states := int(unsafe.Sizeof(flags) << 3)
	led := hw.InputKeyState(1)
	for j := 0; j <= num_states; j++ {
		if flags & 0x01 != 0x00 {
			switch led {
			case hw.INPUT_KEYSTATE_CAPS:
				if err := evSetLEDState(this.handle,EV_LED_CAPSL,state); err != nil {
					return err
				}
			case hw.INPUT_KEYSTATE_SCROLL:
				if err := evSetLEDState(this.handle,EV_LED_SCROLLL,state); err != nil {
					return err
				}
			case hw.INPUT_KEYSTATE_NUM:
				if err := evSetLEDState(this.handle,EV_LED_NUML,state); err != nil {
					return err
				}
			}
		}
		flags >>= 1
		led <<= 1
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *evDevice) String() string {
	return fmt.Sprintf("<linux.InputDevice>{ name=\"%s\" phys=%v uniq=%v type=%v bus=%v product=0x%04X vendor=0x%04X version=0x%04X capabilities=%v exclusive=%v fd=%v }", this.name, this.phys, this.uniq, this.device_type, this.bus, this.product, this.vendor, this.version, this.capabilities, this.exclusive, this.handle.Fd())
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS

// Returns true if all event types are supported by the device, else returns
// false
func (this *evDevice) evSupportsEventType(types ...evType) bool {
	count := 0
	for _, capability := range this.capabilities {
		for _, typ := range types {
			if typ == capability {
				count = count + 1
			}
		}
	}
	return (count == len(types))
}
