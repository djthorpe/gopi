/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package grpc

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi/sys/rpc"
	evt "github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ClientPool struct {
	SkipVerify bool
	SSL        bool
	Timeout    time.Duration
	Discovery  gopi.RPCServiceDiscovery
	Service    string
}

type clientpool struct {
	log        gopi.Logger
	skipverify bool
	ssl        bool
	timeout    time.Duration
	pubsub     *evt.PubSub
	discovery  gopi.RPCServiceDiscovery
	services   map[string]*servicetuple
	clients    map[string]gopi.RPCNewClientFunc
	done       chan struct{}
	lock       sync.Mutex
}

type servicetuple struct {
	cancelfunc context.CancelFunc
	deadline   time.Duration
	flags      gopi.RPCFlag
}

////////////////////////////////////////////////////////////////////////////////
// CONFIGURATION

const (
	DEFAULT_DISCOVERY_DEADLINE = 10 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ClientPool) Services() []string {
	if config.Service == "" {
		return nil
	} else {
		return strings.Split(config.Service, ",")
	}
}

func (config ClientPool) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.clientpool>Open{ services=%v SSL=%v skipverify=%v timeout=%v mdns=%v }", config.Services(), config.SSL, config.SkipVerify, config.Timeout, config.Discovery)
	this := new(clientpool)
	this.log = log
	this.skipverify = config.SkipVerify
	this.ssl = config.SSL
	this.timeout = config.Timeout
	this.pubsub = evt.NewPubSub(0)
	this.clients = make(map[string]gopi.RPCNewClientFunc)
	this.discovery = config.Discovery

	// Start the discovery in the background if it's available
	if this.discovery != nil {
		this.done = make(chan struct{})
		this.services = make(map[string]*servicetuple)
		// Check for rediscovery at half the rate (2 mins)
		go this.mdnsEventLoop(DEFAULT_DISCOVERY_DEADLINE * 2)
		// For each service, discovery is run every 1 minute
		for _, service := range config.Services() {
			if err := this.mdnsDiscovery(service, DEFAULT_DISCOVERY_DEADLINE, 0); err != nil {
				return nil, err
			}
		}
	}

	// Success
	return this, nil
}

func (this *clientpool) Close() error {
	this.log.Debug("<grpc.clientpool>Close{}")

	// Cancel all browsing
	for service, _ := range this.services {
		this.mdnsDoCancel(service)
	}

	// Stop discovery
	if this.done != nil {
		this.done <- gopi.DONE
		<-this.done
	}

	// Release resources
	this.pubsub.Close()
	this.pubsub = nil
	this.clients = nil
	this.discovery = nil
	this.done = nil
	this.services = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *clientpool) Connect(service *gopi.RPCServiceRecord, flags gopi.RPCFlag) (gopi.RPCClientConn, error) {
	this.log.Debug2("<grpc.clientpool>Connect{ service=%v flags=%v }", service, flags)

	// Determine the address
	if addr := addressFor(service, flags); addr == "" {
		return nil, gopi.ErrBadParameter
	} else if clientconn_, err := gopi.Open(ClientConn{
		Name:       service.Name,
		Addr:       addr + ":" + fmt.Sprint(service.Port),
		SSL:        this.ssl,
		SkipVerify: this.skipverify,
		Timeout:    this.timeout,
	}, this.log); err != nil {
		return nil, err
	} else if clientconn, ok := clientconn_.(gopi.RPCClientConn); ok == false {
		return nil, gopi.ErrOutOfOrder
	} else {
		// Do connection
		if err := clientconn.Connect(); err != nil {
			return nil, err
		}

		// Emit a connected event
		this.emit(rpc.NewEvent(clientconn, gopi.RPC_EVENT_CLIENT_CONNECTED, service))

		// Return success
		return clientconn, nil
	}
}

func (this *clientpool) Disconnect(conn gopi.RPCClientConn) error {
	this.log.Debug2("<grpc.clientpool>Disconnect{ conn=%v }", conn)

	// Emit a disconnect event
	this.emit(rpc.NewEvent(conn, gopi.RPC_EVENT_CLIENT_DISCONNECTED, nil))

	return conn.Disconnect()
}

////////////////////////////////////////////////////////////////////////////////
// CLIENTS

func (this *clientpool) RegisterClient(service string, callback gopi.RPCNewClientFunc) error {
	this.log.Debug2("<grpc.clientpool>RegisterClient{ service=%v func=%v }", service, callback)
	if service == "" || callback == nil {
		return gopi.ErrBadParameter
	} else if _, exists := this.clients[service]; exists {
		this.log.Debug("<rpc.clientpool>RegisterClient: Duplicate service: %v", service)
		return gopi.ErrBadParameter
	} else {
		this.clients[service] = callback
		return nil
	}
}

