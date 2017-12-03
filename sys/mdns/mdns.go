/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mdsn /* import "github.com/djthorpe/gopi/sys/mdns" */

import (
	"context"
	"fmt"
	"log"

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
	MDNS_DOMAIN = "local."
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register logger
	gopi.RegisterModule(gopi.Module{
		Name: "sys/mdns",
		Type: gopi.MODULE_TYPE_MDNS,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("mdns.domain", MDNS_DOMAIN, "Domain")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			domain, _ := app.AppFlags.GetString("mdns.domain")
			return gopi.Open(Config{Domain: domain}, app.Logger)
		},
	})
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open a logger
func (config Config) Open(log gopi.Logger) (gopi.Driver, error) {

	this := new(driver)
	this.log = log
	if config.Domain == "" {
		this.domain = MDNS_DOMAIN
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
func (this *driver) Register(serviceName, serviceType string, port uint, txtRecords []string) error {
	if server, err := zeroconf.Register(serviceName, serviceType, this.domain, int(port), txtRecords, nil); err != nil {
		return err
	} else {
		this.servers = append(this.servers, server)
		return nil
	}
}

// Browse will find service entries
func (this *driver) Browse(ctx context.Context, serviceType string, callback gopi.RPCServiceDiscoveryFunc) error {
	entries := make(chan *zeroconf.ServiceEntry)
	if err := this.resolver.Browse(ctx, serviceType, this.domain, entries); err != nil {
		return err
	} else {
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				callback()
				log.Println(entry)
			}
			log.Println("No more entries.")
		}(entries)
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	return fmt.Sprintf("sys.mdns{ domain=\"%v\" registrations=%v }", this.domain, "TODO")
}
