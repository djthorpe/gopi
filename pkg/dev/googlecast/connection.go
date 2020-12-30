package googlecast

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type connection struct {
	sync.RWMutex
	sync.WaitGroup
	channel

	conn   *tls.Conn
	cancel context.CancelFunc
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	recvTimeout = 500 * time.Millisecond
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *connection) Connect(key, addr string, timeout time.Duration, errs chan<- error, state chan<- state) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If already connected, return
	if this.conn != nil {
		return gopi.ErrOutOfOrder.WithPrefix("Connect")
	}

	// Get a connection address and port, connect and establish
	// the channel
	if conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
	}, "tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	}); err != nil {
		return fmt.Errorf("%s: %w", addr, err)
	} else {
		this.conn = conn
	}

	// Initialise the channel
	this.channel.Init(key, state)

	// Start the receive loop, which will end on cancel()
	this.WaitGroup.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		defer this.WaitGroup.Done()
		this.recv(ctx, errs)
	}(ctx)
	this.cancel = cancel

	// Send CONNECT message. Errors here require disconnect
	if _, data, err := this.channel.Connect(); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
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

	// Send CLOSE message
	if _, data, err := this.channel.Disconnect(); err != nil {
		result = multierror.Append(result, err)
	} else if err := this.send(data); err != nil {
		result = multierror.Append(result, err)
	}

	// End receive loop and wait
	this.cancel()
	this.WaitGroup.Wait()

	// Close connection
	if err := this.conn.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	// Release resources
	this.conn = nil
	this.cancel = nil

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

func (this *connection) recv(ctx context.Context, errs chan<- error) {
	var length uint32
	for {
		select {
		case <-ctx.Done():
			return
		default:
			timeout := time.Now().Add(recvTimeout)
			if err := this.conn.SetReadDeadline(timeout); err != nil {
				errs <- err
			} else if err := binary.Read(this.conn, binary.BigEndian, &length); err != nil {
				if err == io.EOF || os.IsTimeout(err) {
					// Ignore error
				} else {
					errs <- err
				}
			} else if err := this.decode(length); err != nil {
				errs <- err
			}
		}
	}
}

func (this *connection) decode(length uint32) error {
	payload := make([]byte, length)

	// Ignore zero-sized data
	if length == 0 {
		return nil
	}

	// Receive, decode and send any follow-ups
	if size, err := io.ReadFull(this.conn, payload); err != nil {
		return err
	} else if uint32(size) != length {
		return fmt.Errorf("Received different number of bytes %v read, expected %v", size, length)
	} else if data, err := this.channel.decode(payload); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}
