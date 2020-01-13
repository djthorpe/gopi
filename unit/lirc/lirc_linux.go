// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package lirc

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lirc struct {
	devices  map[uintptr]*lircdev
	filepoll gopi.FilePoll
	bus      gopi.Bus

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// LIRC_CARRIER_FREQUENCY is the default carrier frequency
	LIRC_CARRIER_FREQUENCY = 38000
	// LIRC_DUTY_CYCLE is the default duty cycle
	LIRC_DUTY_CYCLE = 50
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	reDeviceName = regexp.MustCompile("^(\\w+)$")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *lirc) Init(config LIRC) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check incoming parameters
	if config.Filepoll == nil {
		return gopi.ErrBadParameter.WithPrefix("filepoll")
	} else if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.filepoll = config.Filepoll
		this.bus = config.Bus
	}

	// Devices can either be referenced by number (0,1,2) or name (lirc,lirc0, etc)
	if devices := strings.Split(config.Dev, ","); len(devices) == 0 {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.devices = make(map[uintptr]*lircdev, len(devices))

		// Open devices to check for read/write capability
		for _, device := range devices {
			if handle, err := linux.LIRCOpenDevice(device, linux.LIRC_MODE_RCV); err != nil {
				return fmt.Errorf("%s: %w", device, err)
			} else {
				defer handle.Close()
				if features, err := linux.LIRCFeatures(handle.Fd()); err != nil {
					return fmt.Errorf("%s: %w", device, err)
				} else if dev, err := NewDevice(device, features); err != nil {
					return fmt.Errorf("%s: %w", device, err)
				} else if _, exists := this.devices[dev.Fd()]; exists {
					return gopi.ErrInternalAppError
				} else {
					this.devices[dev.Fd()] = dev
				}
			}
		}
	}

	// If there are no devices, then return error
	if len(this.devices) == 0 {
		return gopi.ErrBadParameter.WithPrefix("dev")
	}

	// Here we have a set of devices for recv and sending, so set up watching
	// to read from LIRC here
	// TODO: We need to shutdown more gacefully if any watch files by unwatching
	// all watches setup and also closing any device files open
	for _, device := range this.devices {
		if device.recv {
			if err := this.filepoll.Watch(device.Fd(), gopi.FILEPOLL_FLAG_READ, this.watch); err != nil {
				return err
			}
		}
	}

	// Return success
	return nil
}

func (this *lirc) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Stop watching each device that receives
	for _, device := range this.devices {
		if device.recv {
			if err := this.filepoll.Unwatch(device.Fd()); err != nil {
				return err
			}
		}
		if err := device.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.devices = nil
	this.filepoll = nil
	this.bus = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *lirc) String() string {
	return "<lirc devices=" + fmt.Sprint(this.devices) + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *lirc) watch(fd uintptr, flags gopi.FilePollFlags) {
	if dev, exists := this.devices[fd]; exists {
		switch flags {
		case gopi.FILEPOLL_FLAG_READ:
			if evt, err := dev.Read(this); err != nil {
				this.Log.Error(err)
			} else {
				this.bus.Emit(evt)
			}
		}
	}
}
