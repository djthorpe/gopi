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

	// Framework
	"github.com/djthorpe/gopi"

	// Protocol Buffer definition
	pb "github.com/djthorpe/gopi/protobuf/helloworld"
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

	// Register service with server
	config.Server.Register(this)

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
