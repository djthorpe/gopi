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

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/sys/hw/linux"
	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Input manager
type InputManager struct {
	// Filepoller
	FilePoll linux.FilePollInterface

	// Whether to try and get exclusivity when opening devices
	Exclusive bool
}

// Driver of multiple input devices
type manager struct {
	log      gopi.Logger
	filepoll linux.FilePollInterface
	pubsub   *util.PubSub

	// Whether to try and get exclusivity when opening devices
	exclusive bool

	// List of open devices
	devices []gopi.InputDevice
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config InputManager) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.input.linux.InputManager.Open>{ exclusive=%v }", config.Exclusive)

	// create new input device manager
	this := new(manager)

	if config.FilePoll == nil {
		return nil, gopi.ErrBadParameter
	}

	this.exclusive = config.Exclusive
	this.log = log
	this.filepoll = config.FilePoll
	this.pubsub = util.NewPubSub(0)
	this.devices = make([]gopi.InputDevice, 0)

	// success
	return this, nil
}

// Close Input driver
func (this *manager) Close() error {
	this.log.Debug("<sys.input.linux.InputManager.Close>{ }")

	this.pubsub.Close()

	// Empty
	this.filepoll = nil
	this.pubsub = nil
	this.devices = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	return fmt.Sprintf("<sys.input.linux.InputManager>{ exclusive=%v }", this.exclusive)
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE DEVICES

// OpenDevicesByName can be called often in order to open any newly plugged in
// devices. It will only return any newly opened devices.
func (this *manager) OpenDevicesByName(alias string, flags gopi.InputDeviceType, bus gopi.InputDeviceBus) ([]gopi.InputDevice, error) {
	this.log.Debug2("<sys.input.linux.InputManager.OpenDevicesByName>{ alias=%v flags=%v bus=%v }", alias, flags, bus)

	new_devices := make([]gopi.InputDevice, 0)
	opened_devices := make([]gopi.InputDevice, 0)

	// Discover devices using evFind and add any new ones to the new_devices
	// array, they are left in an opened state
	evFind(func(path string) {
		// Don't consider devices which are already opened
		if this.deviceByPath(path) == nil {
			if input_device, err := gopi.Open(InputDevice{Path: path, Exclusive: this.exclusive}, this.log); err != nil {
				this.log.Warn("OpenDevicesByName: %v: %v", path, err)
			} else {
				this.log.Debug2("OpenDevicesByName: Adding device %v", input_device)
				new_devices = append(new_devices, input_device.(gopi.InputDevice))
			}
		}
	})

	// Now check devices against filters
	for _, device := range new_devices {
		if device.Matches(alias, flags, bus) {
			opened_devices = append(opened_devices, device)
		} else if err := device.Close(); err != nil {
			this.log.Warn("OpenDevicesByName: %v", err)
		}
	}

	// Return newly opened devices
	return opened_devices, nil
}

/*
	     new_devices := make([]*evDevice, 0)

	   	// Discover devices using evFind and add any new ones to the new_devices
	   	// array, they are left in an opened state
	   	evFind(func(path string) {
	   		device := this.evDeviceByPath(path)
	   		if device == nil {
	   			// we open the device here
	   			var err gopi.Error
	   			if gopi_device, ok := gopi.Open2(InputDevice{Path: path, Exclusive: this.exclusive}, this.log, &err).(*evDevice); !ok {
	   				this.log.Warn("<linux.Input>OpenDevicesByName path=%v Error: %v", path, err)
	   				return
	   			} else {
	   				new_devices = append(new_devices, gopi_device)
	   			}
	   		}
	   	})

	   	// Now we check the new devices to see if they match the stated criteria
	   	// and close the device if not
	   	for _, device := range new_devices {

	   		// Check if device matches criteria. If not, then close it
	   		if device.Matches(alias, flags, bus) == false {
	   			if err := device.Close(); err != nil {
	   				this.log.Warn("<linux.Input>OpenDevicesByName Error: %v", err)
	   			}
	   			continue
	   		}

	   		// We have matched devices here, poll them
	   		if err := this.poll.Add(device.GetFd(), POLL_MODE_READ); err != nil {
	   			this.log.Warn("<linux.Input>OpenDevicesByName Error: %v", err)
	   			device.Close()
	   			continue
	   		}

	   		// cleanup obtain the file descriptor
	   		this.devices[device.GetFd()] = device
	   		opened_devices = append(opened_devices, device)
	   	}

	     return opened_devices, nil
	return nil, gopi.ErrNotImplemented
}
*/

func (this *manager) CloseDevice(device gopi.InputDevice) error {
	this.log.Debug2("<sys.input.linux.InputManager.CloseDevice>{ device=%v }", device)
	return gopi.ErrNotImplemented
}

/*
	// Remove device from poll
	if linux_device, ok := device.(*device); ok {
		if err := this.filepoll.Unwatch(linux_device.GetFd()); err != nil {
			return err
		}
	}
		// Remove device from list of devices
		delete(this.devices, int(linux_device.GetFd()))
	// Close device
	if err := device.Close(); err != nil {
		return err
	}

	// Success
	return nil
}
*/

////////////////////////////////////////////////////////////////////////////////
// PUBLISH AND SUBSCRIBE TO INPUT EVENTS

// Subscribe to events emitted. Returns unique subscriber
// identifier and channel on which events are emitted
func (this *manager) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

// Unsubscribe from events emitted
func (this *manager) Unsubscribe(subscriber <-chan gopi.Event) {
	this.pubsub.Unsubscribe(subscriber)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *manager) deviceByPath(path string) gopi.InputDevice {
	for _, d := range this.devices {
		if linux_device, is_linux := d.(*device); is_linux {
			if linux_device.path == path {
				return d
			}
		}
	}
	return nil
}
