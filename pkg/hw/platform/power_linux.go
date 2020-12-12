// +build linux
// +build !rpi

package platform

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Platform) SetPowerState() error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	var s []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return err
	}

	for _, v := range s {
		fmt.Println(v)
	}

	// Return success
	return nil
}
