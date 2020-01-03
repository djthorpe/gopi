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
	"sort"
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
	Listener ListenerIface
}

type discovery struct {
	listener ListenerIface
	stop     chan struct{}

	sync.WaitGroup
	base.Unit
	Publisher
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Discovery) Name() string { return "gopi.mDNS.Discovery" }

func (config Discovery) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(discovery)

	// Check parameters
	if config.Listener == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("Listener")
	}

	// Init
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else {
		this.listener = config.Listener
		this.stop = make(chan struct{})
	}

	// Start background processor
	this.WaitGroup.Add(1)
	go this.ProcessMessages(this.stop)

	// Success
	return this, nil
}

func (this *discovery) Close() error {
	// Wait for process messages ends
	close(this.stop)
	this.Wait()

	// Release resources
	this.listener = nil
	this.stop = nil

	// Return success
	return this.Unit.Close()
}

func (this *discovery) String() string {
	return fmt.Sprintf("<gopi.mDNS.Discovery %v>", this.listener)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCServiceDiscovery

// Lookup service instances by name
func (this *discovery) Lookup(ctx context.Context, service string) ([]gopi.RPCServiceRecord, error) {
	var wait sync.WaitGroup

	// The message should be to lookup service by name
	msg := new(dns.Msg)
	msg.SetQuestion(service+"."+this.listener.Zone(), dns.TypePTR)
	msg.RecursionDesired = false

	// Set up map for service records
	serviceRecords := make(map[string]gopi.RPCServiceRecord)

	// Receive events in background
	bgctx, cancel := context.WithCancel(context.Background())
	wait.Add(1)
	go func(ctx context.Context) {
		names := this.Subscribe(QUEUE_RECORD, 0)
	FOR_LOOP:
		for {
			select {
			case evt := <-names:
				if evt_, ok := evt.(gopi.RPCEvent); ok {
					record := evt_.Service()
					if record.Name != "" && evt_.Type() == gopi.RPC_EVENT_SERVICE_RECORD {
						if record.Service == service {
							serviceRecords[record.Name] = record
						}
					}
				}
			case <-ctx.Done():
				break FOR_LOOP
			}
		}
		this.Unsubscribe(names)
		wait.Done()
	}(bgctx)

	// Perform the query and wait for cancellation
	err := this.listener.QueryAll(ctx, msg, QUERY_REPEAT)

	// Cancel background goroutine and wait until done
	cancel()
	wait.Wait()

	if err != nil && errors.Is(err, context.DeadlineExceeded) == false {
		// Error occurred other than context timeout
		return nil, err
	} else {
		// Gather service records
		serviceArr := make([]gopi.RPCServiceRecord, 0, len(serviceRecords))
		for _, record := range serviceRecords {
			serviceArr = append(serviceArr, record)
		}

		// Success
		return serviceArr, nil
	}
}

// Return list of service names
func (this *discovery) EnumerateServices(ctx context.Context) ([]string, error) {
	var wait sync.WaitGroup

	// The message should be to enumerate services
	msg := new(dns.Msg)
	msg.SetQuestion(DISCOVERY_SERVICE_QUERY+"."+this.listener.Zone(), dns.TypePTR)
	msg.RecursionDesired = false

	// Set up map for service names
	serviceNames := make(map[string]bool)

	// Receive events in background
	bgctx, cancel := context.WithCancel(context.Background())
	wait.Add(1)
	go func(ctx context.Context) {
		names := this.Subscribe(QUEUE_NAME, 0)
	FOR_LOOP:
		for {
			select {
			case evt := <-names:
				if evt_, ok := evt.(gopi.RPCEvent); ok {
					record := evt_.Service()
					if record.Name != "" && evt_.Type() == gopi.RPC_EVENT_SERVICE_NAME {
						serviceNames[record.Name] = true
					}
				}
			case <-ctx.Done():
				break FOR_LOOP
			}
		}
		this.Unsubscribe(names)
		wait.Done()
	}(bgctx)

	// Perform the query and wait for cancellation
	err := this.listener.QueryAll(ctx, msg, QUERY_REPEAT)

	// Cancel background goroutine and wait until done
	cancel()
	wait.Wait()

	if err != nil && errors.Is(err, context.DeadlineExceeded) == false {
		// Error occurred other than context timeout
		return nil, err
	} else {
		// Gather service names
		serviceNamesArr := make([]string, 0, len(serviceNames))
		for serviceName := range serviceNames {
			serviceNamesArr = append(serviceNamesArr, serviceName)
		}

		// Sort service names alphabetically
		sort.Strings(serviceNamesArr)

		// Success
		return serviceNamesArr, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND PROCESSOR

func (this *discovery) ProcessMessages(stop <-chan struct{}) {
	defer this.WaitGroup.Done()

	messages := this.listener.Subscribe(QUEUE_MESSAGES, 0)
	go func() {
		<-stop
		this.listener.Unsubscribe(messages)
	}()
FOR_LOOP:
	for {
		select {
		case msg := <-messages:
			if msg_, ok := msg.(*dns.Msg); ok && msg != nil {
				if len(msg_.Answer) > 0 {
					this.ProcessAnswer(msg_)
				}
			} else {
				break FOR_LOOP
			}
		}
	}
}

func (this *discovery) ProcessAnswer(msg *dns.Msg) {
	service := NewService(this.listener.Zone())
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
	// Emit the event
	if service.Service == DISCOVERY_SERVICE_QUERY {
		this.Emit(QUEUE_NAME, NewEvent(this, gopi.RPC_EVENT_SERVICE_NAME, service.RPCServiceRecord, service.TTL))
	} else {
		this.Emit(QUEUE_RECORD, NewEvent(this, gopi.RPC_EVENT_SERVICE_RECORD, service.RPCServiceRecord, service.TTL))
	}
}
