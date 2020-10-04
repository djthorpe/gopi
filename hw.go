package gopi

import (
	"fmt"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	PlatformType uint32

	// SPIMode is the SPI Mode
	SPIMode uint8
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Platform interface {
	Product() string                           // Product returns product name
	Type() PlatformType                        // Type returns flags identifying platform type
	SerialNumber() string                      // SerialNumber returns unique serial number for host
	Uptime() time.Duration                     // Uptime returns uptime for host
	LoadAverages() (float64, float64, float64) // LoadAverages returns 1, 5 and 15 minute load averages
	NumberOfDisplays() uint                    // NumberOfDisplays returns the number of possible displays for this host
	AttachedDisplays() []uint                  // AttachedDisplays returns array of displays which are actually attached
}

// SPI implements the SPI interface for sensors, etc.
type SPI interface {
	Mode() SPIMode                        // Get SPI mode
	MaxSpeedHz() uint32                   // Get SPI speed
	BitsPerWord() uint8                   // Get Bits Per Word
	SetMode(SPIMode) error                // Set SPI mode
	SetMaxSpeedHz(uint32) error           // Set SPI speed
	SetBitsPerWord(uint8) error           // Set Bits Per Word
	Transfer(send []byte) ([]byte, error) // Read/Write
	Read(len uint32) ([]byte, error)      // Read
	Write(send []byte) error              // Write
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

const (
	SPI_MODE_CPHA SPIMode = 0x01
	SPI_MODE_CPOL SPIMode = 0x02
	SPI_MODE_0    SPIMode = 0x00
	SPI_MODE_1    SPIMode = (0x00 | SPI_MODE_CPHA)
	SPI_MODE_2    SPIMode = (SPI_MODE_CPOL | 0x00)
	SPI_MODE_3    SPIMode = (SPI_MODE_CPOL | SPI_MODE_CPHA)
	SPI_MODE_NONE SPIMode = 0xFF
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

func (m SPIMode) String() string {
	switch m {
	case SPI_MODE_0:
		return "SPI_MODE_0"
	case SPI_MODE_1:
		return "SPI_MODE_1"
	case SPI_MODE_2:
		return "SPI_MODE_2"
	case SPI_MODE_3:
		return "SPI_MODE_3"
	default:
		return "[?? Invalid SPIMode " + fmt.Sprint(uint(m)) + "]"
	}
}
