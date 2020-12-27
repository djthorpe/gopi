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

	displays map[uint32]*display
	instance rpi.VCHIInstance
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Initialise DX
	rpi.DXInit()

	// Initialise TV Service
	if instance := rpi.VCHI_Init(); instance == nil {
		return gopi.ErrInternalAppError.WithPrefix("VCHI_Init")
	} else if _, err := rpi.VCHI_TVInit(instance); err != nil {
		return err
	} else {
		this.instance = instance
	}

	// Make displays
	this.displays = make(map[uint32]*display, 10)

	// Events callback
	rpi.VCTV_RegisterCallback(func(evt rpi.TVDisplayStateFlag, id rpi.DXDisplayId) {
		var display gopi.Display
		if id != 0 {
			display, _ = this.Display(uint32(id))
		}
		if this.Publisher != nil {
			switch evt {
			case rpi.TV_STATE_HDMI_ATTACHED:
				this.Publisher.Emit(NewEvent(display, gopi.DISPLAY_FLAG_ATTACHED), true)
			case rpi.TV_STATE_HDMI_UNPLUGGED:
				this.Publisher.Emit(NewEvent(display, gopi.DISPLAY_FLAG_UNPLUGGED), true)
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

	// Dispose of displays
	for _, display := range this.displays {
		if display != nil {
			if err := display.Dispose(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Stop TV Service
	if this.instance != nil {
		if err := rpi.VCHI_TVStop(this.instance); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Stop DX
	rpi.DXStop()

	// Release resources
	this.displays = nil
	this.instance = nil

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return one display
func (this *Manager) Display(id uint32) (gopi.Display, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if display, exists := this.displays[id]; exists {
		return display, nil
	} else if display, err := NewDisplay(rpi.DXDisplayId(id)); err != nil {
		return nil, err
	} else {
		this.displays[id] = display
		return display, nil
	}
}

// Return all displays
func (this *Manager) Displays() []gopi.Display {
	this.RWMutex.Lock()
	displays, err := rpi.VCHI_TVGetAttachedDevices()
	if err != nil {
		this.Debug("VCHI_TVGetAttachedDevices: ", err)
		return nil
	}
	this.RWMutex.Unlock()

	result := make([]gopi.Display, 0, len(displays))
	for _, id := range displays {
		if display, err := this.Display(uint32(id)); err != nil {
			this.Debug("Displays: ", id, err)
		} else {
			result = append(result, display)
		}
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
