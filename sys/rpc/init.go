/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register rpc/server
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/server",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("rpc.port", 0, "Server Port")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			port, _ := app.AppFlags.GetUint("rpc.port")
			return gopi.Open(Server{Port: port}, app.Logger)
		},
	})

	// Register rpc/discovery
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/discovery",
		Type: gopi.MODULE_TYPE_MDNS,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("mdns.domain", MDNS_DEFAULT_DOMAIN, "Domain")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			domain, _ := app.AppFlags.GetString("mdns.domain")
			return gopi.Open(Config{Domain: domain}, app.Logger)
		},
	})
}
