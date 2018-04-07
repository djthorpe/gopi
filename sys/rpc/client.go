/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"reflect"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	reflection_pb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// ClientConn is the RPC client connection
type ClientConn struct {
	Addr       string
	SSL        bool
	SkipVerify bool
	Timeout    time.Duration // Connection timeout
}

type clientconn struct {
	log        gopi.Logger
	addr       string
	ssl        bool
	skipverify bool
	timeout    time.Duration
	conn       *grpc.ClientConn
}

////////////////////////////////////////////////////////////////////////////////
// CLIENT OPEN AND CLOSE

// Open a client
func (config ClientConn) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.ClientConn>Open(addr=%v,ssl=%v,skipverify=%v,timeout=%v)", config.Addr, config.SSL, config.SkipVerify, config.Timeout)

	// Create a client object
	this := new(clientconn)
	this.addr = config.Addr
	this.ssl = config.SSL
	this.skipverify = config.SkipVerify
	this.timeout = config.Timeout
	this.log = log
	this.conn = nil

	// success
	return this, nil
}

// Close client
func (this *clientconn) Close() error {
	this.log.Debug("<grpc.ClientConn>Close{ addr=%v }", this.addr)
	return this.Disconnect()
}

////////////////////////////////////////////////////////////////////////////////
// RPCClientConn interface implementation

func (this *clientconn) Connect() ([]string, error) {
	this.log.Debug2("<grpc.ClientConn>Connect{ addr=%v }", this.addr)
	if this.conn != nil {
		return nil, errors.New("Cannot call Connect() when connection already made")
	}
	opts := make([]grpc.DialOption, 0, 1)

	// SSL options
	if this.ssl {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: this.skipverify})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// Connection timeout options
	if this.timeout > 0 {
		opts = append(opts, grpc.WithTimeout(this.timeout))
	}

	// Dial connection
	if conn, err := grpc.Dial(this.addr, opts...); err != nil {
		return nil, err
	} else {
		this.conn = conn
	}

	// Get services
	reflection := this.newServerReflectionClient()
	if reflection == nil {
		this.log.Warn("grpc.ClientConn: Unable to create reflection client")
		return nil, nil
	}
	defer reflection.CloseSend()

	if services, err := this.listServices(reflection); err != nil {
		this.log.Warn("grpc.ClientConn: %v", err)
		return nil, nil
	} else {
		return services, nil
	}
}

func (this *clientconn) Disconnect() error {
	this.log.Debug2("<grpc.ClientConn>Disconnect{ addr=%v }", this.addr)

	if this.conn != nil {
		err := this.conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

func (this *clientconn) NewService(constructor reflect.Value) (interface{}, error) {
	this.log.Debug2("<grpc.ClientConn>NewService{ func=%v }", constructor)

	if constructor.Kind() != reflect.Func {
		return nil, gopi.ErrBadParameter
	}

	if service := constructor.Call([]reflect.Value{reflect.ValueOf(this.conn)}); len(service) != 1 {
		return nil, gopi.ErrBadParameter
	} else {
		return service[0].Interface(), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *clientconn) String() string {
	if this.conn != nil {
		return fmt.Sprintf("<grpc.ClientConn>{ connected=true addr=%v ssl=%v skipverify=%v }", this.addr, this.ssl, this.skipverify)
	} else {
		return fmt.Sprintf("<grpc.ClientConn>{ connected=false addr=%v ssl=%v skipverify=%v }", this.addr, this.ssl, this.skipverify)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *clientconn) newServerReflectionClient() reflection_pb.ServerReflection_ServerReflectionInfoClient {
	if this.conn == nil {
		return nil
	}
	ctx := context.Background()
	if this.timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, this.timeout)
	}
	if client, err := reflection_pb.NewServerReflectionClient(this.conn).ServerReflectionInfo(ctx); err != nil {
		this.log.Error("Error: %v", err)
		return nil
	} else {
		return client
	}
}

func (this *clientconn) listServices(c reflection_pb.ServerReflection_ServerReflectionInfoClient) ([]string, error) {
	if err := c.Send(&reflection_pb.ServerReflectionRequest{
		MessageRequest: &reflection_pb.ServerReflectionRequest_ListServices{},
	}); err != nil {
		return nil, err
	}

	if resp, err := c.Recv(); err != nil {
		return nil, err
	} else if modules := resp.GetListServicesResponse(); modules == nil {
		return nil, fmt.Errorf("GetListServicesResponse() error")
	} else {
		module_services := modules.GetService()
		module_names := make([]string, len(module_services))
		for i, service := range module_services {
			module_names[i] = service.Name
		}
		return module_names, nil
	}
}
