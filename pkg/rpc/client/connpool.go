package client

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type connpool struct {
	gopi.Unit
	sync.Mutex
	gopi.Logger
	gopi.ServiceDiscovery

	conns []gopi.Conn
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	reServiceName = regexp.MustCompile("^_(\\w+)\\._(tcp|udp)\\.$")
	reServiceAddr = regexp.MustCompile("^(\\w+):([a-zA-Z]+\\S*)$")
)

/////////////////////////////////////////////////////////////////////
// INIT

func (this *connpool) New(gopi.Config) error {
	if this.ServiceDiscovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("ServiceDiscovery")
	}

	// Return success
	return nil
}

func (this *connpool) Dispose() error {
	var result error

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close all clients
	for _, c := range this.conns {
		if c != nil {
			if err := c.(*conn).Close(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Return success
	return result
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *connpool) Connect(network, addr string) (gopi.Conn, error) {
	switch network {
	case "tcp":
		this.Debugf("Connect: %q,%q", network, addr)
		if conn, err := grpc.Dial(addr, grpc.WithInsecure()); err != nil {
			return nil, err
		} else if client := NewConn(conn); client == nil {
			return nil, gopi.ErrInternalAppError.WithPrefix(addr)
		} else {
			this.Mutex.Lock()
			defer this.Mutex.Unlock()
			this.conns = append(this.conns, client)
			return client, nil
		}
	default:
		return nil, gopi.ErrNotImplemented.WithPrefix(network)
	}
}

func (this *connpool) ConnectService(ctx context.Context, network, service string, flags gopi.ServiceFlag) (gopi.Conn, error) {
	// Default to Connect if the network is unix
	if network == "unix" {
		return this.Connect(network, service)
	} else if network != "tcp" && network != "udp" {
		return nil, gopi.ErrBadParameter.WithPrefix(network)
	}

	// Name to filter for
	name := ""

	// Default to Connect if service is in a host:port format
	if parts := reServiceAddr.FindStringSubmatch(service); len(parts) == 3 {
		service = parts[1]
		name = parts[2]
		this.Debugf("ConnectService service=%q name=%q", service, name)
	} else if host, port, err := net.SplitHostPort(service); err == nil {
		this.Debugf("ConnectService host=%q port=%q", host, port)
		return this.Connect(network, service)
	}

	// Normalize service name, lookup and connect
	if service, err := fqn(service, network); err != nil {
		return nil, err
	} else if records, err := this.ServiceDiscovery.Lookup(ctx, service); err != nil {
		return nil, err
	} else if addr, err := addr(records, name, flags); err != nil {
		return nil, err
	} else {
		return this.Connect(network, addr)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *connpool) String() string {
	str := "<connpool"
	str += " conns=" + fmt.Sprint(this.conns)
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func fqn(service, network string) (string, error) {
	service = "_" + strings.Trim(service, "_") + "._" + network + "."
	if reServiceName.MatchString(service) == false {
		return "", gopi.ErrBadParameter.WithPrefix(service)
	} else {
		return service, nil
	}
}

func addr(r []gopi.ServiceRecord, name string, flags gopi.ServiceFlag) (string, error) {
	for _, record := range r {
		// Filter by name
		if name != "" && name != record.Name() {
			continue
		}
		// If flags is none, then return hostname
		if flags == gopi.SERVICE_FLAG_NONE {
			return fmt.Sprint(record.Host(), ":", record.Port()), nil
		}
		// Get an address
		for _, addr := range record.Addrs() {
			switch {
			case (flags&gopi.SERVICE_FLAG_IP6 != 0 || flags == gopi.SERVICE_FLAG_NONE) && addr.To4() == nil:
				return fmt.Sprintf("%v:%v", addr.To16(), record.Port()), nil
			case (flags&gopi.SERVICE_FLAG_IP4 != 0 || flags == gopi.SERVICE_FLAG_NONE) && addr.To4() != nil:
				return fmt.Sprintf("%v:%v", addr.To16(), record.Port()), nil
			}
		}
	}
	// No address found
	return "", gopi.ErrNotFound.WithPrefix("ConnectService")
}
