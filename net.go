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

// Server is a generic RPC or HTTP server, which can serve responses for
// registered services to clients
type Server interface {
	// Register an RPC or HTTP service with the server
	RegisterService(interface{}, Service) error

	// Start server in background and return
	StartInBackground(network, addr string) error

	// Stop server, when argument is true forcefully disconnects any clients
	Stop(bool) error

	// Addr returns the address of the server, or empty if not connected
	Addr() string

	// SSL returns true if SSL is enabled
	SSL() bool

	// Service returns _http._net or _grpc._net
	Service() string

	// NewStreamContext returns a streaming context which should be used
	// to cancel streaming to clients when the server is shutdown
	NewStreamContext() context.Context
}

// Service defines an RPC or HTTP service. At the moment HTTP services must
// adhere to the http.Handler interface.
type Service interface{}

// ConnPool is a factory of client connections
type ConnPool interface {
	// Connect accepts a network (tcp, udp, unix) and either
	// an IP:port or a path name to a socket and returns a connection
	Connect(network, addr string) (Conn, error)

	// ConnectService accepts a network (tcp, udp, unix) and a service
	// name. If network is 'unix' or service is an IP:port or host:port
	// it will default to a normal Connect. The service name should be
	// alphanumeric and the flags will determine if a connection by hostname,
	// IP4 or IP6 connection is made. In addition, the service parameter can
	// be either <service>:<name> or <name> to connecto to the correct service
	// instance.
	ConnectService(ctx context.Context, network, service string, flags ServiceFlag) (Conn, error)
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
	Zone() string
	Host() string
	Port() uint16
	Addrs() []net.IP
	Txt() []string
}

/////////////////////////////////////////////////////////////////////
// gRPC SERVICES

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

type MetricsService interface {
	Service
}

type MetricsStub interface {
	ServiceStub

	// List returns the array of defined measurements
	List(context.Context) ([]Measurement, error)

	// Stream emits measurements defined by name filter on
	// the provided channel until context is cancelled. Where
	// the name filter is empty, all measurements are emitted
	Stream(context.Context, string, chan<- Measurement) error
}

/////////////////////////////////////////////////////////////////////
// HTTP SERVICES

// HttpStatic serves files and folders from the filesystem
type HttpStatic interface {
	// Serve a folder and child folders with root URL as "path"
	ServeFolder(path, folder string) error
}

// HttpLogger logs request and response metrics
type HttpLogger interface {
	// Log all requests as named measurement
	Log(string) error
}

// HttpTemplate loads and serves templates
type HttpTemplate interface {
	// Serve a template for a path
	ServeTemplate(path, template string) error
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	SERVICE_FLAG_NONE ServiceFlag = 0
	SERVICE_FLAG_IP4  ServiceFlag = (1 << iota)
	SERVICE_FLAG_IP6
	SERVICE_FLAG_MAX = SERVICE_FLAG_IP6
)
