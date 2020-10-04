package gopi

import (
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	PlatformType uint32
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Platform interface {
	// Type returns flags identifying platform type
	Type() PlatformType
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PLATFORM_NONE PlatformType = 0
	// OS
	PLATFORM_DARWIN PlatformType = (1 << iota) >> 1
	PLATFORM_RPI
	PLATFORM_LINUX
	// CPU
	PLATFORM_X86_32
	PLATFORM_X86_64
	PLATFORM_BCM2835_ARM6
	PLATFORM_BCM2836_ARM7
	PLATFORM_BCM2837_ARM8
	PLATFORM_BCM2838_ARM8
	// MIN AND MAX
	PLATFORM_MIN = PLATFORM_DARWIN
	PLATFORM_MAX = PLATFORM_BCM2838_ARM8
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p PlatformType) String() string {
	str := ""
	if p == 0 {
		return p.FlagString()
	}
	for v := PLATFORM_MIN; v <= PLATFORM_MAX; v <<= 1 {
		if p&v == v {
			str += "|" + v.FlagString()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (p PlatformType) FlagString() string {
	switch p {
	case PLATFORM_NONE:
		return "PLATFORM_NONE"
	case PLATFORM_DARWIN:
		return "PLATFORM_DARWIN"
	case PLATFORM_RPI:
		return "PLATFORM_RPI"
	case PLATFORM_LINUX:
		return "PLATFORM_LINUX"
	case PLATFORM_X86_32:
		return "PLATFORM_X86_32"
	case PLATFORM_X86_64:
		return "PLATFORM_X86_64"
	case PLATFORM_BCM2835_ARM6:
		return "PLATFORM_BCM2835_ARM6"
	case PLATFORM_BCM2836_ARM7:
		return "PLATFORM_BCM2836_ARM7"
	case PLATFORM_BCM2837_ARM8:
		return "PLATFORM_BCM2837_ARM8"
	case PLATFORM_BCM2838_ARM8:
		return "PLATFORM_BCM2838_ARM8"
	default:
		return "[?? Invalid PlatformType value]"
	}
}
