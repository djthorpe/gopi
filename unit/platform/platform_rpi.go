// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package platform

import (
	"fmt"

	"github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

type Implementation struct {
	instance rpi.VCHIInstance
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

func (this *platform) Init() error {
	// Initialise TV Service
	if this.instance = rpi.VCHI_Init(); this.instance == nil {
		return gopi.ErrInternalAppError
	} else if _, err := rpi.VCHI_TVInit(this.instance); err != nil {
		return err
	}

	// Initialise BCMHost
	if err := rpi.BCMHostInit(); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *platform) Close() error {
	// Stop TV Service
	if err := rpi.VCHI_TVStop(this.instance); err != nil {
		return err
	}
	return this.Unit.Close()
}

func (this *platform) Type() gopi.PlatformType {
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
func (this *platform) SerialNumber() string {
	if serial, _, err := rpi.VCGetSerialProduct(); err != nil {
		this.Log.Error(err)
		return ""
	} else {
		return fmt.Sprintf("%08X", serial)
	}
}

// Return product name
func (this *platform) Product() string {
	if _, product, err := rpi.VCGetSerialProduct(); err != nil {
		this.Log.Error(err)
		return ""
	} else {
		productinfo := rpi.NewProductInfo(product)
		return fmt.Sprint(productinfo.Model)
	}
}

// Return number of displays
func (this *platform) NumberOfDisplays() uint {
	return uint(rpi.DXNumberOfDisplays())
}

// Return attached displays
func (this *platform) AttachedDisplays() []uint {
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
