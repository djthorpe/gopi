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

// Metrics returns various metrics for host and
// custom metrics
type Metrics interface {
	Driver

	// Uptimes for host and for application
	UptimeHost() time.Duration
	UptimeApp() time.Duration

	// Load Average (1, 5 and 15 minutes)
	LoadAverage() (float64, float64, float64)
}
