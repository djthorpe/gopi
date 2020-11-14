package client

import (
	"context"
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
	grpc "google.golang.org/grpc"
	reflection "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type conn struct {
	sync.Mutex

	conn *grpc.ClientConn
	stub reflection.ServerReflectionClient
}

/////////////////////////////////////////////////////////////////////
// INIT

func NewConn(c *grpc.ClientConn) gopi.Conn {
	// Check incoming parameters
	if c == nil {
		return nil
	}

	// Create a connection, get stub for reflection
	this := new(conn)
	if stub := reflection.NewServerReflectionClient(c); stub == nil {
		return nil
	} else {
		this.conn = c
		this.stub = stub
	}

	return this
}

/////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *conn) Addr() string {
	if this.conn != nil {
		return this.conn.Target()
	} else {
		return ""
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *conn) Close() error {
	var result error

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.conn != nil {
		if err := this.conn.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	this.conn = nil

	return result
}

func (this *conn) ListServices(ctx context.Context) ([]string, error) {
	var services []string

	// Exclusive lock
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Create stream
	stream, err := this.stub.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer stream.CloseSend()

	// Send message
	if err := stream.Send(&reflection.ServerReflectionRequest{
		MessageRequest: &reflection.ServerReflectionRequest_ListServices{},
	}); err != nil {
		return nil, err
	}

	// Receive response
	resp, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	// Enumerate services
	for _, service := range resp.GetListServicesResponse().GetService() {
		services = append(services, service.Name)
	}

	// Return success
	return services, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *conn) String() string {
	str := "<conn"
	str += fmt.Sprintf(" addr=%q", this.Addr())
	if this.conn != nil {
		str += " state=" + fmt.Sprint(this.conn.GetState())
	} else {
		str += " state=CLOSED"
	}
	return str + ">"
}
