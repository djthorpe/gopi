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
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/mdns/discovery",
		Type: gopi.UNIT_RPC_DISCOVERY,
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
				return gopi.New(Discovery{
					Domain:    app.Flags().GetString("mdns.domain", gopi.FLAG_NS_DEFAULT),
					Interface: iface,
					Flags:     flags,
				}, app.Log().Clone("gopi/mdns/discovery"))
			}
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
