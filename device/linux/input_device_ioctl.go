/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"os"
	"unsafe"
	"syscall"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
 #include <linux/input.h>
 static int _EVIOCGNAME(int len)        { return EVIOCGNAME(len); }
 static int _EVIOCGPHYS(int len)        { return EVIOCGPHYS(len); }
 static int _EVIOCGUNIQ(int len)        { return EVIOCGUNIQ(len); }
 static int _EVIOCGPROP(int len)        { return EVIOCGPROP(len); }
 static int _EVIOCGKEY(int len)         { return EVIOCGKEY(len); }
 static int _EVIOCGLED(int len)         { return EVIOCGLED(len); }
 static int _EVIOCGSND(int len)         { return EVIOCGSND(len); }
 static int _EVIOCGSW(int len)          { return EVIOCGSW(len); }
 static int _EVIOCGBIT(int ev, int len) { return EVIOCGBIT(ev, len); }
 static int _EVIOCGABS(int abs)         { return EVIOCGABS(abs); }
 static int _EVIOCSABS(int abs)         { return EVIOCSABS(abs); }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Internal constants
const (
	MAX_IOCTL_SIZE_BYTES = 256
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	EVIOCGNAME = uintptr(C._EVIOCGNAME(MAX_IOCTL_SIZE_BYTES)) // get device name
	EVIOCGPHYS = uintptr(C._EVIOCGPHYS(MAX_IOCTL_SIZE_BYTES)) // get physical location
	EVIOCGUNIQ = uintptr(C._EVIOCGUNIQ(MAX_IOCTL_SIZE_BYTES)) // get unique identifier
	EVIOCGPROP = uintptr(C._EVIOCGPROP(MAX_IOCTL_SIZE_BYTES)) // get device properties
	EVIOCGID   = uintptr(C.EVIOCGID)                          // get device ID
)

////////////////////////////////////////////////////////////////////////////////
// IOCTL FUNCTIONS

// Get name of the device
func evGetName(handle *os.File) (string, error) {
	name := new([MAX_IOCTL_SIZE_BYTES]C.char)
	err := evIoctl(handle.Fd(), uintptr(EVIOCGNAME), unsafe.Pointer(name))
	if err != 0 {
		return "", err
	}
	return C.GoString(&name[0]), nil
}

// Get physical connection string
func evGetPhys(handle *os.File) (string, error) {
	name := new([MAX_IOCTL_SIZE_BYTES]C.char)
	err := evIoctl(handle.Fd(), uintptr(EVIOCGPHYS), unsafe.Pointer(name))
	if err != 0 {
		return "", err
	}
	return C.GoString(&name[0]), nil
}

// Get device information (bus, vendor, product, version)
func evGetInfo(handle *os.File) (uint16,uint16,uint16,uint16,error) {
	info := [4]uint16{ }
	err := evIoctl(handle.Fd(), uintptr(EVIOCGID), unsafe.Pointer(&info))
	if err != 0 {
		return uint16(0),uint16(0),uint16(0),uint16(0),err
	}
	return info[0],info[1],info[2],info[3],nil
}

// Get device capabilities
func evGetSupportedEventTypes(handle *os.File) ([]evType,error) {
	evbits := new([EV_MAX >> 3]byte)
	err := evIoctl(handle.Fd(),uintptr(C._EVIOCGBIT(C.int(0), C.int(EV_MAX))), unsafe.Pointer(evbits))
	if err != 0 {
		return nil,err
	}
	capabilities := make([]evType,0)
	evtype := evType(0)
	for i := 0; i < len(evbits); i++ {
		evbyte := evbits[i]
		for j := 0; j < 8; j++ {
			if evbyte & 0x01 != 0x00 {
				capabilities = append(capabilities,evtype)
			}
			evbyte = evbyte >> 1
			evtype++
		}
	}
	return capabilities,nil
}

// Call ioctl
func evIoctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}

