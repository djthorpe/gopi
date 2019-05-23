/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"context"
	"net"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// RPCServiceRecord defines a service which can be registered or discovered
// on the network
type RPCServiceRecord interface {
	Name() string
	Subtype() string
	Service() string
	Port() uint
	Text() []string
	Host() string
	IP4() []net.IP
	IP6() []net.IP
	TTL() time.Duration
}

// RPCEventType is an enumeration of event types
type RPCEventType uint

// RPCFlag is a set of flags modifying behavior of client/service
type RPCFlag uint

// RPCNewClientFunc creates a new client with a network connection
// returns nil otherwise
type RPCNewClientFunc func(RPCClientConn) RPCClient

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// RPCServiceDiscovery is the driver for discovering services on the network using
// mDNS or another mechanism
type RPCServiceDiscovery interface {
	Driver
	Publisher

	// Register a service record on the network, and cache it
	Register(RPCServiceRecord) error

	// Lookup service instances by name
	Lookup(ctx context.Context, service string) ([]RPCServiceRecord, error)

	// Return list of service names
	EnumerateServices(ctx context.Context) ([]string, error)

	// Return all cached service instances for a service name
	ServiceInstances(service string) []RPCServiceRecord
}

// RPCService is a driver which implements all the necessary methods to
// handle remote calls
type RPCService interface {
	Driver

	// CancelRequests is called by the server to gracefully end any
	// on-going streaming requests, but before the service is shutdown
	CancelRequests() error
}

// RPCServer is the server which serves RPCModule methods to
// a remote RPCClient
type RPCServer interface {
	Driver
	Publisher

	// Starts an RPC server in currently running thread.
	// The method will not return until Stop is called
	// which needs to be done in a different thread
	Start() error

	// Stop RPC server. If halt is true then it immediately
	// ends the server without waiting for current requests to
	// be served
	Stop(halt bool) error

	// Return address the server is bound to, or nil if
	// the server is not running
	Addr() net.Addr

	// Return service record, or nil when the service record
	// cannot be generated. The service should be of the format
	// _<service>._tcp and the subtype can only be alphanumeric
	Service(service, subtype, name string, text ...string) RPCServiceRecord
}

// RPCClientPool implements a pool of client connections for communicating
// with an RPC server and aides discovery new service records
type RPCClientPool interface {
	Driver
	Publisher

	// Connect and disconnect
	Connect(service RPCServiceRecord, flags RPCFlag) (RPCClientConn, error)
	ConnectAddr(addr string, flags RPCFlag) (RPCClientConn, error)
	Disconnect(RPCClientConn) error

	// Register clients and create new ones given a service name
	RegisterClient(string, RPCNewClientFunc) error
	NewClient(string, RPCClientConn) RPCClient

	// Lookup service records by parameter - returns records
	// which match either name or addr up to max number of records
	// Can wait for new records and block until cancelled
	Lookup(ctx context.Context, name, addr string, max int) ([]RPCServiceRecord, error)
}

// RPCClientConn implements a single client connection for
// communicating with an RPC server
type RPCClientConn interface {
	Driver

	// Mutex lock for the connection
	Lock()
	Unlock()

	// Properties
	Addr() string
	Connected() bool
	Timeout() time.Duration
	Services() ([]string, error)
}

// RPCClient contains the set of RPC methods. Currently
// anything can be an RPCClient
type RPCClient interface{}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RPC_EVENT_NONE            RPCEventType = iota
	RPC_EVENT_SERVER_STARTED               // RPC Server started
	RPC_EVENT_SERVER_STOPPED               // RPC Server stopped
	RPC_EVENT_SERVICE_ADDED                // Service instance lookup (new)
	RPC_EVENT_SERVICE_UPDATED              // Service instance lookup (updated)
	RPC_EVENT_SERVICE_REMOVED              // Service instance lookup (removed)
	RPC_EVENT_SERVICE_EXPIRED              // Service instance lookup (expired)
	RPC_EVENT_SERVICE_NAME                 // Service name discovered
	RPC_EVENT_CLIENT_CONNECTED
	RPC_EVENT_CLIENT_DISCONNECTED
)

const (
	RPC_FLAG_NONE     RPCFlag = 0
	RPC_FLAG_INET_UDP RPCFlag = (1 << iota) // Use UDP protocol (TCP assumed otherwise)
	RPC_FLAG_INET_V4  RPCFlag = (1 << iota) // Use V4 addressing
	RPC_FLAG_INET_V6  RPCFlag = (1 << iota) // Use V6 addressing
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t RPCEventType) String() string {
	switch t {
	case RPC_EVENT_NONE:
		return "RPC_EVENT_NONE"
	case RPC_EVENT_SERVER_STARTED:
		return "RPC_EVENT_SERVER_STARTED"
	case RPC_EVENT_SERVER_STOPPED:
		return "RPC_EVENT_SERVER_STOPPED"
	case RPC_EVENT_SERVICE_ADDED:
		return "RPC_EVENT_SERVICE_ADDED"
	case RPC_EVENT_SERVICE_UPDATED:
		return "RPC_EVENT_SERVICE_UPDATED"
	case RPC_EVENT_SERVICE_REMOVED:
		return "RPC_EVENT_SERVICE_REMOVED"
	case RPC_EVENT_SERVICE_EXPIRED:
		return "RPC_EVENT_SERVICE_EXPIRED"
	case RPC_EVENT_SERVICE_NAME:
		return "RPC_EVENT_SERVICE_NAME"
	case RPC_EVENT_CLIENT_CONNECTED:
		return "RPC_EVENT_CLIENT_CONNECTED"
	case RPC_EVENT_CLIENT_DISCONNECTED:
		return "RPC_EVENT_CLIENT_DISCONNECTED"
	default:
		return "[?? Invalid RPCEventType value]"
	}
}
