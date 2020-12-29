package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"reflect"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	reflection "google.golang.org/grpc/reflection"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type server struct {
	gopi.Unit
	sync.Mutex
	gopi.Logger

	srv      *grpc.Server
	listener net.Listener
	ssl      bool
	cancels  []context.CancelFunc
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *server) Define(cfg gopi.Config) error {
	cfg.FlagString("ssl.cert", "", "SSL certificate file")
	cfg.FlagString("ssl.key", "", "SSL key file")
	cfg.FlagDuration("timeout", 0, "Connection timeout")
	return nil
}

func (this *server) New(cfg gopi.Config) error {
	opts := []grpc.ServerOption{}
	if opts, ssl, err := appendServerCredentialOption(cfg, opts); err != nil {
		return err
	} else if opts, err := appendConnectionTimeoutOption(cfg, opts); err != nil {
		return err
	} else if server := grpc.NewServer(opts...); server == nil {
		return gopi.ErrBadParameter
	} else {
		this.srv = server
		this.ssl = ssl
	}

	// Register reflection service
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

	// Send cancels
	for _, cancel := range this.cancels {
		cancel()
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
	if this.Logger != nil {
		this.Logger.Debug("RegisterService: ", reflect.TypeOf(service))
	}

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

func (this *server) NewStreamContext() context.Context {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	this.cancels = append(this.cancels, cancel)
	return ctx
}

func (this *server) Addr() string {
	if this.listener != nil {
		return this.listener.Addr().String()
	} else {
		return ""
	}
}

func (this *server) SSL() bool {
	if this.listener != nil {
		return this.ssl
	} else {
		return false
	}
}

func (this *server) Service() string {
	if this.listener != nil {
		return "_grpc._" + this.listener.Addr().Network()
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

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func appendServerCredentialOption(cfg gopi.Config, opts []grpc.ServerOption) ([]grpc.ServerOption, bool, error) {
	cert := cfg.GetString("ssl.cert")
	key := cfg.GetString("ssl.key")
	ssl := false
	if cert != "" || key != "" {
		if creds, err := credentials.NewServerTLSFromFile(cert, key); err != nil {
			return nil, false, err
		} else {
			opts = append(opts, grpc.Creds(creds))
			ssl = true
		}
	}
	return opts, ssl, nil
}

func appendConnectionTimeoutOption(cfg gopi.Config, opts []grpc.ServerOption) ([]grpc.ServerOption, error) {
	if timeout := cfg.GetDuration("timeout"); timeout > 0 {
		opts = append(opts, grpc.ConnectionTimeout(timeout))
	}
	return opts, nil
}
