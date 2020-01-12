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
	stop     gopi.Channel

	base.Unit
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
	}

	// Start background message processor
	if this.stop = this.listener.Subscribe(QUEUE_MESSAGES, this.EventHandler); this.stop == 0 {
		return nil, gopi.ErrInternalAppError
	}

	// Success
	return this, nil
}

func (this *discovery) Close() error {

	// Wait for EventHandler to end
	this.listener.Unsubscribe(this.stop)

	// Release resources
	this.listener = nil

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
	var lock sync.Mutex

	msg := new(dns.Msg)
	serviceRecords := map[string]gopi.RPCServiceRecord{}

	// Receive events in background
	receive := this.listener.Subscribe(QUEUE_RECORD, func(value interface{}) {
		lock.Lock()
		defer lock.Unlock()
		if evt, ok := value.(gopi.RPCEvent); ok == false || evt == nil {
			// No nothing
		} else if record := evt.Service(); record.Name == "" {
			// No nothing
		} else if evt.Type() != gopi.RPC_EVENT_SERVICE_RECORD {
			// No nothing
		} else if record.Service != service {
			// No nothing
		} else {
			serviceRecords[record.Name] = record
		}
	})

	// Perform the query and wait for cancellation
	msg.SetQuestion(service+"."+this.listener.Zone(), dns.TypePTR)
	msg.RecursionDesired = false
	err := this.listener.QueryAll(ctx, msg, QUERY_REPEAT)

	// Wait until unsubscribe
	this.listener.Unsubscribe(receive)

	// If error wasn't deadline exceeded, then return error
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
	var lock sync.Mutex

	msg := new(dns.Msg)
	serviceNames := make(map[string]bool)

	// Receive events in background
	fmt.Println("EnumerateServices LISTENER=", this.listener)
	receive := this.listener.Subscribe(QUEUE_NAME, func(value interface{}) {
		lock.Lock()
		defer lock.Unlock()

		if evt, ok := value.(gopi.RPCEvent); ok == false || evt == nil {
			// No nothing
		} else if record := evt.Service(); record.Name == "" {
			// No nothing
		} else if evt.Type() != gopi.RPC_EVENT_SERVICE_NAME {
			// No nothing
		} else {
			serviceNames[record.Name] = true
		}
	})

	// Perform the query and wait for cancellation
	msg.SetQuestion(DISCOVERY_SERVICE_QUERY+"."+this.listener.Zone(), dns.TypePTR)
	msg.RecursionDesired = false
	err := this.listener.QueryAll(ctx, msg, QUERY_REPEAT)

	// Wait until unsubscribe
	this.listener.Unsubscribe(receive)

	// If error wasn't deadline exceeded, then return error
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

func (this *discovery) EventHandler(value interface{}) {
	if msg, ok := value.(*dns.Msg); ok && msg == nil {
		return
	} else if len(msg.Answer) == 0 {
		return
	} else if err := this.ProcessAnswer(msg); err != nil {
		this.Log.Error(err.(error))
	}
}

func (this *discovery) ProcessAnswer(msg *dns.Msg) error {
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
		return nil
	}
	if strings.HasSuffix(service.Service, "._tcp") == false && strings.HasSuffix(service.Service, "._udp") == false {
		return nil
	}

	// Emit the event
	if service.Service == DISCOVERY_SERVICE_QUERY {
		this.listener.Emit(QUEUE_NAME, NewEvent(this, gopi.RPC_EVENT_SERVICE_NAME, service.RPCServiceRecord, service.TTL))
	} else {
		this.listener.Emit(QUEUE_RECORD, NewEvent(this, gopi.RPC_EVENT_SERVICE_RECORD, service.RPCServiceRecord, service.TTL))
	}

	// Success
	return nil
}
