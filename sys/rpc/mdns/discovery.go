/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"context"
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi/sys/rpc"
	evt "github.com/djthorpe/gopi/util/event"
	"github.com/djthorpe/zeroconf"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// The configuration
type Config struct {
	Domain string
}

// The driver for the logging
type driver struct {
	log      gopi.Logger
	domain   string
	servers  []*zeroconf.Server
	resolver *zeroconf.Resolver
	pubsub   *evt.PubSub
}

///////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	MDNS_DEFAULT_DOMAIN = "local."
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create discovery object
func (config Config) Open(log gopi.Logger) (gopi.Driver, error) {

	this := new(driver)
	this.log = log
	if config.Domain == "" {
		this.domain = MDNS_DEFAULT_DOMAIN
	} else {
		this.domain = config.Domain
	}

	log.Debug("sys.rpc.mDNS.Open{ domain=%v }", this.domain)

	this.servers = make([]*zeroconf.Server, 0, 1)

	if resolver, err := zeroconf.NewResolver(); err != nil {
		return nil, err
	} else {
		this.resolver = resolver
	}

	// Publish/Subscribe
	this.pubsub = evt.NewPubSub(0)

	// success
	return this, nil
}

// Close discovery object
func (this *driver) Close() error {
	this.log.Debug("sys.rpc.mDNS.Close{ domain=%v }", this.domain)

	// Close servers
	for _, server := range this.servers {
		server.Shutdown()
	}

	// Unsubscribe
	this.pubsub.Close()

	// Empty methods
	this.servers = nil
	this.pubsub = nil
	this.resolver = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// DRIVER INTERFACE METHODS

// Register a service and announce the service when queries occur
func (this *driver) Register(service *gopi.RPCServiceRecord) error {
	if server, err := zeroconf.Register(service.Name, service.Type, this.domain, int(service.Port), service.Text, nil); err != nil {
		return err
	} else {
		this.servers = append(this.servers, server)
		return nil
	}
}

// Browse will find service entries, will block until ctx timeout
// or cancel
func (this *driver) Browse(ctx context.Context, serviceType string) error {
	entries := make(chan *zeroconf.ServiceEntry)
	if err := this.resolver.Browse(ctx, serviceType, this.domain, entries); err != nil {
		return err
	} else {
		for entry := range entries {
			this.Emit(&gopi.RPCServiceRecord{
				Name: entry.Instance,
				Type: entry.Service,
				Port: uint(entry.Port),
				Text: entry.Text,
				Host: entry.HostName,
				IP4:  entry.AddrIPv4,
				IP6:  entry.AddrIPv6,
				TTL:  time.Duration(entry.TTL) * time.Second,
			})
		}
		return nil
	}
}

func (this *driver) DefaultServiceType(network string) string {
	return "_gopi._" + network
}

////////////////////////////////////////////////////////////////////////////////
// PUBSUB

// Subscribe to events emitted
func (this *driver) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

// Unsubscribe from events emitted
func (this *driver) Unsubscribe(subscriber <-chan gopi.Event) {
	this.pubsub.Unsubscribe(subscriber)
}

// Emit an event
func (this *driver) Emit(record *gopi.RPCServiceRecord) {
	this.pubsub.Emit(rpc.NewEvent(this, gopi.RPC_EVENT_SERVICE_RECORD, record))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	return fmt.Sprintf("sys.mdns{ domain=\"%v\" registrations=%v }", this.domain, "TODO")
}
