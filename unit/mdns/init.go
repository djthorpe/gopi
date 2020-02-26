/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {

	// gopi/mdns/listener
	gopi.UnitRegister(gopi.UnitConfig{
		Name: Listener{}.Name(),
		Config: func(app gopi.App) error {
			app.Flags().FlagString("mdns.domain", "local", "mDNS domain")
			app.Flags().FlagString("mdns.iface", "", "mDNS network interface")
			app.Flags().FlagBool("mdns.ip4", true, "mDNS uses IPv4")
			app.Flags().FlagBool("mdns.ip6", true, "mDNS uses IPv6")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			if iface, err := interfaceForString(app.Flags().GetString("mdns.iface", gopi.FLAG_NS_DEFAULT)); err != nil {
				return nil, err
			} else {
				flags := gopi.RPCFlag(0)
				if app.Flags().GetBool("mdns.ip4", gopi.FLAG_NS_DEFAULT) {
					flags |= gopi.RPC_FLAG_INET_V4
				}
				if app.Flags().GetBool("mdns.ip6", gopi.FLAG_NS_DEFAULT) {
					flags |= gopi.RPC_FLAG_INET_V6
				}
				return gopi.New(Listener{
					Domain:    app.Flags().GetString("mdns.domain", gopi.FLAG_NS_DEFAULT),
					Interface: iface,
					Flags:     flags,
				}, app.Log().Clone(Listener{}.Name()))
			}
		},
	})

	// gopi/mdns/discovery
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Discovery{}.Name(),
		Type:     gopi.UNIT_RPC_DISCOVERY,
		Requires: []string{Listener{}.Name()},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Discovery{
				Listener: app.UnitInstance(Listener{}.Name()).(ListenerIface),
			}, app.Log().Clone(Discovery{}.Name()))
		},
	})

	// gopi/mdns/servicedb
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     ServiceDB{}.Name(),
		Type:     gopi.UNIT_RPC_DISCOVERY,
		Pri:      1,
		Requires: []string{Listener{}.Name(), Discovery{}.Name(), "bus"},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(ServiceDB{
				Listener: app.UnitInstance(Listener{}.Name()).(ListenerIface),
				Bus:      app.Bus(),
			}, app.Log().Clone(ServiceDB{}.Name()))
		},
	})

	// gopi/mdns/register
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Register{}.Name(),
		Type:     gopi.UNIT_RPC_REGISTER,
		Requires: []string{Listener{}.Name()},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Register{
				Listener: app.UnitInstance(Listener{}.Name()).(ListenerIface),
			}, app.Log().Clone(Register{}.Name()))
		},
	})
}

func interfaceForString(name string) (net.Interface, error) {
	if name == "" {
		return net.Interface{}, nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, err
	}
	names := ""
	for _, iface := range ifaces {
		if iface.Name == name {
			return iface, nil
		}
		names += strconv.Quote(iface.Name) + ","
	}
	return net.Interface{}, fmt.Errorf("Invalid -mdns.iface flag (values: %v)", strings.Trim(names, ","))
}
