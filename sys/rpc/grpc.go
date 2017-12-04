/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/third_party/grpc-go"
	_ "github.com/djthorpe/gopi/third_party/grpc-go/reflection"
)

// Server is the RPC server configuration
type Server struct{}

type server struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// SERVER OPEN AND CLOSE

// Open a logger
func (config Server) Open(log gopi.Logger) (gopi.Driver, error) {

	this := new(server)
	this.log = log

	// success
	return this, nil
}

// Close a logger
func (this *server) Close() error {
	// Return success
	return nil
}
