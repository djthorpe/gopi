package client

import (
	"context"
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
	multierror "github.com/hashicorp/go-multierror"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	reflection "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type conn struct {
	sync.Mutex
	*grpc.ClientConn

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
		this.ClientConn = c
		this.stub = stub
	}

	return this
}

/////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *conn) Addr() string {
	if this.ClientConn != nil {
		return this.ClientConn.Target()
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

	if this.ClientConn != nil {
		if err := this.ClientConn.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	this.ClientConn = nil

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

func (this *conn) NewStub(service string) gopi.ServiceStub {

	// Exclusive lock
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Lookup stub and create a new one
	if this.ClientConn == nil {
		return nil
	} else if stub := graph.NewServiceStub(service); stub == nil {
		return nil
	} else {
		stub.New(this)
		return stub
	}
}

func (this *conn) Err(err error) error {
	switch grpc.Code(err) {
	case codes.Canceled:
		return context.Canceled
	case codes.Unavailable:
		return gopi.ErrUnexpectedResponse
	case codes.DeadlineExceeded:
		return context.DeadlineExceeded
	default:
		return err
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *conn) String() string {
	str := "<conn"
	str += fmt.Sprintf(" addr=%q", this.Addr())
	if this.ClientConn != nil {
		str += " state=" + fmt.Sprint(this.ClientConn.GetState())
	} else {
		str += " state=CLOSED"
	}
	return str + ">"
}
