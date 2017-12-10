/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"errors"
	"fmt"
	"net"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	reflection "google.golang.org/grpc/reflection"
)

// Server is the RPC server configuration
type Server struct {
	SSLKey         string
	SSLCertificate string
	Port           uint
}

type server struct {
	log        gopi.Logger
	port       uint
	server     *grpc.Server
	addr       net.Addr
	eventchans []chan gopi.RPCEvent
}

////////////////////////////////////////////////////////////////////////////////
// SERVER OPEN AND CLOSE

// Open a logger
func (config Server) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<grpc.Server>Open(port=%v,sslcert=%v,sslkey=%v)", config.Port, config.SSLCertificate, config.SSLKey)

	this := new(server)
	this.log = log
	this.port = config.Port

	if config.SSLKey != "" || config.SSLCertificate != "" {
		if creds, err := credentials.NewServerTLSFromFile(config.SSLCertificate, config.SSLKey); err != nil {
			return nil, err
		} else {
			this.server = grpc.NewServer(grpc.Creds(creds))
		}
	} else {
		this.server = grpc.NewServer()
	}

	this.addr = nil
	this.eventchans = make([]chan gopi.RPCEvent, 0)

	// Register reflection service on gRPC server.
	reflection.Register(this.server)

	// success
	return this, nil
}

// Close server
func (this *server) Close() error {
	this.log.Debug2("<grpc.Server>Close()")
	return this.Stop(true)
}

////////////////////////////////////////////////////////////////////////////////
// SERVE

func (this *server) Start(module ...gopi.RPCModule) error {
	// Check for serving
	if this.addr != nil {
		return errors.New("Cannot call Start() when server already started")
	} else if lis, err := net.Listen("tcp", portString(this.port)); err != nil {
		return err
	} else {
		this.addr = lis.Addr()
		this.emitEvent(&Event{gopi.RPC_EVENT_SERVER_STARTED})
		this.log.Debug("<grpc.Server>{ addr=%v }", this.addr)
		err := this.server.Serve(lis) // blocking call
		this.emitEvent(&Event{gopi.RPC_EVENT_SERVER_STOPPED})
		this.addr = nil
		return err
	}
}

func (this *server) Stop(halt bool) error {
	// Stop server
	if this.addr != nil {
		if halt {
			this.log.Debug("<grpc.Server>Stop()")
			this.server.Stop()
		} else {
			this.log.Debug("<grpc.Server>GracefulStop()")
			this.server.GracefulStop()
		}
	}

	// Return success
	return nil
}

func (this *server) Addr() net.Addr {
	return this.addr
}

///////////////////////////////////////////////////////////////////////////////
// EVENTS

// Events() creates a new channel on which to emit events and the channel
// is returned. Subsequent events are emitted on all event channels
func (this *server) Events() chan gopi.RPCEvent {
	eventchan := make(chan gopi.RPCEvent)
	this.eventchans = append(this.eventchans, eventchan)
	return eventchan
}

// emitEvent broadcasts events onto listening channels
func (this *server) emitEvent(evt gopi.RPCEvent) {
	for _, c := range this.eventchans {
		if c != nil {
			c <- evt
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// SERVICE

func (this *server) Service(name string) *gopi.RPCService {
	// Can't return a service unless the server is started
	if this.addr == nil {
		return nil
	}
	// Can't return non-TCP
	if this.addr.Network() != "tcp" {
		return nil
	}
	// Can't register if name is blank
	if strings.TrimSpace(name) == "" {
		return nil
	}
	// Return service
	if addr, ok := this.addr.(*net.TCPAddr); ok == false {
		return nil
	} else {
		return &gopi.RPCService{
			Name: strings.TrimSpace(name),
			Type: "_gopi._tcp",
			Port: uint(addr.Port),
		}
	}

}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *server) String() string {
	if this.addr != nil {
		return fmt.Sprintf("<grpc.Server>{ serving,addr=%v }", this.addr)
	} else if this.port == 0 {
		return fmt.Sprintf("<grpc.Server>{ idle }")
	} else {
		return fmt.Sprintf("<grpc.Server>{ idle,port=%v }", this.port)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func portString(port uint) string {
	if port == 0 {
		return ""
	} else {
		return fmt.Sprint(":", port)
	}
}
