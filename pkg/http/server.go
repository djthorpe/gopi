package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	fcgi "github.com/djthorpe/gopi/v3/pkg/http/fcgi"
	multierror "github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Server struct {
	gopi.Unit
	gopi.Logger
	sync.RWMutex
	sync.WaitGroup

	cert, key  *string
	fcgi       *bool
	ssl        bool
	httpserver *http.Server
	fcgiserver *fcgi.Server
	mux        *http.ServeMux
	timeout    *time.Duration
	handler    http.Handler
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
	this.fcgi = cfg.FlagBool("http.fcgi", false, "Serve over FastCGI unix socket")
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

	// FCGI will not work with SSL
	if *this.fcgi && this.ssl {
		return fmt.Errorf("SSL and FCGI are not compatible")
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
	this.httpserver = nil
	this.fcgiserver = nil
	this.mux = nil

	// Return any errors
	return result
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Server) Addr() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.httpserver != nil {
		return this.httpserver.Addr
	} else if this.fcgiserver != nil {
		return this.fcgiserver.Addr
	} else {
		return ""
	}
}

func (this *Server) Flags() gopi.ServiceFlag {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	f := gopi.SERVICE_FLAG_NONE
	if this.httpserver != nil {
		f |= gopi.SERVICE_FLAG_HTTP
	}
	if this.ssl {
		f |= gopi.SERVICE_FLAG_TLS
	}
	if this.fcgiserver != nil {
		f |= gopi.SERVICE_FLAG_FCGI
		if this.fcgiserver.Network == "unix" {
			f |= gopi.SERVICE_FLAG_SOCKET
		}
	}
	return f
}

func (this *Server) Service() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if f := this.Flags(); f&gopi.SERVICE_FLAG_HTTP != 0 {
		return "_http._tcp"
	} else if f&gopi.SERVICE_FLAG_FCGI != 0 {
		return "_fcgi._tcp"
	} else {
		return ""
	}
}

// Start serves HTTP in foreground. Network should be empty or "unix"
// for FGCI or "tcp" otherwise, In the former case, the address is
// either empty (using standard ports) or ":0" which means a free
// port is used and can be determined using the Addr method once the
// server has started. In the case for fcgi it should be a filename
// so that a unix socket can be created for communication.
func (this *Server) StartInBackground(network, addr string) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If already started, return out of order error
	if this.httpserver != nil || this.fcgiserver != nil {
		return gopi.ErrOutOfOrder.WithPrefix("StartInBackground")
	}

	// Check network parameter
	if network == "" {
		// Set network from addr
		if _, _, err := net.SplitHostPort(addr); err == nil {
			network = "tcp"
		}
	} else if *this.fcgi && (network != "unix" && network != "tcp") {
		return gopi.ErrBadParameter.WithPrefix("StartInBackground: ", network)
	} else if network != "tcp" {
		return gopi.ErrBadParameter.WithPrefix("StartInBackground: ", network)
	}
	// Set addr parameter if TCP
	if network == "tcp" {
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
	}
	// Create server
	if *this.fcgi {
		this.fcgiserver = &fcgi.Server{
			Network: network,
			Addr:    addr,
			Handler: this,
		}
	} else {
		this.httpserver = &http.Server{
			Addr:              addr,
			Handler:           this,
			ReadHeaderTimeout: *this.timeout,
			WriteTimeout:      *this.timeout,
			IdleTimeout:       *this.timeout,
		}
	}

	// Timeout for server is 500ms
	errs := make(chan error)

	// Start server in background
	this.WaitGroup.Add(1)
	go func() {
		var result error
		if this.fcgiserver != nil {
			result = this.fcgiserver.ListenAndServe()
		} else if this.ssl {
			result = this.httpserver.ListenAndServeTLS(*this.cert, *this.key)
		} else {
			result = this.httpserver.ListenAndServe()
		}
		if errors.Is(result, http.ErrServerClosed) == false {
			// Pass the error condition
			select {
			case errs <- result:
				break
			default:
				this.Debug("StartInBackground: ", result)
			}
		}
		this.WaitGroup.Done()
	}()

	// Wait for an error to occur in the server for 0.5s or else
	// return success
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	select {
	case err := <-errs:
		return err
	case <-timer.C:
		return nil
	}
}

func (this *Server) Stop(force bool) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If not started, return nil
	if this.httpserver == nil && this.fcgiserver == nil {
		return nil
	}

	// Close or Shutdown
	var result error
	if this.fcgiserver != nil {
		if err := this.fcgiserver.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.httpserver != nil {
		if force {
			if err := this.httpserver.Close(); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), *this.timeout)
			defer cancel()
			if err := this.httpserver.Shutdown(ctx); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Wait for server stop
	this.WaitGroup.Wait()

	// Set server objects as nil and return any errors
	this.httpserver = nil
	this.fcgiserver = nil

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
		return gopi.ErrBadParameter.WithPrefix("RegisterService: ", "service")
	} else if path == nil {
		if handler_, ok := service.(Transport); ok == false {
			return gopi.ErrBadParameter.WithPrefix("RegisterService: ", "Does not implement SetHandler")
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
			return gopi.ErrBadParameter.WithPrefix("RegisterService: ", "path")
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
	if f := this.Flags(); f != 0 {
		str += " flags=" + fmt.Sprint(f)
	}
	if addr := this.Addr(); addr != "" {
		str += " addr=" + strconv.Quote(addr)
	}
	return str + ">"
}
