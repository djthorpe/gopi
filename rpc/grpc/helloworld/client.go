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

	// Framework
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi/sys/rpc/grpc"

	// Protocol buffers
	pb "github.com/djthorpe/gopi/rpc/protobuf/helloworld"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
	pb.GreeterClient
	conn gopi.RPCClientConn
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewGreeterClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &Client{pb.NewGreeterClient(conn.(grpc.GRPCClientConn).GRPCConn()), conn}
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

func (this *Client) SayHello(name string) (string, error) {
	this.conn.Lock()
	defer this.conn.Unlock()

	// Perform SayHello
	if reply, err := this.GreeterClient.SayHello(this.NewContext(), &pb.HelloRequest{Name: name}); err != nil {
		return "", err
	} else {
		return reply.Message, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	return fmt.Sprintf("<helloworld.MyGreeterClient>{ conn=%v }", this.conn)
}
