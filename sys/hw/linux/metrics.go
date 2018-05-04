// +build linux

/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"fmt"
	"syscall"
	"time"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Metrics struct{}

type metrics struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// Timestamp for module creation
	ts = time.Now()
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open creates a new metrics object, returns error if not possible
func (config Metrics) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.linux.Metrics>Open{}")

	// create new driver
	this := new(metrics)
	this.log = log

	// return driver
	return this, nil
}

// Close connection
func (this *metrics) Close() error {
	this.log.Debug("<sys.hw.linux.Metrics>Close{}")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RETURN UPTIME

func (this *metrics) UptimeHost() time.Duration {
	if info := this.sysinfo(); info != nil {
		return time.Second * time.Duration(info.Uptime)
	} else {
		return 0
	}
}

func (this *metrics) UptimeApp() time.Duration {
	return time.Since(ts)
}

////////////////////////////////////////////////////////////////////////////////
// GET SYSTEM INFO STRUCTURE

func (this *metrics) sysinfo() *syscall.Sysinfo_t {
	info := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&info); err != nil {
		this.log.Error("<sys.hw.linux.Metrics>sysinfo: %v", err)
		return nil
	} else {
		return &info
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *metrics) String() string {
	return fmt.Sprintf("<sys.hw.linux.Metrics>{}")
}
