package gopi

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

/*
	This file contains definitions for communication over networks
	and implementation of services available over networks:

	* HTTP and gRPC Servers
	* Services
	* Service Discovery

	There are also some example gRPC services (Ping, Input, Metrics)
	which can be used "out of the box".
*/

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

	// Addr returns the address of the server, the path for the file socket
	// or empty if not connected
	Addr() string

	// Returns information about the server
	Flags() ServiceFlag

	// Service returns _http._net or _grpc._net or empty if networking
	// is file socket based
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
// GRPC SERVICES

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
	// Serve static files in child folders with root URL as "path"
	Serve(string) error
}

// HttpTemplate loads and serves templates
type HttpTemplate interface {
	// Serve templates with root URL as "path" with files
	// in folder "docroot"
	Serve(path, docroot string) error

	// Register a document renderer
	RegisterRenderer(HttpRenderer) error
}

// HttpLogger logs request and response metrics
type HttpLogger interface {
	// Log all requests as named measurement
	Log(string) error
}

// HttpRenderer returns content to process with template
// for a request
type HttpRenderer interface {
	// IsModifiedSince should return true if content that
	// would be served for this request by the renderer and has
	// been modified since a certain time and rendering should
	// occur for that path. It should return false if this
	// renderer should not serve the request
	IsModifiedSince(docroot string, req *http.Request, t time.Time) bool

	// ServeContent returns the serving contexttemplate name, content object
	// last modified time for caching or zero-time if no
	// caching should occur, and an error. If the error is a
	// HttpError then the error return to the client is sent
	// correctly or else client gets InternalServerError
	// on error
	ServeContent(docroot string, req *http.Request) (HttpRenderContext, error)
}

// HttpRenderContext represents information used to render
// a response through a template. If no template is returned
// then the content is served without a template. The type
// is returned on the Content-Type field of the response.
// The response is cached if the Modified field is not zero.
type HttpRenderContext struct {
	Template string
	Type     string
	Content  interface{}
	Modified time.Time
}

// HttpError provides the correct error code to the client which
// can be returned by the ServeContent method in order to more correctly
// respond to the client
type HttpError interface {
	// Code returns the HTTP status code
	Code() int

	// Error returns the error message to be returned to the client
	Error() string

	// Path returns the redirect path for http.PermanentRedirect status
	Path() string
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	SERVICE_FLAG_NONE   ServiceFlag = 0
	SERVICE_FLAG_IP4    ServiceFlag = (1 << iota) // IP4 Addressing
	SERVICE_FLAG_IP6                              // IP6 Addressing
	SERVICE_FLAG_SOCKET                           // Unix File Socket transport
	SERVICE_FLAG_TLS                              // TLS (SSL) Communication
	SERVICE_FLAG_FCGI                             // FastCGI Communiction
	SERVICE_FLAG_HTTP                             // HTTP Protocol
	SERVICE_FLAG_GRPC                             // gRPC Protocol
	SERVICE_FLAG_MIN    = SERVICE_FLAG_IP4
	SERVICE_FLAG_MAX    = SERVICE_FLAG_GRPC
)

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f ServiceFlag) String() string {
	if f == SERVICE_FLAG_NONE {
		return f.FlagString()
	}
	str := ""
	for v := SERVICE_FLAG_MIN; v <= SERVICE_FLAG_MAX; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.Trim(str, "|")
}

func (f ServiceFlag) FlagString() string {
	switch f {
	case SERVICE_FLAG_NONE:
		return "SERVICE_FLAG_NONE"
	case SERVICE_FLAG_IP4:
		return "SERVICE_FLAG_IP4"
	case SERVICE_FLAG_IP6:
		return "SERVICE_FLAG_IP6"
	case SERVICE_FLAG_SOCKET:
		return "SERVICE_FLAG_SOCKET"
	case SERVICE_FLAG_TLS:
		return "SERVICE_FLAG_TLS"
	case SERVICE_FLAG_FCGI:
		return "SERVICE_FLAG_FCGI"
	case SERVICE_FLAG_HTTP:
		return "SERVICE_FLAG_HTTP"
	case SERVICE_FLAG_GRPC:
		return "SERVICE_FLAG_GRPC"
	default:
		return "[?? Invalid ServiceFlag value]"
	}
}
