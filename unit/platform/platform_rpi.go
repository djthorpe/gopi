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

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

func (this *platform) Init() error {
	// Initialise
	if err := rpi.BCMHostInit(); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *platform) Type() gopi.PlatformType {
	return gopi.PLATFORM_RPI | gopi.PLATFORM_LINUX
}

// Return serial number
func (this *platform) SerialNumber() string {
	if serial, _, err := rpi.VCGetSerialRevision(); err != nil {
		this.Log.Error(err)
		return ""
	} else {
		return fmt.Sprintf("%08X", serial)
	}
}
