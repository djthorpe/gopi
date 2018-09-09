/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Timer implements a time keeping driver
type Timer interface {
	Driver
	Publisher

	// Schedule a timeout (one shot)
	NewTimeout(duration time.Duration, userInfo interface{}) error

	// Schedule an interval, which can fire immediately
	NewInterval(duration time.Duration, userInfo interface{}, immediately bool) error

	// Schedule a backoff timer with maximum backoff duration
	NewBackoff(duration time.Duration, max_duration time.Duration, userInfo interface{}) error
}
