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
	return "<lirc" +
		" rcv_mode=" + fmt.Sprint(this.RcvMode()) +
		" send_mode=" + fmt.Sprint(this.SendMode()) +
		" devices=" + fmt.Sprint(this.devices) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// SEND AND RECV MODE PARAMETERS

func (this *lirc) RcvMode() gopi.LIRCMode {
	mode := gopi.LIRC_MODE_NONE
	// Ensure all recv devices have the same mode, or return gopi.LIRC_MODE_NONE
	for _, device := range this.devices {
		if device.recv == false {
			continue
		} else if mode == gopi.LIRC_MODE_NONE {
			mode = device.RcvMode()
		} else if mode != device.RcvMode() {
			return gopi.LIRC_MODE_NONE
		}
	}
	return mode
}

func (this *lirc) SendMode() gopi.LIRCMode {
	mode := gopi.LIRC_MODE_NONE
	// Ensure all send devices have the same mode, or return gopi.LIRC_MODE_NONE
	for _, device := range this.devices {
		if device.send == false {
			continue
		} else if mode == gopi.LIRC_MODE_NONE {
			mode = device.SendMode()
		} else if mode != device.SendMode() {
			return gopi.LIRC_MODE_NONE
		}
	}
	return mode
}

func (this *lirc) SetRcvMode(mode gopi.LIRCMode) error {
	set :=false

	// Set mode for all recv devices
	for _, device := range this.devices {
		if device.recv == false {
			continue
		} else if err := device.SetRcvMode(mode); err != nil {
			return fmt.Errorf("%s: %w",device.Name(),err)
		} else {
			set = true
		}
	}

	// No recv devices
	if set == false {
		return gopi.ErrNotImplemented.WithPrefix("SetRcvMode")
	}
	
	// Success
	return nil
}

func (this *lirc) SetSendMode(mode gopi.LIRCMode) error {
	set :=false
	
	// Set mode for all recv devices
	for _, device := range this.devices {
		if device.send == false {
			continue
		} else if err := device.SetSendMode(mode); err != nil {
			return fmt.Errorf("%s: %w",device.Name(),err)
		} else {
			set = true
		}
	}

	// No send devices
	if set == false {
		return gopi.ErrNotImplemented.WithPrefix("SetSendMode")
	}

	// Success
	return nil
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
