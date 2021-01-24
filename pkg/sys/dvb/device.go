// +build dvb

package dvb

import (
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Device struct {
	Adapter uint
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DVB_ADAPTER_GLOB  = "/dev/dvb/adapter"
	DVB_ADAPTER_UNITS = "demux,dvr,frontend,net"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Devices() []*Device {
	adapters, err := filepath.Glob(DVB_ADAPTER_GLOB + "*")
	if err != nil {
		return nil
	}
	devices := make([]uint, 0, len(adapters))
	for _, path := range adapters {
		if n, err := strconv.ParseUint(strings.TrimPrefix(path, DVB_ADAPTER_GLOB), 0, 32); err == nil {
			devices = append(devices, Device{uint(n)})
		}
	}
	return devices, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Call ioctl
func dvb_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
