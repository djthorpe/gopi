/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"net"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	dns "github.com/miekg/dns"
)

type service struct {
	Zone string
	TTL  time.Duration
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
	// Set TTL from PTR
	this.TTL = time.Duration(ptr.Hdr.Ttl) * time.Second
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
