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

// Server is a generic gRPC server, which can serve registered services
type Server interface {
	RegisterService(interface{}, Service) error   // Register an RPC service
	StartInBackground(network, addr string) error // Start server in background and return
	Stop(bool) error                              // Stop server, when argument is true forcefully disconnects any clients
	Addr() string                                 // Addr returns the address of the server, or empty if not connected
}

// Service defines an RPC service, which can cancel any on-going streams
// when server stops
type Service interface {
	CancelStreams()
}

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

	// Methods
	ListServices(context.Context) ([]string, error) // Return a list of services supported
	NewStub(string) ServiceStub                     // Return the stub for a named service
}

// ServiceStub is a client-side stub used to invoke remote service methods
type ServiceStub interface {
	New(Conn)
}

/////////////////////////////////////////////////////////////////////
// SERVICE DISCOVERY

type ServiceDiscovery interface {
	// NewServiceRecord returns a record from service, name, port, txt and
	// flags for binding to IP4, IP6 or both
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
	Service() string
	Name() string
	HostPort() []string
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

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	SERVICE_FLAG_NONE ServiceFlag = 0
	SERVICE_FLAG_IP4  ServiceFlag = (1 << iota)
	SERVICE_FLAG_IP6
	SERVICE_FLAG_MAX = SERVICE_FLAG_IP6
)
