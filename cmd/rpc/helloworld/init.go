/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package helloworld

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register service/helloworld:grpc
	gopi.RegisterModule(gopi.Module{
		Name:     "service/helloworld:grpc",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Service{
				Server: app.ModuleInstance("rpc/server").(gopi.RPCServer),
			}, app.Logger)
		},
	})

	// Register client/helloworld:grpc
	gopi.RegisterModule(gopi.Module{
		Name: "client/helloworld:grpc",
		Type: gopi.MODULE_TYPE_CLIENT,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return nil, gopi.ErrNotImplemented
		},
	})
}
