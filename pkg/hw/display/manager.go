// +build rpi

package display

import (
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	sync.RWMutex
	gopi.Logger
	gopi.Publisher
	gopi.Platform

	instance rpi.VCHIInstance
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Initialise TV Service
	if instance := rpi.VCHI_Init(); instance == nil {
		return gopi.ErrInternalAppError.WithPrefix("VCHI_Init")
	} else if _, err := rpi.VCHI_TVInit(instance); err != nil {
		return err
	} else {
		this.instance = instance
	}

	// Events callback
	rpi.VCTV_RegisterCallback(func(evt rpi.TVDisplayStateFlag, id rpi.DXDisplayId) {
		if this.Publisher != nil {
			switch evt {
			case rpi.TV_STATE_HDMI_ATTACHED:
				this.Publisher.Emit(NewEvent(NewDisplay(id), gopi.DISPLAY_FLAG_ATTACHED), true)
			case rpi.TV_STATE_HDMI_UNPLUGGED:
				this.Publisher.Emit(NewEvent(NewDisplay(id), gopi.DISPLAY_FLAG_ATTACHED), true)
			}
		}
	})

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Unregister callback
	rpi.VCTV_RegisterCallback(nil)

	if this.instance != nil {
		if err := rpi.VCHI_TVStop(this.instance); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.instance = nil

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return all displays
func (this *Manager) Displays() []gopi.Display {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	displays, err := rpi.VCHI_TVGetAttachedDevices()
	if err != nil {
		this.Debug("VCHI_TVGetAttachedDevices: ", err)
		return nil
	}

	result := make([]gopi.Display, 0, len(displays))
	for _, id := range displays {
		result = append(result, NewDisplay(id))
	}
	return result
}

func (this *Manager) PowerOn(display gopi.Display) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	return rpi.VCHI_TVHDMIPowerOnPreferred(rpi.DXDisplayId(display.Id()))
}

func (this *Manager) PowerOff(display gopi.Display) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	return rpi.VCHI_TVPowerOff(rpi.DXDisplayId(display.Id()))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<display.manager"
	return str + ">"
}
