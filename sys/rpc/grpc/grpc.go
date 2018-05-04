/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package grpc

import (
	"context"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

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

/////////////////////////////////////////////////////////////////////
// UTILITY FUNCTIONS

func IsErrCanceled(err error) bool {
	if err == nil {
		return false
	}
	if err == context.Canceled {
		return true
	}
	return grpc.Code(err) == codes.Canceled
}

func IsErrDeadlineExceeded(err error) bool {
	if err == nil {
		return false
	}
	if err == context.DeadlineExceeded {
		return true
	}
	return grpc.Code(err) == codes.DeadlineExceeded
}
