package tool

import (
	"context"
	"net"
	"strconv"
	"strings"

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

	addr, name *string
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
	this.name = cfg.FlagString("name", "", "HTTP Service Name")
	return nil
}

func (this *server) New(cfg gopi.Config) error {
	*this.name = strings.TrimSpace(*this.name)
	if *this.name == "" {
		*this.name = cfg.Version().Name()
	}
	return this.Server.StartInBackground("tcp", *this.addr)
}

func (this *server) Run(ctx context.Context) error {
	// Determine port
	port := uint16(0)
	if _, port_, err := net.SplitHostPort(this.Server.Addr()); err == nil {
		if port_, err := strconv.ParseUint(port_, 0, 16); err == nil {
			port = uint16(port_)
		}
	}

	// Set TXT record
	txt := []string{}
	if this.Server.SSL() {
		txt = append(txt, "ssl=1")
	} else {
		txt = append(txt, "ssl=0")
	}

	// Register if ServiceDisovery is enabled
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	if this.ServiceDiscovery != nil && port != 0 {
		record, err := this.ServiceDiscovery.NewServiceRecord("_http._tcp.", *this.name, port, txt, 0)
		if err != nil {
			this.Debug(err)
		} else {
			this.Debug("Started server: ", record)
		}
		go func() {
			if this.ServiceDiscovery != nil && record != nil {
				if err := this.ServiceDiscovery.Serve(ctx2, []gopi.ServiceRecord{record}); err != nil {
					this.Debug(err)
				}
			}
		}()
	} else {
		this.Debug("ServiceDiscovery is not enabled")
	}

	// Wait for interupt
	<-ctx.Done()

	// Return success
	return nil
}
