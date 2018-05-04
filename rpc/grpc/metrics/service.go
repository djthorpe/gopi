/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	"fmt"

	// Framework
	"github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi/sys/rpc/grpc"
	context "golang.org/x/net/context"

	// Protocol buffers
	pb "github.com/djthorpe/gopi/rpc/protobuf/metrics"
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
	log.Debug("<grpc.metrics.service>Open{ server=%v}", config.Server)

	this := new(service)
	this.log = log

	// Register service with GRPC server
	pb.RegisterMetricsServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.metrics.service>Close{}")

	// No resources to release

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPCService implementation

func (this *service) CancelRequests() error {
	// No need to cancel any requests since none are streaming
	return nil
}

func (this *service) Ping(ctx context.Context, request *pb.EmptyRequest) (*pb.EmptyReply, error) {
	// Simple ping method to show server is "up"
	return &pb.EmptyReply{}, nil
}

func (this *service) HostMetrics(ctx context.Context, request *pb.EmptyRequest) (*pb.HostMeticsReply, error) {
	// Return host metrics
	return &pb.HostMeticsReply{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Stringify

func (this *service) String() string {
	return fmt.Sprintf("grpc.metrics.service{}")
}
