package tool

import (
	"context"
	"net"
	"strconv"
	"strings"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type server struct {
	gopi.Unit
	gopi.Server
	gopi.Logger
	gopi.ServiceDiscovery

	addr, name *string
	version    string
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func Server(name string, args []string, objs ...interface{}) int {
	srv := []interface{}{new(server)}
	return CommandLine(name, args, append(srv, objs...)...)
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *server) Define(cfg gopi.Config) error {
	this.addr = cfg.FlagString("addr", "", "Address for server")
	this.name = cfg.FlagString("name", "", "Service name")
	return nil
}

func (this *server) New(cfg gopi.Config) error {
	// Check to make sure server is available
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Server")
	}

	// Set defaults for name and version
	*this.name = strings.TrimSpace(*this.name)
	if *this.name == "" {
		*this.name = cfg.Version().Name()
	}
	this.version, _, _ = cfg.Version().Version()

	// Return success
	return nil
}

func (this *server) Run(ctx context.Context) error {
	// Start server in background
	if err := this.Server.StartInBackground("", *this.addr); err != nil {
		return err
	} else {
		this.Debug("Started server: ", this.Server)
	}

	// Determine port
	port := uint16(0)
	if _, port_, err := net.SplitHostPort(this.Server.Addr()); err == nil {
		if port_, err := strconv.ParseUint(port_, 0, 16); err == nil {
			port = uint16(port_)
		}
	}

	// Service
	service := this.Server.Service()

	// Set TXT record
	txt := []string{}
	if this.Server.Flags()&gopi.SERVICE_FLAG_TLS != 0 {
		txt = append(txt, "ssl=1")
	} else {
		txt = append(txt, "ssl=0")
	}
	if this.version != "" {
		txt = append(txt, "v="+this.version)
	}

	// Register if ServiceDisovery is enabled and not socket-based
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	if port != 0 {
		if this.ServiceDiscovery != nil {
			record, err := this.ServiceDiscovery.NewServiceRecord(service, *this.name, port, txt, 0)
			if err != nil {
				this.Debug("Error: ", err)
			} else {
				this.Debug("Advertising server: ", record)
			}
			go func() {
				if this.ServiceDiscovery != nil && record != nil {
					if err := this.ServiceDiscovery.Serve(ctx2, []gopi.ServiceRecord{record}); err != nil {
						this.Print("Error: ", err)
					}
				}
			}()
		} else {
			this.Debug("Notice: ServiceDiscovery is not enabled")
		}
	}

	// Wait for interupt
	<-ctx.Done()

	// Return success
	return nil
}
