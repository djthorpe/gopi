// +build linux
// +build !darwin

package linux

import (
	"io/ioutil"
	"net"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	TEMPERATURE_PATH = "/sys/class/thermal"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// SerialNumber returns the mac address of the hardware, if available
func SerialNumber() string {
	if ifaces, err := net.Interfaces(); err != nil {
		return ""
	} else if len(ifaces) == 0 {
		return ""
	} else {
		macs := make([]string, 0, len(ifaces))
		for _, iface := range ifaces {
			if iface.Flags&net.FlagLoopback == net.FlagLoopback {
				continue
			}
			if iface.HardwareAddr != nil {
				macs = append(macs, iface.HardwareAddr.String())
			}
		}
		// Sort mac addresses alphabetically so we always return the same
		// serial number
		sort.Slice(macs, func(i, j int) bool {
			return macs[i] > macs[j]
		})
		if len(macs) > 0 {
			return macs[0]
		}
	}
	// Failure
	return ""
}

// Uptime returns the duration the machine has been switched on for
func Uptime() time.Duration {
	if info := sysinfo(); info != nil {
		return time.Second * time.Duration(info.Uptime)
	} else {
		return 0
	}
}

// Return load averages
func LoadAverage() (float64, float64, float64) {
	if info := sysinfo(); info != nil {
		return float64(info.Loads[0]) / float64(1<<16), float64(info.Loads[1]) / float64(1<<16), float64(info.Loads[2]) / float64(1<<16)
	} else {
		return 0, 0, 0
	}
}

func TemperatureZones() map[string]float32 {
	zones, err := ioutil.ReadDir(TEMPERATURE_PATH)
	if err != nil {
		return nil
	}
	temps := make(map[string]float32, len(zones))
	for _, zone := range zones {
		if data, err := ioutil.ReadFile(filepath.Join(TEMPERATURE_PATH, zone.Name(), "temp")); err != nil {
			continue
		} else if name, err := ioutil.ReadFile(filepath.Join(TEMPERATURE_PATH, zone.Name(), "type")); err != nil {
			continue
		} else {
			data_ := strings.TrimSpace(string(data))
			name_ := strings.TrimSpace(string(name))
			if value, err := strconv.ParseFloat(data_, 32); err != nil {
				continue
			} else {
				temps[name_] = float32(value) / float32(1000)
			}
		}
	}
	return temps
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func sysinfo() *syscall.Sysinfo_t {
	info := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&info); err != nil {
		return nil
	} else {
		return &info
	}
}
