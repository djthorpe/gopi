/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package helloworld

import (
	"context"
	"fmt"
	"reflect"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/sys/rpc/grpc"

	// Protocol buffers
	pb "github.com/djthorpe/gopi/rpc/protobuf/helloworld"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server gopi.RPCServer
}

type service struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.Service.helloworld>Open{ server=%v}", config.Server)

	this := new(service)
	this.log = log

	// Register service with GRPC server
	pb.RegisterGreeterServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.Service.helloworld>Close{}")

	// No resources to release

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPCService implementation

func (this *service) GRPCHook() reflect.Value {
	return reflect.ValueOf(pb.RegisterGreeterServer)
}

func (this *service) CancelRequests() error {
	// No need to cancel any requests since none are streaming
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Stringify

func (this *service) String() string {
	return fmt.Sprintf("grpc.service.helloworld{}")
}

////////////////////////////////////////////////////////////////////////////////
// SayHello method

func (this *service) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	if req.Name == "" {
		req.Name = "World"
	}
	return &pb.HelloReply{
		Message: "Hello, " + req.Name,
	}, nil
}
