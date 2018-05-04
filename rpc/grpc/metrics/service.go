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
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi/sys/rpc/grpc"
	"github.com/golang/protobuf/ptypes"
	context "golang.org/x/net/context"

	// Protocol buffers
	pb "github.com/djthorpe/gopi/rpc/protobuf/metrics"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server  gopi.RPCServer
	Metrics gopi.Metrics
}

type service struct {
	log     gopi.Logger
	metrics gopi.Metrics
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.metrics.service>Open{ server=%v metrics=%v }", config.Server, config.Metrics)

	this := new(service)
	this.log = log
	this.metrics = config.Metrics

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

func (this *service) HostMetrics(ctx context.Context, request *pb.EmptyRequest) (*pb.HostMetricsReply, error) {
	// Obtain hostname
	if hostname, err := os.Hostname(); err != nil {
		return nil, err
	} else {
		// Return host metrics
		reply := &pb.HostMetricsReply{
			Hostname:      hostname,
			HostUptime:    ptypes.DurationProto(this.metrics.UptimeHost()),
			ServiceUptime: ptypes.DurationProto(this.metrics.UptimeApp()),
		}
		return reply, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Stringify

func (this *service) String() string {
	return fmt.Sprintf("grpc.metrics.service{}")
}
