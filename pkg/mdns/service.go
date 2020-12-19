package mdns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type service struct {
	service string
	zone    string
	name    string
	ttl     time.Duration
	host    []target
	a       []net.IP
	aaaa    []net.IP
	txt     []string
}

type target struct {
	target   string
	port     uint16
	priority uint16
}

///////////////////////////////////////////////////////////////////////////////
// INIT

func NewService(zone string) *service {
	this := new(service)
	this.zone = zone
	return this
}

///////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this *service) Service() string {
	return strings.TrimSuffix(this.service, this.zone)
}

func (this *service) Name() string {
	name := ""
	if this.Service() == fqn(queryServices) {
		name = strings.TrimSuffix(this.name, this.zone)
	} else {
		name = strings.TrimSuffix(this.name, this.service)
		if name_, err := Unquote(unfqn(name)); err != nil {
		} else {
			name = name_
		}
	}
	return name
}

func (this *service) Instance() string {
	return fqn(this.name)
}

func (this *service) HostPort() []string {
	hostport := make([]string, 0, len(this.host))
	for _, target := range this.host {
		hostport = append(hostport, net.JoinHostPort(target.target, fmt.Sprint(target.port)))
	}
	return hostport
}

func (this *service) Addrs() []net.IP {
	addrs := []net.IP{}
	addrs = append(addrs, this.a...)
	addrs = append(addrs, this.aaaa...)
	return addrs
}

func (this *service) Txt() []string {
	return this.txt
}

///////////////////////////////////////////////////////////////////////////////
// SET PROPERTIES

func (this *service) SetPTR(ptr *dns.PTR) {
	this.service = ptr.Hdr.Name
	this.name = ptr.Ptr
	this.ttl = time.Duration(ptr.Hdr.Ttl) * time.Second
}

func (this *service) SetSRV(host string, port uint16, priority uint16) {
	this.host = append(this.host, target{host, port, priority})
}

func (this *service) SetTXT(txt []string) {
	this.txt = txt
}

func (this *service) SetA(ip net.IP) {
	this.a = append(this.a, ip)
}

func (this *service) SetAAAA(ip net.IP) {
	this.aaaa = append(this.aaaa, ip)
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	str := "<dns.servicerecord"
	if service := this.Service(); service != "" {
		str += fmt.Sprintf(" service=%q", service)
	}
	if name := this.Name(); name != "" {
		str += fmt.Sprintf(" name=%q", name)
	}
	if hostport := this.HostPort(); len(hostport) > 0 {
		str += fmt.Sprintf(" host=%v", hostport)
	}
	if ips := this.Addrs(); len(ips) > 0 {
		str += fmt.Sprintf(" addrs=%v", ips)
	}
	if txt := this.Txt(); len(this.txt) > 0 {
		str += fmt.Sprintf(" txt=%v", txt)
	}
	if this.ttl != 0 {
		str += fmt.Sprintf(" ttl=%v", this.ttl)
	}
	return str + ">"
}
