package client

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type connpool struct {
	gopi.Unit
	sync.Mutex

	conns []gopi.Conn
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *connpool) Dispose() error {
	var result error

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close all clients
	for _, c := range this.conns {
		if c != nil {
			if err := c.(*conn).Close(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Return success
	return result
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *connpool) Connect(network, addr string) (gopi.Conn, error) {
	switch network {
	case "tcp":
		if conn, err := grpc.Dial(addr, grpc.WithInsecure()); err != nil {
			return nil, err
		} else if client := NewConn(conn); client == nil {
			return nil, gopi.ErrInternalAppError.WithPrefix(addr)
		} else {
			this.Mutex.Lock()
			defer this.Mutex.Unlock()
			this.conns = append(this.conns, client)
			return client, nil
		}
	default:
		return nil, gopi.ErrNotImplemented.WithPrefix(network)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *connpool) String() string {
	str := "<connpool"
	str += " conns=" + fmt.Sprint(this.conns)
	return str + ">"
}
