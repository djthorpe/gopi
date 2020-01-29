// +build darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package darwin

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"time"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
	#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
	#include <sys/utsname.h>
	#include <stdio.h>
	#include <stdlib.h>
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
// PUBLIC METHODS

// Returns CPU from sysctl
func CPUType() uint64 {
	var cputype uint64
	if err := sysctlbyname("hw.cputype", &cputype); err != nil {
		return 0
	} else {
		return cputype
	}
}

func CPU64Bit() bool {
	var capable uint64
	if err := sysctlbyname("hw.cpu64bit_capable", &capable); err != nil {
		return false
	} else {
		return capable != 0
	}
}

// SerialNumber returns the serial number of the hardware, if available
func SerialNumber() string {
	serial := C.getserial()
	defer C.free(unsafe.Pointer(serial))
	if serial == nil {
		return ""
	} else {
		return C.GoString(serial)
	}
}

// Uptime returns the duration the machine has been switched on for
func Uptime() time.Duration {
	tv := syscall.Timeval32{}
	if err := sysctlbyname("kern.boottime", &tv); err != nil {
		return 0
	} else {
		return time.Since(time.Unix(int64(tv.Sec), int64(tv.Usec)*1000))
	}
}

// LoadAverage returns load averages (1,5 and 15 minute averages)
func LoadAverage() (float64, float64, float64) {
	avg := []C.double{0, 0, 0}
	if C.getloadavg(&avg[0], C.int(len(avg))) == C.int(-1) {
		return 0, 0, 0
	} else {
		return float64(avg[0]), float64(avg[1]), float64(avg[2])
	}
}

// Product returns the name of the hardware
func Product() string {
	model := make([]byte, 80)
	if err := sysctlbyname("hw.model", model); err != nil {
		return ""
	} else {
		return C.GoString((*C.char)(unsafe.Pointer(&model[0])))
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

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
