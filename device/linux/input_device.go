/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"os"
	"fmt"
	"time"
)

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC TYPES

type InputDevice struct {
	// The device path
	Path string
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE TYPES

// Represents an input device such as a keyboard, mouse or touchscreen
type evDevice struct {

	// The Name of the input device
	name string

	// The Physical ID of the input device
	phys string

	// logging object
	log *util.LoggerDevice

	// The type of device, or NONE if not known
	device_type hw.InputDeviceType

	// The bus which the device is attached to, or NONE if not known
	bus hw.InputDeviceBus

	// Product and version
	product uint16
	vendor uint16
	version uint16

	// Capabilities
	capabilities []evType

	// Handle to the device
	handle *os.File

	// Polling
	poll *evPoll
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new InputDevice object or return error
func (config InputDevice) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	var err error

	log.Debug("<linux.InputDevice>Open path=%v",config.Path)
	this := new(evDevice)
	this.log = log
	if this.handle, err = os.OpenFile(config.Path, os.O_RDWR, 0); err != nil {
		return nil, err
	}

	// Get name of device
	if this.name, err = evGetName(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	}

	// Get phys of device (physical connection ID)
	if this.phys, err = evGetPhys(this.handle); err != nil {
		this.handle.Close()
		return nil, err
	}

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
	case this.evSupportsEventType(EV_KEY,EV_LED,EV_REP):
		this.device_type = hw.INPUT_TYPE_KEYBOARD
	case this.evSupportsEventType(EV_KEY,EV_REL):
		this.device_type = hw.INPUT_TYPE_MOUSE
	case this.evSupportsEventType(EV_KEY,EV_ABS,EV_MSC):
		this.device_type = hw.INPUT_TYPE_JOYSTICK
	case this.evSupportsEventType(EV_KEY,EV_ABS):
		this.device_type = hw.INPUT_TYPE_TOUCHSCREEN
	}

	return this, nil
}

// Close InputDevice
func (this *evDevice) Close() error {
	this.log.Debug("<linux.InputDevice>Close device=%v",this)

	// shutdown polling
	if this.poll != nil {
		if err := this.poll.Close(); err != nil {
			return err
		}
		this.poll = nil
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

// Return true if the device matches an alias and/or a device type and/or bus
func (this *evDevice) Matches(alias string,device_type hw.InputDeviceType,device_bus hw.InputDeviceBus) bool {
	// Check the device type. We use NONE or ANY to match any device
	// type. The input argument can be OR'd in order to match more than one
	// device type.
	if device_type == hw.INPUT_TYPE_NONE {
		device_type = hw.INPUT_TYPE_ANY
	}
	if device_type != hw.INPUT_TYPE_ANY {
		if this.device_type & device_type == 0 {
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
	if alias == this.name {
		return true
	}
	if alias == this.phys {
		return true
	}
	return false
}

// Starts watching for events from the device
func (this *evDevice) Watch(callback hw.InputEventCallback) error {
	var err error

	if this.poll != nil {
		if err := this.poll.Close(); err != nil {
			return err
		}
	}
	if this.poll, err = evNewPoll(this.handle); err != nil {
		return err
	}

	go this.evWatch()

	return nil
}

// Background thread for polling
func (this *evDevice) evWatch() error {
	this.log.Debug2("<linux.InputDevice>WatchOpen device=%v",this)
	err := this.poll.evPoll(func (event *evEvent) {
		this.log.Debug2("Event = %v",event)
	})
	this.log.Debug2("<linux.InputDevice>WatchClose error=%v",err)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *evDevice) String() string {
	return fmt.Sprintf("<linux.InputDevice>{ name=\"%s\" phys=%v type=%v bus=%v product=0x%04X vendor=0x%04X version=0x%04X capabilities=%v fd=%v }", this.name, this.phys, this.device_type, this.bus, this.product, this.vendor, this.version, this.capabilities, this.handle.Fd())
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

