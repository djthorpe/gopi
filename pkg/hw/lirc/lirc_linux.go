// +build linux

package lirc

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
	multierror "github.com/hashicorp/go-multierror"

	_ "github.com/djthorpe/gopi/v3/pkg/event"
	_ "github.com/djthorpe/gopi/v3/pkg/file"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lirc struct {
	gopi.Unit
	gopi.FilePoll
	gopi.Publisher
	gopi.Logger
	sync.Mutex

	devices map[uintptr]*lircdev
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// LIRC_DUTY_CYCLE is the default duty cycle
	LIRC_DUTY_CYCLE = 50
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *lirc) Define(cfg gopi.Config) error {
	cfg.FlagString("lirc.dev", "0,1", "Comma-separated list of LIRC devices")
	return nil
}

func (this *lirc) New(cfg gopi.Config) error {
	if this.FilePoll == nil || this.Publisher == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing FilePoll or Publisher")
	}

	if devices, err := getDevices(cfg.GetString("lirc.dev")); err != nil {
		return err
	} else if len(devices) == 0 {
		return gopi.ErrBadParameter.WithPrefix("No LIRC devices")
	} else {
		this.devices = devices
	}

	// Watch receive devices
	for _, device := range this.devices {
		if device.recv {
			if err := this.FilePoll.Watch(device.Fd(), gopi.FILEPOLL_FLAG_READ, this.ReadEvent); err != nil {
				return err
			}
		}
	}

	// Set the duty cycle to default values
	for _, device := range this.devices {
		if device.recv {
			if err := device.SetRcvDutyCycle(LIRC_DUTY_CYCLE); err != nil && errors.Is(err, gopi.ErrNotImplemented) == false {
				return err
			}
		}
		if device.send {
			if err := device.SetSendDutyCycle(LIRC_DUTY_CYCLE); err != nil && errors.Is(err, gopi.ErrNotImplemented) == false {
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

	// Result captures any errors on disposing
	var result error

	// Stop watching each device that receives
	for _, device := range this.devices {
		if device.recv {
			if err := this.FilePoll.Unwatch(device.Fd()); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Close devices
	for _, device := range this.devices {
		if err := device.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.devices = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - EVENTS

func (this *lirc) ReadEvent(fd uintptr, flags gopi.FilePollFlags) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, exists := this.devices[fd]; exists {
		if evt, err := device.ReadEvent(fd, flags); err != nil {
			this.Print("ReadEvent: ", device.Name(), err)
		} else if err := this.Publisher.Emit(evt, true); err != nil {
			this.Print("ReadEvent: ", device.Name(), err)
		}
	} else {
		this.Print("Not watching fd=", fd)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SEND & RECEIVE MODES

func (this *lirc) RecvMode() gopi.LIRCMode {
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

func (this *lirc) SetRecvMode(mode gopi.LIRCMode) error {
	var result error
	var set bool

	// Set mode for all recv devices
	for _, device := range this.devices {
		if device.recv {
			if err := device.SetRcvMode(mode); err != nil {
				result = multierror.Append(result, err)
			}
			set = true
		}
	}
	if set {
		return nil
	} else {
		return gopi.ErrNotImplemented
	}
}

func (this *lirc) SetSendMode(mode gopi.LIRCMode) error {
	var result error
	var set bool

	// Set mode for all send devices
	for _, device := range this.devices {
		if device.send {
			if err := device.SetSendMode(mode); err != nil {
				result = multierror.Append(result, err)
			}
			set = true
		}
	}
	if set {
		return nil
	} else {
		return gopi.ErrNotImplemented
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - DUTY CYCLES

func (this *lirc) RecvDutyCycle() uint32 {
	return 0
}

func (this *lirc) SendDutyCycle() uint32 {
	return 0
}

func (this *lirc) SetRecvDutyCycle(uint32) error {
	return gopi.ErrNotImplemented

}

func (this *lirc) SetSendDutyCycle(uint32) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *lirc) RecvResolutionMicros() uint32 {
	return 0
}

func (this *lirc) SetRecvTimeoutMs(uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetRecvTimeoutReports(bool) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetRecvCarrierHz(uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetRecvCarrierRangeHz(min, max uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetSendCarrierHz(uint32) error {
	return gopi.ErrNotImplemented

}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SEND

func (this *lirc) PulseSend([]uint32) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *lirc) String() string {
	str := "<lirc"
	if len(this.devices) > 0 {
		str += " devices="
		for _, device := range this.devices {
			str += fmt.Sprint(device) + " "
		}
	}
	if mode := this.RecvMode(); mode != gopi.LIRC_MODE_NONE {
		str += " recv_mode=" + fmt.Sprint(mode)
	}
	if mode := this.SendMode(); mode != gopi.LIRC_MODE_NONE {
		str += " send_mode=" + fmt.Sprint(mode)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func getDevices(value string) (map[uintptr]*lircdev, error) {
	devices := make(map[uintptr]*lircdev)

	for _, device := range strings.Split(value, ",") {
		handle, err := linux.LIRCOpenDevice(device, linux.LIRC_MODE_RCV)
		if os.IsNotExist(err) {
			// Skip any LIRC devices which don't exist
			continue
		}
		if err != nil {
			return nil, err
		}
		defer handle.Close()

		// Get features and create a new LIRC device
		if features, err := linux.LIRCFeatures(handle.Fd()); err != nil {
			return nil, err
		} else if device, err := NewDevice(handle.Name(), features); err != nil {
			return nil, err
		} else {
			key := device.Fd()
			if _, exists := devices[key]; exists {
				return nil, gopi.ErrDuplicateEntry
			} else {
				devices[key] = device
			}
		}
	}

	// Return success
	return devices, nil
}
