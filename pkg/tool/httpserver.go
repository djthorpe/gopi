package tool

import (
	"context"
	"net"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/http"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type server struct {
	gopi.Unit
	gopi.Server
	gopi.Logger
	gopi.ServiceDiscovery

	addr *string
	name string
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func HttpServer(name string, args []string, objs ...interface{}) int {
	srv := []interface{}{new(server)}
	return CommandLine(name, args, append(srv, objs...)...)
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *server) Define(cfg gopi.Config) error {
	this.addr = cfg.FlagString("addr", ":0", "Address for HTTP Server")
	return nil
}

func (this *server) New(cfg gopi.Config) error {
	this.name = cfg.Version().Name()
	return nil
}

func (this *server) Run(ctx context.Context) error {
	addr := *this.addr
	if err := this.Server.StartInBackground("tcp", addr); err != nil {
		return err
	}

	// Determine port
	port := uint16(0)
	if _, port_, err := net.SplitHostPort(this.Server.Addr()); err == nil {
		if port_, err := strconv.ParseUint(port_, 0, 16); err == nil {
			port = uint16(port_)
		}
	}

	// Register if ServiceDisovery is enabled
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	if this.ServiceDiscovery != nil && port != 0 {
		go func() {
			if record, err := this.ServiceDiscovery.NewServiceRecord("_http._tcp.", this.name, port, nil, 0); err != nil {
				this.Debug(err)
			} else if err := this.ServiceDiscovery.Serve(ctx2, []gopi.ServiceRecord{record}); err != nil {
				this.Debug(err)
			}
		}()
	} else {
		this.Debug("ServiceDiscovery is not enabled")
	}

	// Wait for interupt
	this.Debug("Started server, http://localhost" + this.Server.Addr() + "/")
	<-ctx.Done()

	// Return success
	return nil
}
