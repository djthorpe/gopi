// +build rpi
// +build !darwin

package platform

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Implementation struct {
	gopi.Unit
	i rpi.VCHIInstance
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Implementation) New(cfg gopi.Config) error {
	// Initialise TV Service
	if this.i = rpi.VCHI_Init(); this.i == nil {
		return gopi.ErrInternalAppError
	} else if _, err := rpi.VCHI_TVInit(this.i); err != nil {
		return err
	}

	// Initialise BCMHost
	if err := rpi.BCMHostInit(); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Implementation) Dispose() error {
	// Stop TV Service
	return rpi.VCHI_TVStop(this.i)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Platform) Type() gopi.PlatformType {
	platform := gopi.PLATFORM_RPI | gopi.PLATFORM_LINUX
	if _, product, err := rpi.VCGetSerialProduct(); err != nil {
		return platform
	} else {
		productinfo := rpi.NewProductInfo(product)
		switch productinfo.Processor {
		case rpi.RPI_PROCESSOR_BCM2835:
			platform |= gopi.PLATFORM_BCM2835_ARM6
		case rpi.RPI_PROCESSOR_BCM2836:
			platform |= gopi.PLATFORM_BCM2836_ARM7
		case rpi.RPI_PROCESSOR_BCM2837:
			platform |= gopi.PLATFORM_BCM2837_ARM8
		case rpi.RPI_PROCESSOR_BCM2838:
			platform |= gopi.PLATFORM_BCM2838_ARM8
		}
		return platform
	}
}

// Return serial number
func (this *Platform) SerialNumber() string {
	if serial, _, err := rpi.VCGetSerialProduct(); err != nil {
		return ""
	} else {
		return fmt.Sprintf("%08X", serial)
	}
}

// Return product name
func (this *Platform) Product() string {
	if _, product, err := rpi.VCGetSerialProduct(); err != nil {
		return ""
	} else {
		productinfo := rpi.NewProductInfo(product)
		return fmt.Sprint(productinfo.Model)
	}
}

// Return number of displays
func (this *Platform) NumberOfDisplays() uint {
	return uint(rpi.DXNumberOfDisplays())
}

// Return attached displays
func (this *Platform) AttachedDisplays() []uint {
	if displays, err := rpi.VCHI_TVGetAttachedDevices(); err != nil {
		return nil
	} else {
		displays_ := make([]uint, len(displays))
		for i, display := range displays {
			displays_[i] = uint(display)
		}
		return displays_
	}
}
