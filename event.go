/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "time"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// Unique ID for each timer created
	TimerId uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Timer unit - sends out messages on the event bus
type Timer interface {
	Unit

	// Create periodic event at regular intervals
	NewTicker(time.Duration) TimerId
}
