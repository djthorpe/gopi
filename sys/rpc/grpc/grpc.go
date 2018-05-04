/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package grpc

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "google.golang.org/grpc"
)

// GRPCServer interface is an RPCServer which also
// returns gRPC-specific properties
type GRPCServer interface {
	gopi.RPCServer

	// Return the gRPC Server object
	GRPCServer() *grpc.Server
}

// GRPCClientConn is an RPCClientConn which also
// returns gRPC-specific properties
type GRPCClientConn interface {
	gopi.RPCClientConn

	// Return the gRPC ClientConn object
	GRPCConn() *grpc.ClientConn
}
