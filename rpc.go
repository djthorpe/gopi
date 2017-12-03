/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"

import (
	"context"
	"net"
)

// RPCService defines a service which can be registered
// or discovered on the network
type RPCService struct {
	Name string
	Type string
	Port uint
	Text []string
	Host string
	IP4  net.IP
	IP6  net.IP
}

// RPCBrowseFunc is the callback function for when a service is discovered
// on the network. It's called with a nil parameter when no more services
// are found
type RPCBrowseFunc func(service *RPCService)

// RPCDiscoveryDriver is the driver for discovering
// services on the network using mDNS or another mechanism
type RPCDiscoveryDriver interface {
	// Register a service record on the network
	Register(service *RPCService) error

	// Browse for service records on the network with context
	Browse(ctx context.Context, serviceType string, callback RPCBrowseFunc) error

	// Return a list of current services
	Services() []*RPCService
}
