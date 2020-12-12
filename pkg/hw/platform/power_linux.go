// +build linux
// +build !rpi

package platform

import (
	gopi "github.com/djthorpe/gopi/v3"
	dbus "github.com/godbus/dbus/v5"
)

const (
	dbusNode = "org.freedesktop.login1"
	dbusPath = "/org/freedesktop/login1"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Platform) SetPowerState() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	askForAuth := false
	if obj := conn.Object(dbusNode, dbusPath); obj == nil {
		return gopi.ErrNotFound.WithPrefix(dbusNode)
	} else if result := obj.Call("org.freedesktop.login1.Manager.Reboot", 0, askForAuth); result.Err != nil {
		return result.Err
	}

	// Return success
	return nil
}
