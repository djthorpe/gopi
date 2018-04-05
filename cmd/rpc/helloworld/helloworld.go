/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package helloworld

import "github.com/djthorpe/gopi"

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register rpc/service:helloworld
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/service:helloworld",
		Type: gopi.MODULE_TYPE_SERVICE,
	})
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE
