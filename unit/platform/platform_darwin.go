// +build darwin

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
	darwin "github.com/djthorpe/gopi/v2/sys/darwin"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

func (this *platform) Init() error {
	// No special init for darwin
	return nil
}

func (this *platform) Platform() gopi.PlatformType {
	return gopi.PLATFORM_DARWIN
}

// Return serial number
func (this *platform) SerialNumber() string {
	return darwin.SerialNumber()
}

// Return uptime
func (this *platform) Uptime() time.Duration {
	return darwin.Uptime()
}

// Return 1, 5 and 15 minute load averages
func (this *platform) LoadAverages() (float64, float64, float64) {
	return darwin.LoadAverage()
}
