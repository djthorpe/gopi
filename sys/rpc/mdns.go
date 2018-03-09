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

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/third_party/zeroconf"
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
}

///////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	MDNS_DEFAULT_DOMAIN = "local."
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open a logger
func (config Config) Open(log gopi.Logger) (gopi.Driver, error) {

	this := new(driver)
	this.log = log
	if config.Domain == "" {
		this.domain = MDNS_DEFAULT_DOMAIN
	} else {
		this.domain = config.Domain
	}
	this.servers = make([]*zeroconf.Server, 0, 1)

	if resolver, err := zeroconf.NewResolver(); err != nil {
		return nil, err
	} else {
		this.resolver = resolver
	}

	// success
	return this, nil
}

// Close a logger
func (this *driver) Close() error {
	// Close servers
	for _, server := range this.servers {
		server.Shutdown()
	}
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE METHODS

// Register a service and announce the service when queries occur
func (this *driver) Register(service *gopi.RPCServiceRecord) error {
	if server, err := zeroconf.Register(service.Name, service.Type, this.domain, int(service.Port), service.Text, nil); err != nil {
		return err
	} else {
		this.servers = append(this.servers, server)
		return nil
	}
}

// Browse will find service entries
func (this *driver) Browse(ctx context.Context, serviceType string, callback gopi.RPCBrowseFunc) error {
	entries := make(chan *zeroconf.ServiceEntry)
	if err := this.resolver.Browse(ctx, serviceType, this.domain, entries); err != nil {
		return err
	} else {
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				callback(&gopi.RPCServiceRecord{
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
			callback(nil)
		}(entries)
		return nil
	}
}

func (this *driver) DefaultServiceType(network string) string {
	return "_gopi._" + network
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	return fmt.Sprintf("sys.mdns{ domain=\"%v\" registrations=%v }", this.domain, "TODO")
}
