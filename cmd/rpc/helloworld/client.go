/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package helloworld

import (
	"fmt"
	// Framework
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
}

type client struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the client connection
func (config Client) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.client.helloworld>Open{}")

	this := new(client)
	this.log = log

	// Success
	return this, nil
}

// Close the client connection
func (this *client) Close() error {
	this.log.Debug("<grpc.client.helloworld>Close{}")

	// No resources to release

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *client) String() string {
	return fmt.Sprintf("grpc.client.helloworld{}")
}
