package mdns

import (
	"strings"

	"github.com/djthorpe/gopi/v3"
	"github.com/miekg/dns"
)

type services struct {
	services []*service
}

// Parse DNS message and capture service records
func NewServices(msg *dns.Msg, zone string) *services {
	this := new(services)
	sections := append(append(msg.Answer, msg.Ns...), msg.Extra...)
	for _, answer := range sections {
		switch rr := answer.(type) {
		case *dns.PTR:
			this.services = append(this.services, NewService(zone))
			this.services[0].SetPTR(rr)
		case *dns.SRV:
			if len(this.services) > 0 {
				this.services[0].SetSRV(rr.Target, rr.Port, rr.Priority)
			}
		case *dns.TXT:
			if len(this.services) > 0 {
				this.services[0].SetTXT(rr.Txt)
			}
		case *dns.A:
			if len(this.services) > 0 {
				this.services[0].SetA(rr.A)
			}
		case *dns.AAAA:
			if len(this.services) > 0 {
				this.services[0].SetAAAA(rr.AAAA)
			}
		}
	}

	return this
}

// Services returns all service records relevant for zone
func (this *services) Services() []gopi.ServiceRecord {
	result := make([]gopi.ServiceRecord, 0, len(this.services))
	for _, service := range this.services {
		if strings.HasSuffix(service.service, service.zone) {
			result = append(result, service)
		}
	}
	return result
}
