/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register service/helloworld:grpc
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/service/metrics:grpc",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server", "metrics"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Service{
				Server:  app.ModuleInstance("rpc/server").(gopi.RPCServer),
				Metrics: app.ModuleInstance("metrics").(gopi.Metrics),
			}, app.Logger)
		},
	})

	// Register the client
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/client/metrics:grpc",
		Type:     gopi.MODULE_TYPE_CLIENT,
		Requires: []string{"rpc/clientpool"},
		Run: func(app *gopi.AppInstance, _ gopi.Driver) error {
			clientpool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
			if clientpool == nil {
				return gopi.ErrAppError
			} else {
				clientpool.RegisterClient("mutablelogic.Metrics", NewClient)
				return nil
			}
		},
	})
}
