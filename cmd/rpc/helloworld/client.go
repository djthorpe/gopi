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
	pb "github.com/djthorpe/gopi/protobuf/helloworld"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MyGreeterClient struct {
	pb.GreeterClient
	conn gopi.RPCClientConn
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewGreeterClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &MyGreeterClient{pb.NewGreeterClient(conn.(grpc.GRPCClientConn).Conn()), conn}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *MyGreeterClient) Conn() gopi.RPCClientConn {
	return this.conn
}

////////////////////////////////////////////////////////////////////////////////
// CALLS

func (this *MyGreeterClient) SayHello(name string) (string, error) {
	this.conn.Lock()
	defer this.conn.Unlock()

	// TODO Need to add a deadline
	if reply, err := this.GreeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: name}); err != nil {
		return "", err
	} else {
		return reply.Message, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *MyGreeterClient) String() string {
	return fmt.Sprintf("<helloworld.MyGreeterClient>{ conn=%v }", this.conn)
}
