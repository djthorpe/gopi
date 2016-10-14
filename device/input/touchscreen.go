/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md

	This package provides input mechanisms, including the touchscreen
	interface for the official Raspberry Pi LED.
*/
package input

import (
	"image"
	"time"
)

////////////////////////////////////////////////////////////////////////////////

// Exported event structure, which is fired when an event occurs
type TouchEvent struct {
	Slot       uint32
	Timestamp  time.Duration
	Point      image.Point
	LastPoint  image.Point
	Identifier int
}

// Non-exported raw event data structure sent over the wire
type rawEvent struct {
	Second      uint32
	Microsecond uint32
	Type        uint16
	Code        uint16
	Value       uint32
}

////////////////////////////////////////////////////////////////////////////////

const (
	// References:

	// Event types
	// https://www.kernel.org/doc/Documentation/input/event-codes.txt
	EV_SYN uint16 = 0x0000
	EV_KEY uint16 = 0x0001
	EV_ABS uint16 = 0x0003

	// Button information
	BTN_TOUCH         uint16 = 0x014A
	BTN_TOUCH_RELEASE uint32 = 0x00000000
	BTN_TOUCH_PRESS   uint32 = 0x00000001

	// Multi-Touch Types
	// https://www.kernel.org/doc/Documentation/input/multi-touch-protocol.txt
	ABS_X              uint16 = 0x0000
	ABS_Y              uint16 = 0x0001
	ABS_MT_SLOT        uint16 = 0x002F // 47 MT slot being modified
	ABS_MT_POSITION_X  uint16 = 0x0035 // 53 Center X of multi touch position
	ABS_MT_POSITION_Y  uint16 = 0x0036 // 54 Center Y of multi touch position
	ABS_MT_TRACKING_ID uint16 = 0x0039 // 57 Unique ID of initiated contact
)
