/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	// Frameworks
	"fmt"
	"net"

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

func (this *server) Start() error {
	if this.serving {
		return grpc.ErrServerStopped
	} else if lis, err := net.Listen("tcp", portString(this.port)); err != nil {
		return err
	} else if err := this.server.Serve(lis); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *server) Stop(halt bool) error {
	if this.serving {
		if halt {
			this.server.Stop()
		} else {
			this.server.GracefulStop()
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *server) String() string {
	return fmt.Sprintf("<gopi.rpcserver.grpc>{ port=%v }", this.port)
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
