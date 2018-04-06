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
	"os"
	"reflect"
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	evt "github.com/djthorpe/gopi/util/event"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	reflection "google.golang.org/grpc/reflection"
)

// Server is the RPC server configuration
type Server struct {
	SSLKey         string
	SSLCertificate string
	Port           uint
}

type server struct {
	log    gopi.Logger
	port   uint
	server *grpc.Server
	addr   net.Addr
	pubsub *evt.PubSub
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	grpclog_once sync.Once
	grpclogger   grpclog.LoggerV2
)

////////////////////////////////////////////////////////////////////////////////
// SERVER OPEN AND CLOSE

// Open the server
func (config Server) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.Server>Open(port=%v,sslcert=\"%v\",sslkey=\"%v\")", config.Port, config.SSLCertificate, config.SSLKey)

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

	// Fan out events to subscribers
	this.pubsub = evt.NewPubSub(0)

	// Register reflection service on gRPC server.
	reflection.Register(this.server)

	// Set GRPC logger
	grpclog_once.Do(func() {
		// TODO: Create grpclogger if it doesn't yet exist
		// then route through to 'log'
		//grpclog.SetLoggerV2(grpclogger)
	})

	// success
	return this, nil
}

// Close server
func (this *server) Close() error {
	this.log.Debug("<grpc.Server>Close( addr=%v )", this.addr)

	// Ungracefully stop the server
	err := this.Stop(true)
	if err != nil {
		this.log.Warn("grpc.Server: %v", err)
	}

	// Release resources
	this.pubsub.Close()
	this.pubsub = nil
	this.addr = nil
	this.server = nil

	// Return any error that occurred
	return err
}

////////////////////////////////////////////////////////////////////////////////
// SERVE

func (this *server) Start() error {
	this.log.Debug2("<grpc.Server>Start()")

	// Check for serving
	if this.addr != nil {
		return errors.New("Cannot call Start() when server already started")
	} else if lis, err := net.Listen("tcp", portString(this.port)); err != nil {
		return err
	} else {
		// Start server
		this.addr = lis.Addr()
		this.emit(&Event{source: this, t: gopi.RPC_EVENT_SERVER_STARTED})
		this.log.Debug("<grpc.Server>{ addr=%v }", this.addr)
		err := this.server.Serve(lis) // blocking call
		this.emit(&Event{source: this, t: gopi.RPC_EVENT_SERVER_STOPPED})
		this.addr = nil
		return err
	}
}

func (this *server) Stop(halt bool) error {
	// Stop server
	if this.addr != nil {
		if halt {
			this.log.Debug2("<grpc.Server>Stop()")
			this.server.Stop()
		} else {
			this.log.Debug2("<grpc.Server>GracefulStop()")
			this.server.GracefulStop()
		}
	}

	// Return success
	return nil
}

// Addr returns the currently listening address or will return
// nil if the server is not serving requests
func (this *server) Addr() net.Addr {
	return this.addr
}

func (this *server) Register(service gopi.RPCService) error {
	this.log.Debug2("<grpc.Server>Register( service=%v )", service)
	if this.addr != nil {
		this.log.Error("Register: Unable to register whilst server is started")
		return gopi.ErrOutOfOrder
	}

	// Register service with GRPC
	if callback := service.GRPCHook(); callback.Kind() != reflect.Func {
		return errors.New("Expected callback to be a function")
	} else {
		callback.Call([]reflect.Value{
			reflect.ValueOf(this.server),
			reflect.ValueOf(service),
		})
	}

	// Success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// EVENTS

// Subscribe to events from the server
func (this *server) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

// Unsubscribe from events from the server
func (this *server) Unsubscribe(c <-chan gopi.Event) {
	this.pubsub.Unsubscribe(c)
}

// emit broadcasts events onto listening channels
func (this *server) emit(evt gopi.RPCEvent) {
	this.pubsub.Emit(evt)
}

///////////////////////////////////////////////////////////////////////////////
// SERVICE

func (this *server) Service(service string, text ...string) *gopi.RPCServiceRecord {
	if hostname, err := os.Hostname(); err != nil {
		this.log.Error("<grpc.Server>Service: %v", err)
		return nil
	} else {
		return this.ServiceWithName(service, hostname, text...)
	}
}

func (this *server) ServiceWithName(service, name string, text ...string) *gopi.RPCServiceRecord {
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
		return &gopi.RPCServiceRecord{
			Name: strings.TrimSpace(name),
			Type: serviceType(service, addr.Network()),
			Port: uint(addr.Port),
			Text: text,
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

func serviceType(service, network string) string {
	return "_" + service + "._" + network
}
