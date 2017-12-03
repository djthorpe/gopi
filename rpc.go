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
	"fmt"
	"net"
	"strings"
)

// RPCService defines a service which can be registered
// or discovered on the network
type RPCService struct {
	Name string
	Type string
	Port uint
	Text []string
	Host string
	IP4  []net.IP
	IP6  []net.IP
}

// RPCBrowseFunc is the callback function for when a service is discovered
// on the network. It's called with a nil parameter when no more services
// are found
type RPCBrowseFunc func(service *RPCService)

// RPCDiscoveryDriver is the driver for discovering
// services on the network using mDNS or another mechanism
type RPCServiceDiscovery interface {
	// Register a service record on the network
	Register(service *RPCService) error

	// Browse for service records on the network with context
	Browse(ctx context.Context, serviceType string, callback RPCBrowseFunc) error

	// Return a list of current services
	//Services() []*RPCService
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s *RPCService) String() string {
	p := make([]string, 0, 5)
	if s.Name != "" {
		p = append(p, fmt.Sprintf("name=\"%v\"", s.Name))
	}
	if s.Type != "" {
		p = append(p, fmt.Sprintf("type=%v", s.Type))
	}
	if s.Port > 0 {
		p = append(p, fmt.Sprintf("port=%v", s.Port))
	}
	if s.Host != "" {
		p = append(p, fmt.Sprintf("host=%v", s.Host))
	}
	if len(s.IP4) > 0 {
		p = append(p, fmt.Sprintf("ip4=%v", s.IP4))
	}
	if len(s.IP6) > 0 {
		p = append(p, fmt.Sprintf("ip6=%v", s.IP6))
	}
	if len(s.Text) > 0 {
		p = append(p, fmt.Sprintf("txt=%v", s.Text))
	}
	return fmt.Sprintf("<gopi.RPCService>{ %v }", strings.Join(p, " "))
}
