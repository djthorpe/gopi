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
	"reflect"
	"strings"
	"time"
)

// RPCService defines a service which can be registered
// or discovered on the network
type RPCService struct {
	Name string
	Type string
	Port uint
	Text []string
	Host string
	TTL  time.Duration
	IP4  []net.IP
	IP6  []net.IP
}

// RPCBrowseFunc is the callback function for when a service is discovered
// on the network. It's called with a nil parameter when no more services
// are found
type RPCBrowseFunc func(service *RPCService)

// RPCModule is a set of functions which will service RPC calls remotely
type RPCModule interface {
	// Register the module with server before the server starts
	Register(server RPCServer) error

	// Return the service type string
	ServiceType() string
}

// RPCDiscoveryDriver is the driver for discovering
// services on the network using mDNS or another mechanism
type RPCServiceDiscovery interface {
	Driver
	Publisher

	// Register a service record on the network
	Register(service *RPCService) error

	// Browse for service records on the network with context
	Browse(ctx context.Context, serviceType string) error
}

// RPCServer is the server which serves RPCModule methods to
// a remote RPCClient
type RPCServer interface {
	Driver

	// Starts an RPC server in currently running thread, with
	// current set of modules as RPC calls. The method will not
	// return until Stop is called
	Start(module ...RPCModule) error

	// Stop RPC server. If halt is true then it immediately
	// ends the server without waiting for current requests to
	// be served
	Stop(halt bool) error

	// Return address the server is bound to, or nil if
	// the server is not running
	Addr() net.Addr

	// Return service record, or nil when the service record
	// cannot be generated
	Service(name, service string) *RPCService

	// Fudge is something I will fix later. It implements
	// a hook for calling the grpc register functions
	Fudge(callback reflect.Value, module RPCModule) error

	// Return channel on which server events are emitted
	Events() chan RPCEvent
}

// RPCEventType is an enumeration of event types
type RPCEventType uint

// RPCEvent is an event which is emitted by either discovery or
// server.
type RPCEvent interface {
	Type() RPCEventType
}

// RPCClient implements a client for communicating with an RPC server
type RPCClient interface {
	Driver

	// Connect to the remote server
	Connect() error

	// Disconnect from the remote server
	Disconnect() error

	// Return names of modules which server has registered
	Modules() ([]string, error)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RPC_EVENT_NONE RPCEventType = iota
	RPC_EVENT_SERVER_STARTED
	RPC_EVENT_SERVER_STOPPED
)

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
	if s.TTL > 0 {
		p = append(p, fmt.Sprintf("ttl=%v", s.TTL))
	}
	if len(s.Text) > 0 {
		p = append(p, fmt.Sprintf("txt=%v", s.Text))
	}
	return fmt.Sprintf("<gopi.RPCService>{ %v }", strings.Join(p, " "))
}

func (t RPCEventType) String() string {
	switch t {
	case RPC_EVENT_NONE:
		return "RPC_EVENT_NONE"
	case RPC_EVENT_SERVER_STARTED:
		return "RPC_EVENT_SERVER_STARTED"
	case RPC_EVENT_SERVER_STOPPED:
		return "RPC_EVENT_SERVER_STOPPED"
	default:
		return "[?? Invalid RPCEventType value]"
	}
}
