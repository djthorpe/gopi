/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package grpc

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register GRPC rpc/server
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/server",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("rpc.port", 0, "Server Port")
			config.AppFlags.FlagString("rpc.sslcert", "", "SSL Certificate Path")
			config.AppFlags.FlagString("rpc.sslkey", "", "SSL Key Path")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			port, _ := app.AppFlags.GetUint("rpc.port")
			key, _ := app.AppFlags.GetString("rpc.sslkey")
			cert, _ := app.AppFlags.GetString("rpc.sslcert")
			return gopi.Open(Server{
				Port:           port,
				SSLCertificate: cert,
				SSLKey:         key,
				ServerOption:   []grpc.ServerOption{},
			}, app.Logger)
		},
	})

	// Register GRPC rpc/clientpool module
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/clientpool",
		Type:     gopi.MODULE_TYPE_OTHER,
		Requires: []string{"mdns"},
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagBool("rpc.insecure", false, "Disable SSL Connection")
			config.AppFlags.FlagBool("rpc.skipverify", true, "Skip SSL Verification")
			config.AppFlags.FlagDuration("rpc.timeout", 0, "Connection timeout")
			config.AppFlags.FlagString("rpc.service", "", "Comma-separated list of service names")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			insecure, _ := app.AppFlags.GetBool("rpc.insecure")
			skipverify, _ := app.AppFlags.GetBool("rpc.skipverify")
			timeout, _ := app.AppFlags.GetDuration("rpc.timeout")
			service, _ := app.AppFlags.GetString("rpc.service")
			if service == "" {
				service = app.Service()
			}
			return gopi.Open(ClientPool{
				Discovery:  app.ModuleInstance("mdns").(gopi.RPCServiceDiscovery),
				SkipVerify: skipverify,
				SSL:        (insecure == false),
				Timeout:    timeout,
				Service:    service,
			}, app.Logger)
		},
	})

	/*
		// Register GRPC rpc/clientconn
		gopi.RegisterModule(gopi.Module{
			Name: "rpc/clientconn",
			Type: gopi.MODULE_TYPE_OTHER,
			Config: func(config *gopi.AppConfig) {
				config.AppFlags.FlagString("rpc.addr", "", "Server Address")
				config.AppFlags.FlagBool("rpc.insecure", false, "Disable SSL Connection")
				config.AppFlags.FlagBool("rpc.skipverify", true, "Skip SSL Verification")
				config.AppFlags.FlagDuration("rpc.timeout", 0, "Connection timeout")
			},
			New: func(app *gopi.AppInstance) (gopi.Driver, error) {
				addr, _ := app.AppFlags.GetString("rpc.addr")
				insecure, _ := app.AppFlags.GetBool("rpc.insecure")
				skipverify, _ := app.AppFlags.GetBool("rpc.skipverify")
				timeout, _ := app.AppFlags.GetDuration("rpc.timeout")
				return gopi.Open(ClientConn{
					Addr:       addr,
					SSL:        (insecure == false),
					SkipVerify: skipverify,
					Timeout:    timeout,
				}, app.Logger)
			},
		})*/
}
