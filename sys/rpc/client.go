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
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	reflection_pb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// Client is the RPC client configuration
type Client struct {
	Host       string
	SSL        bool
	SkipVerify bool
	Timeout    time.Duration // Connection timeout
}

type client struct {
	log        gopi.Logger
	host       string
	ssl        bool
	skipverify bool
	timeout    time.Duration
	conn       *grpc.ClientConn
}

////////////////////////////////////////////////////////////////////////////////
// CLIENT OPEN AND CLOSE

// Open a client
func (config Client) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<grpc.Client>Open(host=%v,ssl=%v,skipverify=%v,timeout=%v)", config.Host, config.SSL, config.SkipVerify, config.Timeout)

	// Create a client object
	this := new(client)
	this.host = config.Host
	this.ssl = config.SSL
	this.skipverify = config.SkipVerify
	this.timeout = config.Timeout
	this.log = log
	this.conn = nil

	// success
	return this, nil
}

// Close client
func (this *client) Close() error {
	this.log.Debug2("<grpc.Client>Close()")
	return this.Disconnect()
}

////////////////////////////////////////////////////////////////////////////////
// RPCClient interface implementation

func (this *client) Connect() error {
	if this.conn != nil {
		return errors.New("Cannot call Connect() when connection already made")
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
	if conn, err := grpc.Dial(this.host, opts...); err != nil {
		return err
	} else {
		this.conn = conn
	}

	// Success
	return nil
}

func (this *client) Disconnect() error {
	if this.conn != nil {
		this.log.Debug("<grpc.Client>Disconnect()")
		err := this.conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

func (this *client) Modules() ([]string, error) {
	if this.conn == nil {
		return nil, errors.New("Disconnected")
	}
	if client := this.newServerReflectionClient(); client == nil {
		return nil, errors.New("Unable to create client")
	} else {
		defer client.CloseSend()
		if services, err := this.listServices(client); err != nil {
			return nil, err
		} else {
			return services, nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *client) String() string {
	if this.conn != nil {
		return fmt.Sprintf("<grpc.Client>{ connected=true host=%v ssl=%v skipverify=%v }", this.host, this.ssl, this.skipverify)
	} else {
		return fmt.Sprintf("<grpc.Client>{ connected=false host=%v ssl=%v skipverify=%v }", this.host, this.ssl, this.skipverify)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *client) newServerReflectionClient() reflection_pb.ServerReflection_ServerReflectionInfoClient {
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

func (this *client) listServices(c reflection_pb.ServerReflection_ServerReflectionInfoClient) ([]string, error) {
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
