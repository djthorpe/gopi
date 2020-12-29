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
	timeout   *time.Duration
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
		Handler:           http.NewServeMux(),
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

func (this *Server) NewStreamContext() context.Context {
	return nil
}

func (this *Server) RegisterService(path interface{}, service gopi.Service) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.server == nil {
		return gopi.ErrOutOfOrder.WithPrefix("RegisterService")
	} else if path_, ok := path.(string); ok == false {
		return gopi.ErrBadParameter.WithPrefix("RegisterService", "path")
	} else if handler_, ok := service.(http.Handler); ok == false {
		return gopi.ErrBadParameter.WithPrefix("RegisterService", "service")
	} else {
		mux := this.server.Handler.(*http.ServeMux)
		mux.Handle(path_, handler_)
	}

	// Return success
	return nil
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
