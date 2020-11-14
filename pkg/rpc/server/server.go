package server

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
	grpc "google.golang.org/grpc"
	reflection "google.golang.org/grpc/reflection"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type server struct {
	gopi.Unit
	sync.Mutex

	srv      *grpc.Server
	listener net.Listener
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *server) New(cfg gopi.Config) error {
	// Create a gRPC server
	if server := grpc.NewServer(); server == nil {
		return gopi.ErrBadParameter
	} else {
		this.srv = server
	}

	// Register reflection
	reflection.Register(this.srv)

	// Return success
	return nil
}

func (this *server) Dispose() error {
	var result error

	// force stop the server
	if this.listener != nil {
		if err := this.Stop(true); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.listener = nil
	this.srv = nil

	// Return success
	return result
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *server) StartInBackground(network, addr string) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Create listener
	if this.listener != nil {
		return gopi.ErrOutOfOrder
	} else if listener, err := net.Listen(network, addr); err != nil {
		return err
	} else {
		this.listener = listener
	}

	// Serve!
	go func() {
		if err := this.srv.Serve(this.listener); err != nil {
			// Should emit any errors on a channel
			fmt.Fprintln(os.Stderr, "TODO", err)
		}
	}()

	// Return success
	return nil
}

func (this *server) Stop(force bool) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check for listener
	if this.listener == nil {
		return nil
	}

	// Perform stop
	if force {
		this.srv.Stop()
	} else {
		this.srv.GracefulStop()
	}

	// Close listener
	this.listener = nil

	// Return success
	return nil
}

func (this *server) RegisterService(fn interface{}, service gopi.Service) error {
	// Check parameters
	if fn == nil {
		return gopi.ErrBadParameter.WithPrefix("fn")
	}
	if service == nil {
		return gopi.ErrBadParameter.WithPrefix("service")
	}
	if value := reflect.ValueOf(fn); value.Kind() != reflect.Func {
		return gopi.ErrBadParameter.WithPrefix("fn")
	} else {
		value.Call([]reflect.Value{reflect.ValueOf(this.srv), reflect.ValueOf(service)})
	}

	// Return success
	return nil
}

func (this *server) Addr() string {
	if this.listener != nil {
		return this.listener.Addr().String()
	} else {
		return ""
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *server) String() string {
	str := "<server"
	if this.listener != nil {
		str += fmt.Sprintf(" addr=%q", this.listener.Addr())
	}

	for k, v := range this.srv.GetServiceInfo() {
		str += " " + k + "=["
		for i, method := range v.Methods {
			if i > 0 {
				str += ","
			}
			str += method.Name
		}
		str += "]"
	}

	return str + ">"
}
