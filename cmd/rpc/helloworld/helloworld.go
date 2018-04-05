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

	// Framework
	"github.com/djthorpe/gopi"

	// Protobuf
	pb "github.com/djthorpe/gopi/protobuf/helloworld"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register rpc/service:helloworld
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/service:helloworld",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server"},
	})
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	// No configuration parameters
}

type service struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.Service.helloworld>Open")

	this := new(service)
	this.log = log

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.Service.helloworld>Close")

	// No resources to release

	// Success
	return nil
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

/*

func (this *HelloworldService) Register(server gopi.RPCServer) error {
	// Check to make sure we satisfy the interface
	var _ hw.GreeterServer = (*HelloworldService)(nil)
	return server.Fudge(reflect.ValueOf(hw.RegisterGreeterServer), this)
}
*/
