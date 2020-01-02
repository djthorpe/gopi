/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"context"
	"net"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// RPCServiceRecord defines a service which can be registered or discovered
// on the network
type RPCServiceRecord struct {
	Name    string
	Service string
	Subtype string
	Host    string
	Port    uint16
	Addrs   []net.IP
	Txt     []string
}

// RPCFlag is a set of flags modifying behavior of client/service
type RPCFlag uint

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type RPCServiceDiscovery interface {
	// Lookup service instances by name
	Lookup(ctx context.Context, service string) ([]RPCServiceRecord, error)

	// Return list of service names
	EnumerateServices(ctx context.Context) ([]string, error)

	// Return all cached service instances for a service name
	ServiceInstances(service string) []RPCServiceRecord
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
