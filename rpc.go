/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// RPCServiceRecord defines a service which can be registered or discovered
// on the network
type RPCServiceRecord struct {
	Name    string
	Service string
	Host    string
	Port    uint16
	Addrs   []net.IP
	Txt     []string
}

type (
	RPCFlag      uint // RPCFlag is a set of flags modifying behavior
	RPCEventType uint // RPCEventType is an enumeration of event types
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// RPCServiceDiscovery will lookup services and classes of service
type RPCServiceDiscovery interface {
	// Lookup service instances by name
	Lookup(ctx context.Context, service string) ([]RPCServiceRecord, error)

	// Return list of service names
	EnumerateServices(ctx context.Context) ([]string, error)

	// Implements gopi.Unit
	Unit
}

// RPCServiceRegister will register services
type RPCServiceRegister interface {
	// Register service record, and de-register when deadline is exceeded
	Register(ctx context.Context, record RPCServiceRecord) error

	// Implements gopi.Unit
	Unit
}

// RPCServer serves requests for one or more services
type RPCServer interface {
	// Start an RPC server in currently running thread.
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

	// Implements gopi.Unit
	Unit
}

// RPCService is a driver which implements all the necessary
// methods to handle remote calls
type RPCService interface {
	// CancelRequests is called by the server to gracefully end any
	// on-going streaming requests, but before the service is shutdown
	CancelRequests() error

	// Implements gopi.Unit
	Unit
}

// RPCClientPool represents a pool of clients and connections to
// remove services
type RPCClientPool interface {
	// Lookup service records by parameter - returns records
	// which match either a service name up to max number of
	// records (or zero for unlimited). Will wait for new records
	// and block until cancelled
	Lookup(ctx context.Context, addr string, max uint) ([]RPCServiceRecord, error)

	// Connect to service by service record
	Connect(RPCServiceRecord, RPCFlag) (RPCClientConn, error)

	// Connect to service by address and port
	ConnectAddr(net.IP, uint16) (RPCClientConn, error)

	// Connect to service by unix socket
	ConnectFifo(string) (RPCClientConn, error)

	// Disconnect from service
	Disconnect(RPCClientConn) error

	// Implements gopi.Unit
	Unit
}

// RPCClientConn implements a single client connection for
// communicating with an RPC server
type RPCClientConn interface {
	// Return connection information, or nil if not connected
	Addr() string

	// Return service names supported by connection
	Services() ([]string, error)

	// Mutex locking to ensure one request at a time
	Lock()
	Unlock()

	// Implements gopi.Unit
	Unit
}

type RPCClientStub interface {
	// Return the connection associated with the client
	Conn() RPCClientConn

	// Implements gopi.Unit
	Unit
}

// RPCEvent is emitted on service discovery and server events
type RPCEvent interface {
	// Type of event
	Type() RPCEventType

	// Service record associated with event
	Service() RPCServiceRecord

	// Time-to-live value for event
	TTL() time.Duration

	// Implements gopi.Event
	Event
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RPC_FLAG_NONE          RPCFlag = 0
	RPC_FLAG_INET_UDP      RPCFlag = (1 << iota) >> 1 // Use UDP protocol (TCP assumed otherwise)
	RPC_FLAG_INET_V4                                  // Use V4 addressing
	RPC_FLAG_INET_V6                                  // Use V6 addressing
	RPC_FLAG_SERVICE_FIRST                            // Use first service
	RPC_FLAG_SERVICE_ANY                              // Use any service
	RPC_FLAG_MIN           = RPC_FLAG_INET_UDP
	RPC_FLAG_MAX           = RPC_FLAG_SERVICE_ANY
)

const (
	RPC_EVENT_NONE            RPCEventType = iota
	RPC_EVENT_SERVER_STARTED               // RPC Server started
	RPC_EVENT_SERVER_STOPPED               // RPC Server stopped
	RPC_EVENT_SERVICE_ADDED                // Service instance lookup (new)
	RPC_EVENT_SERVICE_UPDATED              // Service instance lookup (updated)
	RPC_EVENT_SERVICE_REMOVED              // Service instance lookup (removed)
	RPC_EVENT_SERVICE_EXPIRED              // Service instance lookup (expired)
	RPC_EVENT_SERVICE_NAME                 // Service name discovered
	RPC_EVENT_SERVICE_RECORD               // Service record lookup
	RPC_EVENT_CLIENT_CONNECTED
	RPC_EVENT_CLIENT_DISCONNECTED
	RPC_EVENT_MAX = RPC_EVENT_CLIENT_DISCONNECTED
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f RPCFlag) String() string {
	str := ""
	if f == 0 {
		return f.FlagString()
	}
	for v := RPC_FLAG_MIN; v <= RPC_FLAG_MAX; v <<= 1 {
		if f&v == v {
			str += "|" + v.FlagString()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (v RPCFlag) FlagString() string {
	switch v {
	case RPC_FLAG_NONE:
		return "RPC_FLAG_NONE"
	case RPC_FLAG_INET_UDP:
		return "RPC_FLAG_INET_UDP"
	case RPC_FLAG_INET_V4:
		return "RPC_FLAG_INET_V4"
	case RPC_FLAG_INET_V6:
		return "RPC_FLAG_INET_V6"
	case RPC_FLAG_SERVICE_FIRST:
		return "RPC_FLAG_SERVICE_FIRST"
	case RPC_FLAG_SERVICE_ANY:
		return "RPC_FLAG_SERVICE_ANY"
	default:
		return "[?? Invalid PlatformType value]"
	}
}

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
	case RPC_EVENT_SERVICE_RECORD:
		return "RPC_EVENT_SERVICE_RECORD"
	case RPC_EVENT_CLIENT_CONNECTED:
		return "RPC_EVENT_CLIENT_CONNECTED"
	case RPC_EVENT_CLIENT_DISCONNECTED:
		return "RPC_EVENT_CLIENT_DISCONNECTED"
	default:
		return "[?? Invalid RPCEventType value]"
	}
}

func (this RPCServiceRecord) String() string {
	str := "<RPCServiceRecord name=" + strconv.Quote(this.Name)
	if this.Service != "" {
		str += " service=" + strconv.Quote(this.Service)
	}
	if this.Host != "" {
		str += " host=" + strconv.Quote(this.Host)
	}
	if this.Port != 0 {
		str += " port=" + fmt.Sprint(this.Port)
	}
	if len(this.Addrs) > 0 {
		str += " addrs="
		for _, addr := range this.Addrs {
			str += addr.String() + ","
		}
		str = strings.TrimSuffix(str, ",")
	}
	if len(this.Txt) > 0 {
		str += " txt=" + fmt.Sprint(this.Txt)
	}
	str += ">"
	return str
}
