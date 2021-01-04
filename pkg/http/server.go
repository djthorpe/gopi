package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Server struct {
	gopi.Unit
	gopi.Logger
	sync.RWMutex
	sync.WaitGroup

	cert, key *string
	ssl       bool
	server    *http.Server
	mux       *http.ServeMux
	timeout   *time.Duration
	handler   http.Handler
}

type Transport interface {
	http.Handler
	SetHandler(fn http.Handler)
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Server) Define(cfg gopi.Config) error {
	this.cert = cfg.FlagString("ssl.cert", "", "SSL Certificate")
	this.key = cfg.FlagString("ssl.key", "", "SSL Key")
	this.timeout = cfg.FlagDuration("http.timeout", 15*time.Second, "HTTP server read and write timeout")
	return nil
}

func (this *Server) New(cfg gopi.Config) error {
	// Check SSL parameters
	if *this.cert != "" || *this.key != "" {
		if _, err := os.Stat(*this.cert); os.IsNotExist(err) {
			return fmt.Errorf("Invalid SSL certificate")
		} else if _, err := os.Stat(*this.key); os.IsNotExist(err) {
			return fmt.Errorf("Invalid SSL private key")
		} else {
			this.ssl = true
		}
	}

	// Set multiplexer and handler chain
	this.mux = http.NewServeMux()

	// Return success
	return nil
}

func (this *Server) Dispose() error {
	var result error
	if err := this.Stop(true); err != nil {
		result = multierror.Append(result, err)
	}

	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Release resources
	this.server = nil
	this.mux = nil

	// Return any errors
	return result
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Server) Addr() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.server == nil {
		return ""
	} else {
		return this.server.Addr
	}
}

func (this *Server) SSL() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.server == nil {
		return false
	} else {
		return this.ssl
	}
}

func (this *Server) Service() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.server != nil {
		return "_http._tcp"
	} else {
		return ""
	}
}

// Start serves HTTP in foreground. Network should always be "tcp"
// and address is either empty (using standard ports) or ":0" which
// means a free port is used and can be determined using the Addr
// method once the server has started.
func (this *Server) StartInBackground(network, addr string) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If already started, return out of order error
	if this.server != nil {
		return gopi.ErrOutOfOrder.WithPrefix("StartInBackground")
	}

	// Only accept "tcp" as network argument
	if network != "tcp" {
		return gopi.ErrBadParameter.WithPrefix("StartInBackground", network)
	}
	if addr == "" {
		if this.ssl {
			addr = ":https"
		} else {
			addr = ":http"
		}
	} else if addr == ":0" {
		if port, err := getFreePort(); err != nil {
			return err
		} else {
			addr = ":" + fmt.Sprint(port)
		}
	}

	// Set server object
	this.server = &http.Server{
		Addr:              addr,
		Handler:           this,
		ReadHeaderTimeout: *this.timeout,
		WriteTimeout:      *this.timeout,
		IdleTimeout:       *this.timeout,
	}

	// Start server in background
	this.WaitGroup.Add(1)
	go func() {
		var result error
		if this.ssl {
			result = this.server.ListenAndServeTLS(*this.cert, *this.key)
		} else {
			result = this.server.ListenAndServe()
		}
		if errors.Is(result, http.ErrServerClosed) == false {
			this.Print(result)
		}
		this.WaitGroup.Done()
	}()

	// Return success
	return nil
}

func (this *Server) Stop(force bool) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If not started, return nil
	if this.server == nil {
		return nil
	}

	// Close or Shutdown
	var result error
	if force {
		if err := this.server.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), *this.timeout)
		defer cancel()
		if err := this.server.Shutdown(ctx); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Wait for server stop
	this.WaitGroup.Wait()

	// Set server object as nil and return any errors
	this.server = nil

	return result
}

// NewStreamContext is unused presently as it's not so useful for HTTP
func (this *Server) NewStreamContext() context.Context {
	return nil
}

// RegisterService currently accepts a path and a http.Handler object
// but in future should also be able to handle http.Transport handlers
// as well
func (this *Server) RegisterService(path interface{}, service gopi.Service) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.mux == nil {
		return gopi.ErrOutOfOrder.WithPrefix("RegisterService")
	} else if handler_, ok := service.(http.Handler); ok == false {
		return gopi.ErrBadParameter.WithPrefix("RegisterService", "service")
	} else if path == nil {
		if handler_, ok := service.(Transport); ok == false {
			return gopi.ErrBadParameter.WithPrefix("RegisterService", "Does not implement SetHandler")
		} else {
			if this.handler == nil {
				handler_.SetHandler(this.mux)
			} else {
				handler_.SetHandler(this.handler)
			}
			this.handler = handler_
		}
	} else {
		if path_, ok := path.(string); ok == false {
			return gopi.ErrBadParameter.WithPrefix("RegisterService", "path")
		} else {
			this.mux.Handle(path_, handler_)
		}
	}

	// Return success
	return nil
}

func (this *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// If any handlers are installed call them, or else call the default multiplexer
	if this.handler == nil {
		this.mux.ServeHTTP(w, req)
	} else {
		this.handler.ServeHTTP(w, req)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Server) String() string {
	str := "<server.http"
	if addr := this.Addr(); addr != "" {
		str += " addr=" + strconv.Quote(addr)
	}
	return str + ">"
}
