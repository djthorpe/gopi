/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	PlatformType uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Platform interface {

	// Return platform type
	Platform() PlatformType

	// Return serial number
	SerialNumber() string

	// Return uptime
	Uptime() time.Duration

	// Return 1, 5 and 15 minute load averages
	LoadAverages() (float32, float32, float32)

	// Implements gopi.Unit
	Unit
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PLATFORM_NONE   PlatformType = 0
	PLATFORM_DARWIN PlatformType = (1 << iota) >> 1
	PLATFORM_RPI
	PLATFORM_LINUX
	PLATFORM_MAX = PLATFORM_LINUX
)
