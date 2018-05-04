// +build darwin

/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package darwin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	// Frameworks
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
	tv := syscall.Timeval32{}

	if err := sysctlbyname("kern.boottime", &tv); err != nil {
		this.log.Error("<sys.hw.linux.Metrics>UptimeHost: %v", err)
		return 0
	} else {
		return time.Since(time.Unix(int64(tv.Sec), int64(tv.Usec)*1000))
	}
}

func (this *metrics) UptimeApp() time.Duration {
	return time.Since(ts)
}

////////////////////////////////////////////////////////////////////////////////
// GET SYSTEM INFO

// Generic Sysctl buffer unmarshalling
// from https://github.com/cloudfoundry/gosigar/blob/master/sigar_darwin.go
func sysctlbyname(name string, data interface{}) error {
	if val, err := syscall.Sysctl(name); err != nil {
		return err
	} else {
		buf := []byte(val)
		switch v := data.(type) {
		case *uint64:
			*v = *(*uint64)(unsafe.Pointer(&buf[0]))
			return nil
		}
		bbuf := bytes.NewBuffer([]byte(val))
		return binary.Read(bbuf, binary.LittleEndian, data)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *metrics) String() string {
	return fmt.Sprintf("<sys.hw.linux.Metrics>{}")
}
