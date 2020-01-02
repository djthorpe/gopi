/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package platform

import (
	"fmt"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type Platform struct{}

type platform struct {
	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Platform) Name() string { return "gopi.Platform" }

func (config Platform) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(platform)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *platform) String() string {
	return fmt.Sprintf("<gopi.Platform>")
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

func (this *platform) Platform() gopi.PlatformType {
	return gopi.PLATFORM_NONE
}

// Return serial number
func (this *platform) SerialNumber() string {
	return ""
}

// Return uptime
func (this *platform) Uptime() time.Duration {
	return 0
}

// Return 1, 5 and 15 minute load averages
func (this *platform) LoadAverages() (float32, float32, float32) {
	return 0, 0, 0
}
