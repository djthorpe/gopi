// +build darwin
// +build !linux

package platform

import (
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	darwin "github.com/djthorpe/gopi/v3/pkg/sys/darwin"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	CPU_TYPE_VAX     = 1
	CPU_TYPE_MC680x0 = 6
	CPU_TYPE_X86     = 7
	CPU_TYPE_MIPS    = 8
	CPU_TYPE_MC98000 = 10
	CPU_TYPE_HPPA    = 11
	CPU_TYPE_ARM     = 12
	CPU_TYPE_SPARC   = 14
	CPU_TYPE_I860    = 15
	CPU_TYPE_ALPHA   = 16
	CPU_TYPE_POWERPC = 18
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Platform) Type() gopi.PlatformType {
	platform := gopi.PLATFORM_DARWIN
	cputype := darwin.CPUType()
	cpu64 := darwin.CPU64Bit()
	switch {
	case cputype == CPU_TYPE_X86 && cpu64 == false:
		platform = platform | gopi.PLATFORM_X86_32
	case cputype == CPU_TYPE_X86 && cpu64 == true:
		platform = platform | gopi.PLATFORM_X86_64
	}
	return platform
}

// Return serial number
func (this *Platform) SerialNumber() string {
	return darwin.SerialNumber()
}

// Return uptime
func (this *Platform) Uptime() time.Duration {
	return darwin.Uptime()
}

// Return 1, 5 and 15 minute load averages
func (this *Platform) LoadAverages() (float64, float64, float64) {
	return darwin.LoadAverage()
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
	return darwin.Product()
}
