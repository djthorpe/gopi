package gopi

import (
	"context"
	"net"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type ServiceFlag uint

/////////////////////////////////////////////////////////////////////
// INTERFACES

// Server is a generic RPC server, which can serve responses for
// registered services to clients
type Server interface {
	// Register an RPC service with the server
	RegisterService(interface{}, Service) error

	// Start server in background and return
	StartInBackground(network, addr string) error

	// Stop server, when argument is true forcefully disconnects any clients
	Stop(bool) error

	// Addr returns the address of the server, or empty if not connected
	Addr() string

	// NewStreamContext returns a streaming context which should be used
	// to cancel streaming to clients when the server is shutdown
	NewStreamContext() context.Context
}

// Service defines an RPC service
type Service interface{}

// ConnPool is a factory of client connections
type ConnPool interface {
	Connect(network, addr string) (Conn, error)
}

// Conn is a connection to a remote server
type Conn interface {
	// Addr returns the bound server address, or empty string if connection is closed
	Addr() string

	// Mutex
	Lock()   // Lock during RPC call
	Unlock() // Unlock at end of RPC call

	// ListServices returns a list of all services supported by the
	// remote server
	ListServices(context.Context) ([]string, error)

	// NewStub returns the stub for a named service
	NewStub(string) ServiceStub

	// Err translates service error codes to gopi error types
	Err(error) error
}

// ServiceStub is a client-side stub used to invoke remote service methods
type ServiceStub interface {
	New(Conn)
}

/////////////////////////////////////////////////////////////////////
// SERVICE DISCOVERY

type ServiceDiscovery interface {
	// NewServiceRecord returns a record from service, name, port, txt and
	// flags for IP4, IP6 or both
	NewServiceRecord(string, string, uint16, []string, ServiceFlag) (ServiceRecord, error)

	// EnumerateServices queries for available service names
	EnumerateServices(context.Context) ([]string, error)

	// Lookup queries for records for a service name
	Lookup(context.Context, string) ([]ServiceRecord, error)

	// Serve will respond to service discovery queries and
	// de-register those services when ending
	Serve(context.Context, []ServiceRecord) error
}

type ServiceRecord interface {
	Instance() string
	Service() string
	Name() string
	Host() string
	Port() uint16
	Addrs() []net.IP
	Txt() []string
}

/////////////////////////////////////////////////////////////////////
// SERVICES

type PingService interface {
	Service
}

type PingStub interface {
	ServiceStub

	Ping(ctx context.Context) error
	Version(ctx context.Context) (Version, error)
	ListServices(context.Context) ([]string, error) // Return a list of services supported
}

type InputService interface {
	Service
}

type InputStub interface {
	ServiceStub

	Stream(ctx context.Context, ch chan<- InputEvent) error
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	SERVICE_FLAG_NONE ServiceFlag = 0
	SERVICE_FLAG_IP4  ServiceFlag = (1 << iota)
	SERVICE_FLAG_IP6
	SERVICE_FLAG_MAX = SERVICE_FLAG_IP6
)
