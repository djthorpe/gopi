// +build linux
// +build !rpi

package platform

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Implementation struct{}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Platform) Type() gopi.PlatformType {
	return gopi.PLATFORM_LINUX
}

// Return serial number
func (this *Platform) SerialNumber() string {
	return linux.SerialNumber()
}

// Return product
func (this *Platform) Product() string {
	return "linux"
}
