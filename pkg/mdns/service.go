package mdns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

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

func NewService(zone string) *service {
	this := new(service)
	this.zone = zone
	return this
}

func (this *service) Service() string {
	return strings.TrimSuffix(this.service, this.zone)
}

func (this *service) Name() string {
	return strings.TrimSuffix(this.name, this.zone)
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

func (this *service) String() string {
	str := "<dns.servicerecord"
	if this.service != "" {
		str += fmt.Sprintf(" service=%q", strings.TrimSuffix(this.service, this.zone))
	}
	if this.name != "" {
		str += fmt.Sprintf(" name=%q", strings.TrimSuffix(this.name, this.zone))
	}
	if len(this.host) > 0 {
		str += fmt.Sprintf(" host=%v", this.host)
	}
	if len(this.a) > 0 {
		str += fmt.Sprintf(" a=%v", this.a)
	}
	if len(this.aaaa) > 0 {
		str += fmt.Sprintf(" aaaa=%v", this.aaaa)
	}
	if len(this.txt) > 0 {
		str += fmt.Sprintf(" txt=%v", this.txt)
	}
	if this.ttl != 0 {
		str += fmt.Sprintf(" ttl=%v", this.ttl)
	}
	return str + ">"
}
