// +build linux
// +build !darwin

package platform

import (
	"time"

	// Frameworks
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Platform

// Return uptime
func (this *Platform) Uptime() time.Duration {
	return linux.Uptime()
}

// Return 1, 5 and 15 minute load averages
func (this *Platform) LoadAverages() (float64, float64, float64) {
	return linux.LoadAverage()
}
