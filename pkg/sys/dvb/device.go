// +build dvb

package dvb

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Device struct {
	Adapter  uint
	Frontend []uint
	Dvr      []uint
	Demux    []uint
	Net      []uint
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DVB_ADAPTER_GLOB  = "/dev/dvb/adapter"
	DVB_ADAPTER_UNITS = "demux dvr frontend net"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Devices() []Device {
	adapters, err := filepath.Glob(DVB_ADAPTER_GLOB + "*")
	if err != nil {
		return nil
	}
	devices := make([]Device, 0, len(adapters))
	for _, path := range adapters {
		if n, err := strconv.ParseUint(strings.TrimPrefix(path, DVB_ADAPTER_GLOB), 0, 32); err == nil {
			d := Device{Adapter: uint(n)}
			for _, unit := range strings.Fields(DVB_ADAPTER_UNITS) {
				prefix := filepath.Join(path, unit)
				paths, err := filepath.Glob(prefix + "*")
				if err != nil {
					continue
				}
				for _, path := range paths {
					if m, err := strconv.ParseUint(strings.TrimPrefix(path, prefix), 0, 32); err == nil {
						switch unit {
						case "frontend":
							d.Frontend = append(d.Frontend, uint(m))
						case "demux":
							d.Demux = append(d.Demux, uint(m))
						case "dvr":
							d.Dvr = append(d.Dvr, uint(m))
						case "net":
							d.Net = append(d.Net, uint(m))
						}
					}
				}
			}
			devices = append(devices, d)
		}
	}
	return devices
}

// Path returns the path to an adaptor device
func (d Device) Path(unit string, m uint) string {
	path := filepath.Join(DVB_ADAPTER_GLOB+fmt.Sprint(d.Adapter), unit+fmt.Sprint(m))
	if _, err := os.Stat(path); err != nil {
		return ""
	} else {
		return path
	}
}

// Open frontend, return file
func (d Device) FEOpen(mode int) (*os.File, error) {
	if len(d.Frontend) == 0 {
		return nil, gopi.ErrNotFound
	} else if path := d.Path("frontend", d.Frontend[0]); path == "" {
		return nil, gopi.ErrNotFound
	} else if fh, err := os.OpenFile(path, mode, 0); err != nil {
		return nil, err
	} else {
		return fh, nil
	}
}

// Open demux, return file
func (d Device) DMXOpen() (*os.File, error) {
	if len(d.Demux) == 0 {
		return nil, gopi.ErrNotFound
	} else if path := d.Path("demux", d.Demux[0]); path == "" {
		return nil, gopi.ErrNotFound
	} else if fh, err := os.OpenFile(path, os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return fh, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d Device) String() string {
	str := "<dvb.device"
	str += " adaptor=" + fmt.Sprint(d.Adapter)
	if len(d.Frontend) > 0 {
		str += " frontend=" + fmt.Sprint(d.Frontend)
	}
	if len(d.Demux) > 0 {
		str += " demux=" + fmt.Sprint(d.Demux)
	}
	if len(d.Dvr) > 0 {
		str += " dvr=" + fmt.Sprint(d.Dvr)
	}
	if len(d.Net) > 0 {
		str += " net=" + fmt.Sprint(d.Net)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Call ioctl
func dvb_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
