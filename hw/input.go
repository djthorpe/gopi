/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package hw // import "github.com/djthorpe/gopi/hw"

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type InputDriver interface {
	// Enforces general driver
	gopi.Driver
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Types of input device
const (
	INPUT_TYPE_KEYBOARD    uint8 = 0x01
	INPUT_TYPE_MOUSE       uint8 = 0x02
	INPUT_TYPE_TOUCHSCREEN uint8 = 0x04
	INPUT_TYPE_JOYSTICK    uint8 = 0x08
)
