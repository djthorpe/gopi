// +build linux
// +build !rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package platform

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

type Implementation struct{}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

func (this *platform) Init() error {
	// No special init for linux
	return nil
}

func (this *platform) Type() gopi.PlatformType {
	return gopi.PLATFORM_LINUX
}

// Return serial number
func (this *platform) SerialNumber() string {
	return linux.SerialNumber()
}

// Return number of displays
func (this *platform) NumberOfDisplays() uint {
	return 0
}

// Return attached displays
func (this *platform) AttachedDisplays() []uint {
	return nil
}

// Return product
func (this *platform) Product() string {
	return "linux"
}
