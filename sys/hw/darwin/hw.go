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
	"syscall"
	"unsafe"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
	#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
	#include <sys/utsname.h>
	#include <stdio.h>
	#include <CoreFoundation/CoreFoundation.h>
	#include <IOKit/IOKitLib.h>
	char* getserial() {
	    io_service_t platformExpert = IOServiceGetMatchingService(kIOMasterPortDefault,IOServiceMatching("IOPlatformExpertDevice"));
		if (platformExpert) {
        CFTypeRef serialNumberAsCFString = IORegistryEntryCreateCFProperty(platformExpert,CFSTR(kIOPlatformSerialNumberKey),kCFAllocatorDefault, 0);
        if (serialNumberAsCFString) {
            CFIndex bufsize = CFStringGetLength(serialNumberAsCFString) + 1;
            char* buf = malloc(bufsize);
            if (buf != NULL) {
                Boolean result = CFStringGetCString(serialNumberAsCFString, buf, bufsize, kCFStringEncodingMacRoman);
                if (result) {
                    free((void*)serialNumberAsCFString);
                    IOObjectRelease(platformExpert);
                    return buf;
				} else {
	                free((void*)buf);
				}
            }
        }
        free((void *)serialNumberAsCFString);
        IOObjectRelease(platformExpert);
    }
    return NULL;
}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Hardware struct{}

type hardware struct {
	log    gopi.Logger
	serial string
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Hardware) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.darwin.Hardware.Open{}")

	this := new(hardware)
	this.log = logger

	// Success
	return this, nil
}

// Close
func (this *hardware) Close() error {
	this.log.Debug("sys.darwin.Hardware.Close{ }")

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetName returns the name of the hardware
func (this *hardware) Name() string {
	model := make([]byte, 80)
	if err := sysctlbyname("hw.model", model); err != nil {
		this.log.Error("<sys.hw.darwin.Metrics>Name: %v", err)
		return ""
	} else {
		return string(model)
	}
}

// SerialNumber returns the serial number of the hardware, if available
func (this *hardware) SerialNumber() string {
	this.log.Debug2("<sys.darwin.Hardware>SerialNumber{}")
	serial := C.getserial()
	defer C.free(unsafe.Pointer(serial))
	if serial == nil {
		this.log.Error("<sys.darwin.Hardware>SerialNumber: Error")
		return ""
	}
	return C.GoString(serial)
}

// Return the number of displays which can be opened
func (this *hardware) NumberOfDisplays() uint {
	return 0
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
		case []byte:
			for i := 0; i < len(val) && i < len(v); i++ {
				v[i] = val[i]
			}
			return nil
		}
		bbuf := bytes.NewBuffer([]byte(val))
		return binary.Read(bbuf, binary.LittleEndian, data)
	}
}
