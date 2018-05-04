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
	"net"
	"strconv"
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
	records    map[string]*gopi.RPCServiceRecord
	done       chan struct{}
	lock       sync.Mutex
}

type servicetuple struct {
	cancelfunc context.CancelFunc
	deadline   time.Duration
	flags      gopi.RPCFlag
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *clientpool) String() string {
	return fmt.Sprintf("<grpc.clientpool>{ ssl=%v skipverify=%v timeout=%v clients=%v services=%v records=%v", this.ssl, this.skipverify, this.timeout, this.clients, this.services, this.records)
}

////////////////////////////////////////////////////////////////////////////////
// CONFIGURATION

const (
	DEFAULT_DISCOVERY_DEADLINE = 5 * time.Second
	DEFAULT_DISCOVERY_DELTA    = 60 * time.Second
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
	this.records = make(map[string]*gopi.RPCServiceRecord)

	// Start the discovery in the background if it's available
	if this.discovery != nil {
		this.done = make(chan struct{})
		this.services = make(map[string]*servicetuple)
		// Check for rediscovery every 1 min
		go this.mdnsEventLoop(DEFAULT_DISCOVERY_DELTA)
		// For each service, discovery is run for 5 seconds
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
	this.records = nil

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
	} else if clientconn, ok := clientconn_.(*clientconn); ok == false {
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

	// Emit a disconnect event - event happens just before disconnect occurs
	this.emit(rpc.NewEvent(conn, gopi.RPC_EVENT_CLIENT_DISCONNECTED, nil))

	return conn.(*clientconn).Disconnect()
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
	// Emit happens in goroutine
	go this.pubsub.Emit(evt)
}

////////////////////////////////////////////////////////////////////////////////
// LOOKUP

func (this *clientpool) Lookup(ctx context.Context, name, addr string, max int) ([]*gopi.RPCServiceRecord, error) {
	this.log.Debug2("<grpc.clientpool>Lookup{ name='%v' addr='%v' max='%v' }", name, addr, max)

	// Make a buffered channel of all service records and put them in
	records := make(chan *gopi.RPCServiceRecord, len(this.records)+2)
	matched := make([]*gopi.RPCServiceRecord, 0, max)

	// Queue up the records we know about
	for _, record := range this.records {
		records <- record
	}

	// Subscribe to receive events
	pool_events := this.Subscribe()

	// The loop continues until the context is done, all existing
	// records are consumed and all new records are matched up to
	// the 'max' number of records (or unlimited if max is zero)
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		case record := <-records:
			if record != nil && lookupMatch(name, addr, record) {
				matched = append(matched, record)
				if max != 0 && len(matched) >= max {
					break FOR_LOOP
				}
			}
		case event := <-pool_events:
			if event != nil {
				if rpc_event := event.(gopi.RPCEvent); rpc_event.Type() == gopi.RPC_EVENT_SERVICE_RECORD {
					// Emit the service record for matching
					records <- rpc_event.ServiceRecord()
				}
			}
		}
	}

	// Cleanup
	this.Unsubscribe(pool_events)
	close(records)

	// If we have reached max matched records, then return them
	// without an error or else return 'DeadlineExceeded'
	if len(matched) == 0 {
		return nil, gopi.ErrDeadlineExceeded
	} else if len(matched) < max {
		return matched, gopi.ErrDeadlineExceeded
	} else {
		return matched, nil
	}
}

func lookupMatch(name, addr string, record *gopi.RPCServiceRecord) bool {
	if name != "" && name == record.Name {
		return true
	}
	if addr != "" {
		host, port := lookupMatchSplitAddr(addr)
		if host == "" && port == 0 {
			// Parse error occured here
			return false
		}
		if host != "" {
			// We return false if any host doesn't match
			if lookupMatchHost(host, record) == false && lookupMatchIP(host, record) == false {
				return false
			}
		}
		if port != 0 && port != record.Port {
			return false
		}
		return true
	}
	// No conditions so match
	return true
}

func lookupMatchSplitAddr(addr string) (string, uint) {
	if host, port_, err := net.SplitHostPort(addr); err != nil {
		if strings.HasSuffix(err.Error(), "missing port in address") {
			return addr, 0
		} else {
			return "", 0
		}
	} else if port, err := strconv.ParseUint(strings.TrimPrefix(port_, ":"), 10, 32); err != nil {
		return addr, 0
	} else {
		return host, uint(port)
	}
}

func lookupMatchHost(addr string, record *gopi.RPCServiceRecord) bool {
	// Fully-qualify address
	if strings.HasSuffix(addr, ".") == false {
		addr += "."
	}
	return strings.ToLower(addr) == strings.ToLower(record.Host)
}

func lookupMatchIP(addr string, record *gopi.RPCServiceRecord) bool {
	if ip := net.ParseIP(addr); ip != nil {
		for _, ip4 := range record.IP4 {
			if ip4.Equal(ip) {
				return true
			}
		}
		for _, ip6 := range record.IP6 {
			if ip6.Equal(ip) {
				return true
			}
		}
	}
	return false
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
			if evt != nil {
				this.mdnsHandleEvent(evt.(gopi.RPCEvent))
			}
		case <-timer.C:
			this.mdnsRediscovery()
		}
	}
	timer.Stop()
	this.discovery.Unsubscribe(mdns_events)
	close(this.done)
}

func (this *clientpool) mdnsHandleEvent(evt gopi.RPCEvent) {
	// Check for event being a service record
	if evt.Type() != gopi.RPC_EVENT_SERVICE_RECORD {
		return
	}
	// Now create a hash for the service record
	new_sr := evt.ServiceRecord()
	if hash := mdnsRecordHash(new_sr); hash == "" {
		return
	} else if sr, exists := this.records[hash]; exists == false {
		this.records[hash] = new_sr
		this.emit(rpc.NewEvent(this, evt.Type(), new_sr))
	} else if mdnsRecordEquals(sr, new_sr) {
		// Do nothing
	} else {
		// Copy over the details from the new_sr
		sr.Host = new_sr.Host
		sr.IP4 = new_sr.IP4
		sr.IP6 = new_sr.IP6
		sr.Name = new_sr.Name
		sr.Port = new_sr.Port
		sr.Text = new_sr.Text
		sr.TTL = new_sr.TTL
		sr.Type = new_sr.Type
		this.emit(rpc.NewEvent(this, evt.Type(), sr))

		// If TTL is zero then we delete the record
		if sr.TTL == 0 {
			delete(this.records, hash)
		}
	}
}

func mdnsRecordHash(record *gopi.RPCServiceRecord) string {
	// The hash is a combo of host.port.type
	if record.Host == "" || record.Port == 0 || record.Type == "" {
		return ""
	}
	return strings.Join([]string{record.Host, fmt.Sprint(record.Port), record.Type}, " ")
}

func mdnsRecordEquals(a, b *gopi.RPCServiceRecord) bool {
	if a.Host != b.Host {
		return false
	}
	if a.Type != b.Type {
		return false
	}
	if a.Port != b.Port {
		return false
	}
	if a.TTL != b.TTL {
		return false
	}
	if len(a.Text) != len(b.Text) {
		return false
	}
	for i := range a.Text {
		if a.Text[i] != b.Text[i] {
			return false
		}
	}
	return true
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
			this.log.Debug2("<grpc.clientpool>mdnsDiscovery: Start: %v", service_type)
			if err := this.discovery.Browse(ctx, service_type); err != nil {
				this.log.Debug2("<grpc.clientpool>mdnsDiscovery: %v", err)
			}
			this.log.Debug2("<grpc.clientpool>mdnsDiscovery: Stop: %v", service_type)
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
