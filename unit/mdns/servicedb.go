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
	"math/rand"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	dns "github.com/miekg/dns"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ServiceDB struct {
	Listener ListenerIface
	Bus      gopi.Bus
}

type servicedb struct {
	listener               ListenerIface
	bus                    gopi.Bus
	records, names, errors gopi.Channel
	stop                   chan struct{}
	lookup                 map[string]bool

	Database
	base.Unit
	sync.WaitGroup
	sync.Mutex // mutex for lookup
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (ServiceDB) Name() string { return "gopi/mdns/servicedb" }

func (config ServiceDB) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(servicedb)

	// Init
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if err := this.Init(config); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *servicedb) Init(config ServiceDB) error {
	// Check Listener
	if config.Listener == nil {
		return gopi.ErrBadParameter.WithPrefix("Listener")
	} else {
		this.listener = config.Listener
	}
	// Check Bus
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("Bus")
	} else {
		this.bus = config.Bus
	}

	// Init database
	this.Database.Init()

	// Initiate lookup
	this.stop = make(chan struct{})
	this.lookup = make(map[string]bool)
	go this.backgroundTask(this.stop)

	// Subscribe to records
	this.records = this.listener.Subscribe(QUEUE_RECORD, func(value interface{}) {
		if r, ok := value.(*event); ok {
			if evt := this.Database.RegisterRecord(r); evt != nil {
				this.bus.Emit(evt)
			}
		}
	})

	// Subscribe to names
	this.names = this.listener.Subscribe(QUEUE_NAME, func(value interface{}) {
		if r, ok := value.(*event); ok {
			if evt := this.Database.RegisterName(r); evt != nil {
				if evt.Type() == gopi.RPC_EVENT_SERVICE_NAME {
					this.addLookup(evt.Service().Name)
				}
			}
		}
	})

	// Subscribe to errors
	this.errors = this.listener.Subscribe(QUEUE_ERRORS, func(value interface{}) {
		if err, ok := value.(error); ok {
			this.Log.Error(err)
		}
	})

	// Return success
	return nil
}

func (this *servicedb) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Stop background task and wait until exit
	close(this.stop)
	this.WaitGroup.Wait()

	// Unsubscribe from listener queues
	if this.records != 0 {
		this.listener.Unsubscribe(this.records)
	}
	if this.names != 0 {
		this.listener.Unsubscribe(this.names)
	}
	if this.errors != 0 {
		this.listener.Unsubscribe(this.errors)
	}

	// Close database
	this.Database.Close()

	// Release resources
	this.lookup = nil
	this.records = 0
	this.names = 0
	this.errors = 0
	this.listener = nil
	this.bus = nil

	// Return success
	return this.Unit.Close()
}

func (this *servicedb) String() string {
	return "<" + this.Log.Name() + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCServiceDiscovery

// Return service instances by name
func (this *servicedb) Lookup(ctx context.Context, service string) ([]gopi.RPCServiceRecord, error) {
	// Perform the query
	msg := new(dns.Msg)
	msg.SetQuestion(service+"."+this.listener.Zone(), dns.TypePTR)
	msg.RecursionDesired = false
	if err := this.listener.QueryAll(ctx, msg, QUERY_REPEAT); err != nil && err != context.DeadlineExceeded {
		return nil, err
	} else {
		return this.Database.Records(service), nil
	}
}

// Return list of service names
func (this *servicedb) EnumerateServices(ctx context.Context) ([]string, error) {
	// Perform the query
	msg := new(dns.Msg)
	msg.SetQuestion(DISCOVERY_SERVICE_QUERY+"."+this.listener.Zone(), dns.TypePTR)
	msg.RecursionDesired = false
	if err := this.listener.QueryAll(ctx, msg, QUERY_REPEAT); err != nil && err != context.DeadlineExceeded {
		return nil, err
	} else {
		return this.Database.Names(), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

// addLookup adds a name to the list of names to lookup service records for
func (this *servicedb) addLookup(name string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.lookup[name] = true
}

// popLookup removes an entry from the list of names
func (this *servicedb) popLookup() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if len(this.lookup) > 0 {
		idx := rand.Int() % len(this.lookup)
		for name := range this.lookup {
			if idx == 0 {
				delete(this.lookup, name)
				return name
			} else {
				idx--
			}
		}
	}

	// Nothing to pop
	return ""
}

func (this *servicedb) backgroundTask(stop <-chan struct{}) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()
	lookup := time.NewTicker(DISCOVERY_LOOKUP_DELTA)
	ctx, cancel := context.WithTimeout(context.Background(), DISCOVERY_LOOKUP_DELTA)
FOR_LOOP:
	for {
		select {
		case <-lookup.C:
			if service := this.popLookup(); service != "" {
				if err := this.backgroundLookup(ctx, service); err != nil && err != context.DeadlineExceeded {
					this.Log.Error(fmt.Errorf("backgroundTask: %w", err))
				}
				ctx, cancel = context.WithTimeout(context.Background(), DISCOVERY_LOOKUP_DELTA)
			}
		case <-stop:
			cancel()
			lookup.Stop()
			break FOR_LOOP
		}
	}
}

func (this *servicedb) backgroundLookup(ctx context.Context, service string) error {
	// Perform the query and wait for cancellation
	msg := new(dns.Msg)
	msg.SetQuestion(service+"."+this.listener.Zone(), dns.TypePTR)
	msg.RecursionDesired = false
	return this.listener.QueryAll(ctx, msg, QUERY_REPEAT)
}
