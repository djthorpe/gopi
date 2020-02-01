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
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type inputmanager struct {
	filepoll gopi.FilePoll
	devices  map[uintptr]*device

	base.Unit
	sync.Mutex
}

type device struct {
	bus        uint
	dev        *os.File
	name       string
	cap        []linux.EVType
	exclusive  bool
	deviceType gopi.InputDeviceType
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *inputmanager) Init(config InputManager) error {

	if config.FilePoll == nil {
		return gopi.ErrBadParameter.WithPrefix("filepoll")
	} else {
		this.devices = make(map[uintptr]*device)
		this.filepoll = config.FilePoll
	}

	// Success
	return nil
}

func (this *inputmanager) Close() error {
	for _, device := range this.devices {
		if err := this.CloseDevice(device); err != nil {
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
	if deviceIds, err := linux.EVDevices(); err != nil {
		return nil, err
	} else {
		// Cycle through devices, ignoring ones which don't match
		devices := make([]gopi.InputDevice, 0, len(deviceIds))
		for _, deviceId := range deviceIds {
			if device := this.deviceById(deviceId); device != nil {
				if device.matches(name, flags) {
					devices = append(devices, device)
				}
			} else if device, err := NewDevice(deviceId, false); err == nil {
				matches := device.matches(name, flags)
				device.Close()
				if matches == false {
					continue
				}
				if device, err := this.OpenDeviceEx(deviceId, exclusive); err != nil {
					return nil, err
				} else {
					devices = append(devices, device)
				}
			}
		}
		return devices, nil
	}
}

func (this *inputmanager) OpenDevice(bus uint, exclusive bool) (gopi.InputDevice, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.OpenDeviceEx(bus, exclusive)
}

func (this *inputmanager) OpenDeviceEx(bus uint, exclusive bool) (gopi.InputDevice, error) {
	if device, err := NewDevice(bus, exclusive); err != nil {
		return nil, err
	} else if err := this.filepoll.Watch(device.dev.Fd(), gopi.FILEPOLL_FLAG_READ, this.watch); err != nil {
		device.Close()
		return nil, err
	} else {
		this.devices[device.dev.Fd()] = device
		return device, nil
	}
}

func (this *inputmanager) CloseDevice(d gopi.InputDevice) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if d == nil {
		return gopi.ErrBadParameter.WithPrefix("device")
	} else if device_, ok := d.(*device); ok == false {
		return gopi.ErrBadParameter.WithPrefix("device")
	} else if device_.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("device")
	} else {
		fd := device_.dev.Fd()
		if _, exists := this.devices[fd]; exists == false {
			return gopi.ErrNotFound.WithPrefix("device")
		} else if err := this.filepoll.Unwatch(fd); err != nil {
			return err
		} else if err := device_.Close(); err != nil {
			return err
		} else {
			delete(this.devices, fd)
		}
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *inputmanager) deviceById(id uint) *device {
	for _, device := range this.devices {
		if device.bus == id {
			return device
		}
	}
	// Not found, return nil
	return nil
}

func (this *inputmanager) watch(fd uintptr, flags gopi.FilePollFlags) {
	device := this.devices[fd]
	if device.dev == nil {
		return
	}
	switch flags {
	case gopi.FILEPOLL_FLAG_READ:
		var event linux.EVEvent
		if err := binary.Read(device.dev, binary.LittleEndian, &event); err != nil {
			this.Log.Error(err)
		} else {
			switch event.Type {
			case linux.EV_SYN:
				evDecodeSyn(event, device)
			case linux.EV_KEY:
				evDecodeKey(event, device)
			case linux.EV_ABS:
				evDecodeAbs(event, device)
			case linux.EV_REL:
				evDecodeRel(event, device)
			case linux.EV_MSC:
				evDecodeMsc(event, device)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// DEVICE

func NewDevice(bus uint, exclusive bool) (*device, error) {
	if dev, err := linux.EVOpenDevice(bus); err != nil {
		return nil, err
	} else if name, err := linux.EVGetName(dev.Fd()); err != nil {
		dev.Close()
		return nil, err
	} else if cap, err := linux.EVGetSupportedEventTypes(dev.Fd()); err != nil {
		dev.Close()
		return nil, err
	} else {
		if exclusive {
			if err := linux.EVSetGrabState(dev.Fd(), true); err != nil {
				dev.Close()
				return nil, err
			}
		}
		return &device{
			dev:       dev,
			bus:       bus,
			name:      strings.TrimSpace(name),
			cap:       cap,
			exclusive: exclusive,
		}, nil
	}
}

func (this *device) Close() error {
	if this.dev != nil {
		// remove exclusive access
		if this.exclusive {
			if err := linux.EVSetGrabState(this.dev.Fd(), false); err != nil {
				return err
			}
			this.exclusive = false
		}
		if err := this.dev.Close(); err != nil {
			return err
		}
	}
	this.dev = nil

	// Success
	return nil
}

func (this *device) Name() string {
	return this.name
}

func (this *device) Type() gopi.InputDeviceType {
	// Return cached type
	if this.deviceType != 0 {
		return this.deviceType
	}

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

func (this *device) String() string {
	if this.dev == nil {
		return "<gopi.inputdevice>"
	} else {
		return "<gopi.inputdevice id=" + fmt.Sprint(this.bus) +
			" name=" + strconv.Quote(this.name) +
			" type=" + fmt.Sprint(this.Type()) +
			" capabilities=" + fmt.Sprint(this.cap) +
			">"
	}
}

func (this *device) matches(name string, flags gopi.InputDeviceType) bool {
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

func evTimestamp(event linux.EVEvent, device *device) time.Duration {
	return time.Duration(event.Second)*time.Second + time.Duration(event.Microsecond)*time.Microsecond
}

func evDecodeKey(event linux.EVEvent, device *device) {

}

func evDecodeSyn(event linux.EVEvent, device *device) {

}

func evDecodeAbs(event linux.EVEvent, device *device) {

}

func evDecodeRel(event linux.EVEvent, device *device) {

}

func evDecodeMsc(event linux.EVEvent, device *device) {

}
