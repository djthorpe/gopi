/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"context"
	"errors"
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
	Bus       gopi.Bus
}

type discovery struct {
	errors   chan error
	messages chan *dns.Msg
	bus      gopi.Bus

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
	} else if config.Bus == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("Bus")
	} else {
		this.bus = config.Bus
	}

	// Start background processor
	this.WaitGroup.Add(1)
	go this.ProcessMessages()

	// Success
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

	// Release resources
	this.bus = nil
	this.errors = nil
	this.messages = nil

	// Return success
	return this.Unit.Close()
}

func (this *discovery) String() string {
	return fmt.Sprintf("<gopi.mDNS.Discovery %v>", this.Listener)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCServiceDiscovery

// Lookup service instances by name
func (this *discovery) Lookup(ctx context.Context, service string) ([]gopi.RPCServiceRecord, error) {
	// The message should be to lookup service by name
	msg := new(dns.Msg)
	msg.SetQuestion(service+"."+this.Listener.domain, dns.TypePTR)
	msg.RecursionDesired = false

	// Set a handler which reads RPCEvents, remove when done
	var lock sync.Mutex
	services := make(map[string]gopi.RPCServiceRecord)

	// TODO handler :=
	this.bus.NewHandler("gopi.RPCEvent", func(_ context.Context, evt gopi.Event) {
		evt_ := evt.(gopi.RPCEvent)
		if evt_.Type() == gopi.RPC_EVENT_SERVICE_RECORD {
			if service := evt_.Service(); service.Name != "" {
				lock.Lock()
				defer lock.Unlock()
				services[service.Name] = service
			}
		}
	})
	// TODO
	//defer this.bus.Cancel(handler)

	// Perform the query and wait for cancellation
	if err := this.QueryAll(ctx, msg, 2); err != nil && errors.Is(err, context.DeadlineExceeded) == false {
		return nil, err
	}

	// Return service records
	records := make([]gopi.RPCServiceRecord, 0, len(services))
	for _, service := range services {
		records = append(records, service)
	}
	return records, nil
}

// Return list of service names
func (this *discovery) EnumerateServices(ctx context.Context) ([]string, error) {
	// The message should be to enumerate services
	msg := new(dns.Msg)
	msg.SetQuestion(DISCOVERY_SERVICE_QUERY+"."+this.Listener.domain, dns.TypePTR)
	msg.RecursionDesired = false

	// Set a handler which reads RPCEvents, remove when done
	var lock sync.Mutex
	services := make(map[string]bool)

	// TODO handler :=
	this.bus.NewHandler("gopi.RPCEvent", func(_ context.Context, evt gopi.Event) {
		evt_ := evt.(gopi.RPCEvent)
		if evt_.Type() == gopi.RPC_EVENT_SERVICE_NAME {
			if service := evt_.Service(); service.Name != "" {
				lock.Lock()
				defer lock.Unlock()
				services[service.Name] = true
			}
		}
	})
	// TODO
	//defer this.bus.Cancel(handler)

	// Perform the query and wait for cancellation
	if err := this.QueryAll(ctx, msg, 2); err != nil && errors.Is(err, context.DeadlineExceeded) == false {
		return nil, err
	} else {
		// Create array of services
		serviceNames := make([]string, 0, len(services))
		for service := range services {
			serviceNames = append(serviceNames, service)
		}
		// Success
		return serviceNames, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND PROCESSOR

func (this *discovery) ProcessMessages() {
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
	this.WaitGroup.Done()
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
		this.bus.Emit(NewEvent(this, gopi.RPC_EVENT_SERVICE_NAME, service.RPCServiceRecord))
	} else {
		this.bus.Emit(NewEvent(this, gopi.RPC_EVENT_SERVICE_RECORD, service.RPCServiceRecord))
	}
}
