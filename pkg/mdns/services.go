package mdns

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

type services struct {
	zone string
}

func NewServices(msg *dns.Msg, zone string) *services {
	this := new(services)
	this.zone = zone
	sections := append(append(msg.Answer, msg.Ns...), msg.Extra...)
	for _, answer := range sections {
		switch rr := answer.(type) {
		case *dns.PTR:
			this.SetPTR(rr)
		case *dns.SRV:
			this.SetSRV(rr.Target, rr.Port, rr.Priority)
		case *dns.TXT:
			this.SetTXT(rr.Txt)
		case *dns.A:
			this.SetA(rr.A)
		case *dns.AAAA:
			this.SetAAAA(rr.AAAA)
		}
	}

	return this
}

func (this *services) SetPTR(ptr *dns.PTR) {
	fmt.Println("PTR=", ptr)
}

func (this *services) SetSRV(target string, port uint16, priority uint16) {
	fmt.Println("SRV=", target, port, priority)
}

func (this *services) SetTXT(txt []string) {
	fmt.Println("TXT=", txt)
}

func (this *services) SetA(ip net.IP) {
	fmt.Println("A=", ip)
}

func (this *services) SetAAAA(ip net.IP) {
	fmt.Println("AAAA=", ip)
}
