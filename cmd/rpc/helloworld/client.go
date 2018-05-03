/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package helloworld

import (
	"fmt"

	// Framework
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi/sys/rpc/grpc"

	// Protocol buffers
	pb "github.com/djthorpe/gopi/protobuf/helloworld"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GreeterClient struct {
	pb.GreeterClient
	conn gopi.RPCClientConn
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewGreeterClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &GreeterClient{pb.NewGreeterClient(conn.(grpc.GRPCClientConn).Conn()), conn}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GreeterClient) String() string {
	return fmt.Sprintf("<helloworld.GreeterClient>{ conn=%v }", this.conn)
}
