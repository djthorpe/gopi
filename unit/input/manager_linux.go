// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package input

import (
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type inputmanager struct {
	filepoll gopi.FilePoll
	bus      gopi.Bus
	devices  map[uintptr]gopi.InputDevice

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *inputmanager) Init(config InputManager) error {

	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else if config.FilePoll == nil {
		return gopi.ErrBadParameter.WithPrefix("filepoll")
	} else {
		this.devices = make(map[uintptr]gopi.InputDevice)
		this.filepoll = config.FilePoll
		this.bus = config.Bus
	}

	// Success
	return nil
}

func (this *inputmanager) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close devices
	for _, device := range this.devices {
		if err := this.closeDeviceEx(device); err != nil {
			return err
		}
	}

	// Release resources
	this.devices = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.InputManager

func (this *inputmanager) OpenDevicesByNameType(name string, flags gopi.InputDeviceType, exclusive bool) ([]gopi.InputDevice, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Obtain all devices
	if ids, err := linux.EVDevices(); err != nil {
		return nil, err
	} else {
		// Cycle through devices, ignoring ones which don't match
		devices := make([]gopi.InputDevice, 0, len(ids))
		for _, deviceId := range ids {
			if device := this.deviceById(deviceId); device != nil {
				if device.Matches(name, flags) {
					devices = append(devices, device)
				}
			} else if device, err := gopi.New(Device{deviceId, false, this.bus}, this.Log.Clone(Device{}.Name())); err != nil {
				this.Log.Warn(err)
			} else if matches := device.(gopi.InputDevice).Matches(name, flags); matches {
				device.Close()
				if device, err := this.openDeviceEx(deviceId, exclusive); err != nil {
					return nil, err
				} else {
					devices = append(devices, device)
				}
			} else {
				device.Close()
			}
		}
		return devices, nil
	}
}

func (this *inputmanager) OpenDevice(id uint, exclusive bool) (gopi.InputDevice, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.openDeviceEx(id, exclusive)
}

func (this *inputmanager) openDeviceEx(id uint, exclusive bool) (gopi.InputDevice, error) {
	if device := this.deviceById(id); device != nil {
		return nil, gopi.ErrBadParameter.WithPrefix("id")
	} else if device_, err := gopi.New(Device{id, exclusive, this.bus}, this.Log.Clone(Device{}.Name())); err != nil {
		return nil, err
	} else if device, ok := device_.(gopi.InputDevice); ok == false {
		return nil, gopi.ErrInternalAppError
	} else if err := this.filepoll.Watch(device.Fd(), gopi.FILEPOLL_FLAG_READ, this.watch); err != nil {
		device.(gopi.Unit).Close()
		return nil, err
	} else {
		this.devices[device.Fd()] = device
		return device, nil
	}
}

func (this *inputmanager) CloseDevice(device gopi.InputDevice) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.closeDeviceEx(device)
}

func (this *inputmanager) closeDeviceEx(device gopi.InputDevice) error {
	if device == nil || device.Fd() == 0 {
		return gopi.ErrBadParameter.WithPrefix("device")
	} else if _, exists := this.devices[device.Fd()]; exists == false {
		return gopi.ErrNotFound.WithPrefix("device")
	} else if err := this.filepoll.Unwatch(device.Fd()); err != nil {
		return err
	} else {
		err := device.(gopi.Unit).Close()
		delete(this.devices, device.Fd())
		return err
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *inputmanager) deviceById(id uint) gopi.InputDevice {
	for _, device := range this.devices {
		if device.Id() == id {
			return device
		}
	}
	// Not found, return nil
	return nil
}

func (this *inputmanager) watch(fd uintptr, flags gopi.FilePollFlags) {
	if device_, exists := this.devices[fd]; exists == false {
		return
	} else if flags&gopi.FILEPOLL_FLAG_READ == gopi.FILEPOLL_FLAG_READ {
		device_.(*device).read(this)
	}
}
