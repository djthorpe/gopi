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

	// Whether to try and get exclusivity when opening devices
	exclusive bool
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Pattern for finding event-driven input devices
	INPUT_PATH_DEVICES = "/sys/class/input/event*"
)

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

	// success
	return this, nil
}

// Close Input driver
func (this *manager) Close() error {
	this.log.Debug("<sys.input.linux.InputManager.Close>{ }")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	return fmt.Sprintf("<sys.input.linux.InputManager>{ exclusive=%v }", this.exclusive)
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE DEVICES

func (this *manager) OpenDevicesByName(alias string, flags gopi.InputDeviceType, bus gopi.InputDeviceBus) ([]gopi.InputDevice, error) {
	this.log.Debug2("<sys.input.linux.InputManager.OpenDevicesByName>{ alias=%v flags=%v bus=%v }", alias, flags, bus)

	/*
	     new_devices := make([]*evDevice, 0)
	   	opened_devices := make([]hw.InputDevice, 0)

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
	*/
	return nil, gopi.ErrNotImplemented
}

func (this *manager) CloseDevice(device gopi.InputDevice) error {
	this.log.Debug2("<sys.input.linux.InputManager.CloseDevice>{ device=%v }", device)
	return gopi.ErrNotImplemented
	/*
		// Remove device from poll
		linux_device, ok := device.(*evDevice)
		if ok == true {
			if err := this.poll.Remove(linux_device.GetFd(), POLL_MODE_READ); err != nil {
				return err
			}

			// Remove device from list
			delete(this.devices, int(linux_device.GetFd()))
		}

		// Close device
		if err := device.Close(); err != nil {
			return err
		}

		// Success
		return nil*/
}
