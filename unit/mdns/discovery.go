/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	dns "github.com/miekg/dns"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Discovery struct {
	Domain    string
	Interface net.Interface
	Flags     gopi.RPCFlag
}

type discovery struct {
	errors   chan error
	messages chan *dns.Msg

	Listener
	base.Unit
	sync.Mutex
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DISCOVERY_SERVICE_QUERY = "_services._dns-sd._udp"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Discovery) Name() string { return "gopi.mDNS.Discovery" }

func (config Discovery) FQDomain() string {
	if config.Domain == "" {
		return ""
	} else {
		return strings.Trim(config.Domain, ".") + "."
	}
}

func (config Discovery) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(discovery)
	this.errors = make(chan error)
	this.messages = make(chan *dns.Msg)

	// Init
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if err := this.Listener.Init(config, this.errors, this.messages); err != nil {
		return nil, err
	}

	// Start background processor
	this.WaitGroup.Add(1)
	go this.ProcessMessages()

	return this, nil
}

func (this *discovery) Close() error {
	// Close listener
	if err := this.Listener.Destroy(); err != nil {
		return err
	}

	// Release resources
	close(this.errors)
	close(this.messages)

	// Wait for process messages ends
	this.Wait()

	// Return success
	return nil
}

func (this *discovery) String() string {
	return fmt.Sprintf("<gopi.mDNS.Discovery %v>", this.Listener)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCServiceDiscovery

// Lookup service instances by name
func (this *discovery) Lookup(ctx context.Context, service string) ([]gopi.RPCServiceRecord, error) {
	return nil, gopi.ErrNotImplemented
}

// Return list of service names
func (this *discovery) EnumerateServices(ctx context.Context) ([]string, error) {
	return nil, gopi.ErrNotImplemented
}

// Return all cached service instances for a service name
func (this *discovery) ServiceInstances(service string) []gopi.RPCServiceRecord {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND PROCESSOR

func (this *discovery) ProcessMessages() {
	this.Log.Debug("ProcessMessages started")
FOR_LOOP:
	for {
		select {
		case err := <-this.errors:
			if err == nil {
				break FOR_LOOP
			} else {
				this.Log.Error(err)
			}
		case msg := <-this.messages:
			if msg == nil {
				break FOR_LOOP
			} else if len(msg.Answer) > 0 {
				this.ProcessAnswer(msg)
			}
		}
	}
	this.Log.Debug("ProcessMessages finished")
	this.Done()
}

func (this *discovery) ProcessAnswer(msg *dns.Msg) {
	service := NewService(this.Listener.domain)
	sections := append(append(msg.Answer, msg.Ns...), msg.Extra...)
	for _, answer := range sections {
		switch rr := answer.(type) {
		case *dns.PTR:
			service.SetPTR(rr)
		case *dns.SRV:
			service.SetSRV(rr.Target, rr.Port, rr.Priority)
		case *dns.TXT:
			service.SetTXT(rr.Txt)
		case *dns.A:
			service.SetA(rr.A)
		case *dns.AAAA:
			service.SetAAAA(rr.AAAA)
		default:
			//fmt.Println("OTHER", rr)
		}
	}
	// Ignore any service without a name, or where it doesn't end with _udp or _tcp
	if service.Name == "" {
		return
	}
	if strings.HasSuffix(service.Service, "._tcp") == false && strings.HasSuffix(service.Service, "._udp") == false {
		return
	}
	if service.Service == DISCOVERY_SERVICE_QUERY {
		fmt.Println("GOT NAME", service.Name)
	} else {
		fmt.Println("GOT SERVICE", service)
	}
}
