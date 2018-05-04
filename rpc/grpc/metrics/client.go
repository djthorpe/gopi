/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	"context"
	"fmt"

	// Framework
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi/sys/rpc/grpc"

	// Protocol buffers
	pb "github.com/djthorpe/gopi/rpc/protobuf/metrics"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
	pb.MetricsClient
	conn gopi.RPCClientConn
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &Client{pb.NewMetricsClient(conn.(grpc.GRPCClientConn).GRPCConn()), conn}
}

func (this *Client) NewContext() context.Context {
	if this.conn.Timeout() == 0 {
		return context.Background()
	} else {
		ctx, _ := context.WithTimeout(context.Background(), this.conn.Timeout())
		return ctx
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Client) Conn() gopi.RPCClientConn {
	return this.conn
}

////////////////////////////////////////////////////////////////////////////////
// CALLS

func (this *Client) Ping() error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MetricsClient.Ping(this.NewContext(), &pb.EmptyRequest{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) HostMetrics() (*pb.HostMetricsReply, error) {
	this.conn.Lock()
	defer this.conn.Unlock()

	if reply, err := this.MetricsClient.HostMetrics(this.NewContext(), &pb.EmptyRequest{}); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	return fmt.Sprintf("<grpc.metrics.client>{ conn=%v }", this.conn)
}
