/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"fmt"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	evt "github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ClientPool struct {
	SkipVerify bool
	SSL        bool
	Timeout    time.Duration
}

type clientpool struct {
	log        gopi.Logger
	skipverify bool
	ssl        bool
	timeout    time.Duration
	pubsub     evt.PubSub
	clients    map[string]gopi.RPCNewClientFunc
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ClientPool) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<rpc.clientpool>Open{ SSL=%v skipverify=%v timeout=%v }", config.SSL, config.SkipVerify, config.Timeout)
	this := new(clientpool)
	this.log = log
	this.skipverify = config.SkipVerify
	this.ssl = config.SSL
	this.timeout = config.Timeout
	this.pubsub = evt.NewPubSub(0)
	this.clients = make(map[string]gopi.RPCNewClientFunc)

	// Success
	return this, nil
}

func (this *clientpool) Close() error {
	this.log.Debug("<rpc.clientpool>Close{}")

	// Release resources
	this.pubsub.Close()
	this.pubsub = nil
	this.clients = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *clientpool) Connect(service *gopi.RPCServiceRecord, flags gopi.RPCFlag) (gopi.RPCClientConn, error) {
	this.log.Debug2("<rpc.clientpool>Connect{ service=%v flags=%v }", service, flags)

	// Determine the address
	if addr := addressFor(service, flags); addr == "" {
		return nil, gopi.ErrBadParameter
	} else if clientconn_, err := gopi.Open(ClientConn{
		Name:       service.Name,
		Addr:       addr + ":" + fmt.Sprint(service.Port),
		SSL:        this.ssl,
		SkipVerify: this.skipverify,
		Timeout:    this.timeout,
	}, this.log); err != nil {
		return nil, err
	} else if clientconn, ok := clientconn_.(gopi.RPCClientConn); ok == false {
		return nil, gopi.ErrOutOfOrder
	} else {
		// Do connection
		if err := clientconn.Connect(); err != nil {
			return nil, err
		}

		// TODO: Emit a connected event

		// Return success
		return clientconn, nil
	}
}

func (this *clientpool) Disconnect(conn gopi.RPCClientConn) error {
	this.log.Debug2("<rpc.clientpool>Disconnect{ conn=%v }", conn)

	// TODO: Emit a disconnect event

	return conn.Disconnect()
}

////////////////////////////////////////////////////////////////////////////////
// CLIENTS

func (this *clientpool) RegisterClient(service string, callback gopi.RPCNewClientFunc) error {
	this.log.Debug2("<rpc.clientpool>RegisterClient{ service=%v func=%v }", service, callback)
	if service == "" || callback == nil {
		return gopi.ErrBadParameter
	} else if _, exists := this.clients[service]; exists {
		this.log.Debug("<rpc.clientpool>RegisterClient: Duplicate service: %v", service)
		return gopi.ErrBadParameter
	} else {
		this.clients[service] = callback
		return nil
	}
}

func (this *clientpool) NewClient(service string, conn mutablehome.RPCClientConn) (mutablehome.RPCClient, error) {
	this.log.Debug2("<rpc.clientpool>NewClient{ service=%v conn=%v }", service, conn)

	// Obtain the module with which to create a new client
	if callback, exists := this.clients[service]; exists == false {
		return nil, mutablehome.ErrServiceNotRegistered
	} else {
		return callback(conn)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBSUB

func (this *clientpool) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

func (this *clientpool) Unsubscribe(c <-chan gopi.Event) {
	this.pubsub.Unsubscribe(c)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func addressFor(service *gopi.RPCServiceRecord, flags gopi.RPCFlag) string {
	if flags&gopi.RPC_FLAG_INET_UDP != 0 {
		// We don't support UDP connections
		return ""
	} else if flags&gopi.RPC_FLAG_INET_V6 != 0 {
		if len(service.IP6) == 0 {
			return ""
		} else {
			return service.IP6[0].String()
		}
	} else if flags&gopi.RPC_FLAG_INET_V4 != 0 {
		if len(service.IP4) == 0 {
			return ""
		} else {
			return service.IP4[0].String()
		}
	} else {
		return service.Host
	}
}
