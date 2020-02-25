/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"context"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ServiceDB2 struct {
	Discovery gopi.RPCServiceDiscovery
	Listener  ListenerIface
	Bus       gopi.Bus
}

type servicedb2 struct {
	discovery              gopi.RPCServiceDiscovery
	listener               ListenerIface
	bus                    gopi.Bus
	records, names, errors gopi.Channel

	Database
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (ServiceDB2) Name() string { return "gopi/mdns/servicedb2" }

func (config ServiceDB2) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(servicedb2)

	// Init
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if err := this.Init(config); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *servicedb2) Init(config ServiceDB2) error {

	// Check Discovery
	if config.Discovery == nil {
		return gopi.ErrBadParameter.WithPrefix("Discovery")
	} else {
		this.discovery = config.Discovery
	}

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

	// Subscribe to records
	this.records = this.listener.Subscribe(QUEUE_RECORD, func(value interface{}) {
		if r, ok := value.(gopi.RPCEvent); ok {
			this.Database.RegisterRecord(r)
		}
	})

	// Subscribe to names
	this.names = this.listener.Subscribe(QUEUE_NAME, func(value interface{}) {
		if r, ok := value.(gopi.RPCEvent); ok {
			this.Database.RegisterName(r)
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

func (this *servicedb2) Close() error {

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
	this.records = 0
	this.names = 0
	this.errors = 0
	this.listener = nil
	this.discovery = nil
	this.bus = nil

	// Return success
	return this.Unit.Close()
}

func (this *servicedb2) String() string {
	return "<" + this.Log.Name() + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCServiceDiscovery

// Lookup service instances by name
func (this *servicedb2) Lookup(ctx context.Context, service string) ([]gopi.RPCServiceRecord, error) {
	return this.discovery.Lookup(ctx, service)
}

// Return list of service names
func (this *servicedb2) EnumerateServices(ctx context.Context) ([]string, error) {
	return this.discovery.EnumerateServices(ctx)
}
