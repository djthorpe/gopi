package gopi

import "context"

/////////////////////////////////////////////////////////////////////
// INTERFACES

// Server is a generic gRPC server, which can serve registered services
type Server interface {
	RegisterService(interface{}, Service) error   // Register an RPC service
	StartInBackground(network, addr string) error // Start server in background and return
	Stop(bool) error                              // Stop server, when argument is true forcefully disconnects any clients
	Addr() string                                 // Addr returns the address of the server, or empty if not connected
}

// Service defines an RPC service, which can cancel any on-going requests
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
}

/////////////////////////////////////////////////////////////////////
// SERVICES

type PingService interface {
	Service
}
