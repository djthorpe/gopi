// +build linux
// +build !rpi
// +build !darwin

package platform

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Platform) Type() gopi.PlatformType {
	return gopi.PLATFORM_LINUX
}

// Return serial number
func (this *Platform) SerialNumber() string {
	return linux.SerialNumber()
}

// Return number of displays
func (this *Platform) NumberOfDisplays() uint {
	return 0
}

// Return attached displays
func (this *Platform) AttachedDisplays() []uint {
	return nil
}

// Return product
func (this *Platform) Product() string {
	return "linux"
}
