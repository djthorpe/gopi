package googlecast

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/djthorpe/gopi"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type connection struct {
	sync.RWMutex
	channel

	conn *tls.Conn
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *connection) Connect(addr string, timeout time.Duration) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If already connected, return
	if this.conn != nil {
		return gopi.ErrOutOfOrder.WithPrefix("Connect")
	}

	// TODO start the receive loop somewhere

	// Get a connection address and port, connect and establish
	// the channel
	if conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
	}, "tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	}); err != nil {
		return fmt.Errorf("%s: %w", addr, err)
	} else if _, data, err := this.channel.Connect(); err != nil {
		conn.Close()
		return err
	} else if err := this.send(data); err != nil {
		conn.Close()
		return err
	} else {
		this.conn = conn
	}

	// Success
	return nil
}

func (this *connection) Disconnect() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If already connected, return
	if this.conn == nil {
		return nil
	}

	var result error

	// TODO Send close message

	// TODO stop receive loop

	// Close connection
	if err := this.conn.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	// Release resources
	this.conn = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *connection) Addr() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if this.conn != nil {
		return this.conn.RemoteAddr().String()
	} else {
		return ""
	}
}

func (this *connection) IsConnected() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.conn != nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *connection) String() string {
	str := "<cast.conn"
	if addr := this.Addr(); addr != "" {
		str += " addr=" + addr
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *connection) send(data []byte) error {
	if len(data) == 0 {
		return nil
	} else if err := binary.Write(this.conn, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	} else if _, err := this.conn.Write(data); err != nil {
		return err
	} else {
		return nil
	}
}
