package gopi

import (
	"context"
	"net"
)

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
	EnumerateServices(context.Context) ([]string, error)
	Lookup(context.Context, string) ([]ServiceRecord, error)
	Serve(context.Context, []ServiceRecord) error
}

type ServiceRecord interface {
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
