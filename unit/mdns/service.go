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

	"github.com/miekg/dns"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

type service struct {
	Zone string
	gopi.RPCServiceRecord
}

func NewService(zone string) *service {
	return &service{
		Zone: zone,
	}
}

func (this *service) SetPTR(ptr *dns.PTR) {
	this.Service = strings.TrimSuffix(ptr.Hdr.Name, "."+this.Zone)
	this.Name = strings.TrimSuffix(ptr.Ptr, "."+this.Zone)
	if this.Service != DISCOVERY_SERVICE_QUERY {
		this.Name = strings.TrimSuffix(this.Name, "."+this.Service)
	}
	// Unquote the name if possible
	if name, err := Unquote(this.Name); err == nil {
		this.Name = name
	}
}

func (this *service) SetSRV(host string, port, priority uint16) {
	this.Host = host
	this.Port = port
}

func (this *service) SetTXT(txt []string) {
	this.Txt = txt
}

func (this *service) SetA(addr net.IP) {
	if this.Addrs == nil {
		this.Addrs = make([]net.IP, 0, 1)
	}
	this.Addrs = append(this.Addrs, addr)
}

func (this *service) SetAAAA(addr net.IP) {
	if this.Addrs == nil {
		this.Addrs = make([]net.IP, 0, 1)
	}
	this.Addrs = append(this.Addrs, addr)
}

func (this *service) String() string {
	str := "<service name=" + strconv.Quote(this.Name)
	if this.Service != "" {
		str += " service=" + strconv.Quote(this.Service)
	}
	if this.Subtype != "" {
		str += " subtype=" + strconv.Quote(this.Subtype)
	}
	if this.Host != "" {
		str += " host=" + strconv.Quote(this.Host)
	}
	if this.Port != 0 {
		str += " port=" + fmt.Sprint(this.Port)
	}
	if len(this.Addrs) > 0 {
		str += " addrs="
		for _, addr := range this.Addrs {
			str += addr.String() + ","
		}
		str = strings.TrimSuffix(str, ",")
	}
	if len(this.Txt) > 0 {
		str += " txt=" + fmt.Sprint(this.Txt)
	}
	str += ">"
	return str
}
