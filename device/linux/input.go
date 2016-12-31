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
	"path"
	"path/filepath"
)

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Input configuration
type Input struct{
	// Whether to try and get exclusivity when opening devices
	Exclusive bool
}

// Driver of multiple input devices
type InputDriver struct {
	// Whether to try and get exclusivity when opening devices
	exclusive bool

	// Logger
	log     *util.LoggerDevice

	// Map of input devices (keyed by their file descriptor)
	devices map[int]*evDevice  // input devices

	// Polling mechanism to check for events
	poll    *PollDriver
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Pattern for finding event-driven input devices
	INPUT_PATH_DEVICES = "/sys/class/input/event*"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new Input object, returns error if not possible
func (config Input) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	var err error

	log.Debug("<linux.Input>Open")

	// create new GPIO driver
	this := new(InputDriver)

	// Set logging & devices
	this.exclusive = config.Exclusive
	this.log = log
	this.devices = make(map[int]*evDevice, 0)

	// Create polling mechanism
	if this.poll, err = NewPollDriver(this.log); err != nil {
		return nil, err
	}

	// success
	return this, nil
}

// Close Input driver
func (this *InputDriver) Close() error {
	this.log.Debug("<linux.Input>Close")

	if err := this.poll.Close(); err != nil {
		return err
	}

	for _, device := range this.devices {
		if device == nil {
			continue
		}
		if err := device.Close(); err != nil {
			return err
		}
	}

	return nil
}


////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *InputDriver) String() string {
	return fmt.Sprintf("<linux.Input>{ exclusive=%v open_devices=%v }",this.exclusive,this.devices)
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE DEVICES

func (this *InputDriver) OpenDevicesByName(alias string, flags hw.InputDeviceType, bus hw.InputDeviceBus) ([]hw.InputDevice, error) {

	new_devices := make([]*evDevice, 0)
	opened_devices := make([]hw.InputDevice, 0)

	// Discover devices using evFind and add any new ones to the new_devices
	// array, they are left in an opened state
	evFind(func(path string) {
		device := this.evDeviceByPath(path)
		if device == nil {
			// we open the device here
			gopi_device, err := gopi.Open(InputDevice{Path: path, Exclusive: this.exclusive}, this.log)
			if err != nil {
				this.log.Warn("<linux.Input>OpenDevicesByName path=%v Error: %v", path, err)
				return
			}
			new_devices = append(new_devices, gopi_device.(*evDevice))
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
}

func (this *InputDriver) CloseDevice(device hw.InputDevice) error {
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
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Find all input devices
func evFind(callback func(string)) error {
	files, err := filepath.Glob(INPUT_PATH_DEVICES)
	if err != nil {
		return err
	}
	for _, file := range files {
		callback(path.Clean(path.Join("/", "dev", "input", path.Base(file))))
	}
	return nil
}

// Return opened device by path
func (this *InputDriver) evDeviceByPath(path string) *evDevice {
	for _, device := range this.devices {
		if path == device.GetPath() {
			return device
		}
	}
	return nil
}
