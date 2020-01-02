// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package platform

import (
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

func (this *platform) Init() error {
	// No special init for linux
	return nil
}

func (this *platform) Platform() gopi.PlatformType {
	return gopi.PLATFORM_LINUX
}

// Return serial number
func (this *platform) SerialNumber() string {
	return linux.SerialNumber()
}

// Return uptime
func (this *platform) Uptime() time.Duration {
	return linux.Uptime()
}

// Return 1, 5 and 15 minute load averages
func (this *platform) LoadAverages() (float64, float64, float64) {
	return linux.LoadAverage()
}
