package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

type Server struct {
	gopi.Unit
	sync.RWMutex
	sync.WaitGroup

	cert, key *string
	ssl       bool
	errs      chan error
	server    *http.Server
	timeout   *time.Duration
}

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

	// Create channel for errors
	this.errs = make(chan error)

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

	// Close error channel
	close(this.errs)

	// Release resources
	this.server = nil
	this.errs = nil

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

// Start serves HTTP in foreground
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
		if this.ssl {
			this.errs <- this.server.ListenAndServeTLS(*this.cert, *this.key)
		} else {
			this.errs <- this.server.ListenAndServe()
		}
		this.server = nil
		this.WaitGroup.Done()
	}()

	// Return success
	return nil
}

func (this *Server) RegisterService(interface{}, gopi.Service) error {
	return gopi.ErrNotImplemented
}

func (this *Server) Stop(force bool) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If not started, return nil
	if this.server == nil {
		return nil
	}

	cancel := func() {}
	ctx := context.Background()
	if force {
		ctx, cancel = context.WithTimeout(ctx, time.Second)
	}
	defer cancel()
	err := this.server.Shutdown(ctx)
	this.WaitGroup.Wait()
	return err
}

func (this *Server) NewStreamContext() context.Context {
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
