// +build linux,rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import "net"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// SerialNumber returns the mac address of the hardware, if available
func SerialNumber() string {
	if ifaces, err := net.Interfaces(); err != nil {
		return ""
	} else if len(ifaces) == 0 {
		return ""
	} else {
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 {
				continue
			}
			if iface.Flags&net.FlagLoopback == 0 {
				continue
			}
			if iface.HardwareAddr != nil {
				return iface.HardwareAddr.String()
			}
		}
	}
	// Failure
	return ""
}