func (this *clientpool) NewClient(service string, conn gopi.RPCClientConn) gopi.RPCClient {
	this.log.Debug2("<grpc.clientpool>NewClient{ service=%v conn=%v }", service, conn)

	// Obtain the module with which to create a new client
	if callback, exists := this.clients[service]; exists == false {
		this.log.Debug("<grpc.clientpool>NewClient: Not Found: %v", service)
		return nil
	} else {
		return callback(conn)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBSUB

func (this *clientpool) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

func (this *clientpool) Unsubscribe(c <-chan gopi.Event) {
	this.pubsub.Unsubscribe(c)
}

func (this *clientpool) emit(evt gopi.Event) {
	this.pubsub.Emit(evt)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func addressFor(service *gopi.RPCServiceRecord, flags gopi.RPCFlag) string {
	if flags&gopi.RPC_FLAG_INET_UDP != 0 {
		// We don't support UDP connections
		return ""
	} else if flags&gopi.RPC_FLAG_INET_V6 != 0 {
		if len(service.IP6) == 0 {
			return ""
		} else {
			return service.IP6[0].String()
		}
	} else if flags&gopi.RPC_FLAG_INET_V4 != 0 {
		if len(service.IP4) == 0 {
			return ""
		} else {
			return service.IP4[0].String()
		}
	} else {
		return service.Host
	}
}

////////////////////////////////////////////////////////////////////////////////
// DISCOVERY METHODS

func (this *clientpool) mdnsEventLoop(delta time.Duration) {
	mdns_events := this.discovery.Subscribe()
	timer := time.NewTicker(delta)
FOR_LOOP:
	for {
		select {
		case <-this.done:
			break FOR_LOOP
		case evt := <-mdns_events:
			if rpc_evt, ok := evt.(gopi.RPCEvent); ok && rpc_evt != nil {
				this.emit(rpc.NewEvent(this, rpc_evt.Type(), rpc_evt.ServiceRecord()))
			}
		case <-timer.C:
			this.mdnsRediscovery()
		}
	}
	timer.Stop()
	this.discovery.Unsubscribe(mdns_events)
	close(this.done)
}

func (this *clientpool) mdnsRediscovery() {
	// Critical section
	this.lock.Lock()
	defer this.lock.Unlock()

	for service, tuple := range this.services {
		if tuple.cancelfunc == nil {
			this.log.Debug2("<grpc.clientpool>mdnsRediscovery{ service=%v }", service)
			if err := this.mdnsDiscovery(service, tuple.deadline, tuple.flags); err != nil {
				this.log.Error("<grpc.clientpool>mdnsRediscovery{ service=%v }: %v", service, err)
			}
		}
	}
}

func (this *clientpool) mdnsDoCancel(service string) {
	if tuple, exists := this.services[service]; exists {
		if tuple.cancelfunc != nil {
			tuple.cancelfunc()
		}
	}
}

func (this *clientpool) mdnsSetCancel(service string, cancel context.CancelFunc) *servicetuple {
	// Critical section
	this.lock.Lock()
	defer this.lock.Unlock()

	// Make a servicetuple if it doesn't exist yet
	if _, exists := this.services[service]; exists == false {
		this.services[service] = &servicetuple{
			cancelfunc: nil,
		}
	}

	// Set tuple parameters
	tuple := this.services[service]

	if cancel != nil {
		tuple.cancelfunc = cancel
	} else {
		tuple.cancelfunc = nil
	}

	// Return the service tuple
	return tuple
}

func (this *clientpool) mdnsDiscovery(service string, deadline time.Duration, flags gopi.RPCFlag) error {

	// If no discovery module then return error
	if this.discovery == nil {
		return gopi.ErrAppError
	}

	this.log.Debug2("<grpc.clientpool>mdnsDiscovery{ service=%v deadline=%v flags=%v }", service, deadline, flags)

	// Cancel current browse
	this.mdnsDoCancel(service)

	// Obtain service type
	if service_type, err := gopi.RPCServiceType(service, flags); err != nil {
		return err
	} else {
		// Fire up service discovery
		go func() {
			// Create context and tuple
			ctx, cancel := mdnsContextWithDeadline(deadline)
			tuple := this.mdnsSetCancel(service, cancel)

			// Set other tuple parameters
			tuple.deadline = deadline
			tuple.flags = flags

			// Start the browse and on end of the browse, return any error
			this.log.Debug2("<grpc.clientpool>mdnsDiscovery: STARTED for %v", service_type)
			if err := this.discovery.Browse(ctx, service_type); err != nil {
				this.log.Debug2("<grpc.clientpool>mdnsDiscovery: %v", err)
			}
			this.log.Debug2("<grpc.clientpool>mdnsDiscovery: STOPPED for %v", service_type)
			this.mdnsSetCancel(service, nil)
		}()
	}

	// Return success
	return nil
}

func mdnsContextWithDeadline(deadline time.Duration) (context.Context, context.CancelFunc) {
	if deadline == 0 {
		return context.WithCancel(context.Background())
	} else {
		return context.WithTimeout(context.Background(), deadline)
	}
}
