/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"fmt"
	"net"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	grpc "google.golang.org/grpc"
	reflection "google.golang.org/grpc/reflection"
)

// Server is the RPC server configuration
type Server struct {
	Port uint
}

type server struct {
	log     gopi.Logger
	port    uint
	server  *grpc.Server
	addr    net.Addr
	serving bool
}

////////////////////////////////////////////////////////////////////////////////
// SERVER OPEN AND CLOSE

// Open a logger
func (config Server) Open(log gopi.Logger) (gopi.Driver, error) {

	this := new(server)
	this.log = log
	this.port = config.Port
	this.server = grpc.NewServer()
	this.serving = false
	this.addr = nil

	// Register reflection service on gRPC server.
	reflection.Register(this.server)

	// success
	return this, nil
}

// Close a logger
func (this *server) Close() error {
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SERVE

func (this *server) Start(module ...gopi.RPCModule) error {
	// Check for serving
	if this.serving {
		return grpc.ErrServerStopped
	} else if lis, err := net.Listen("tcp", portString(this.port)); err != nil {
		return err
	} else {
		this.addr = lis.Addr()
		this.serving = true
		err := this.server.Serve(lis)
		this.serving = false
		this.addr = nil
		return err
	}
}

func (this *server) StartInBackground(module ...gopi.RPCModule) error {
	// Check for serving
	if this.serving {
		return grpc.ErrServerStopped
	} else if lis, err := net.Listen("tcp", portString(this.port)); err != nil {
		return err
	} else {
		this.addr = lis.Addr()
		this.serving = true
		go func() {
			if err := this.server.Serve(lis); err != nil {
				this.log.Error("<gopi.rpcserver.grpc> Error: %v", err)
			}
			this.serving = false
			this.addr = nil
		}()
	}
	return nil
}

func (this *server) Stop(halt bool) error {
	// Stop server
	if this.serving {
		if halt {
			this.server.Stop()
		} else {
			this.server.GracefulStop()
		}
	}

	// Return success
	return nil
}

func (this *server) Addr() net.Addr {
	return this.addr
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *server) String() string {
	if this.serving {
		return fmt.Sprintf("<gopi.rpcserver.grpc>{ serving,addr=%v }", this.addr)
	} else {
		return fmt.Sprintf("<gopi.rpcserver.grpc>{ idle,port=%v }", this.port)
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
