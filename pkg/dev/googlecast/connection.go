package googlecast

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type connection struct {
	sync.RWMutex
	sync.WaitGroup
	channel

	conn   *tls.Conn
	lock   sync.RWMutex // Additional lock just for conn value
	cancel context.CancelFunc
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	recvTimeout = 500 * time.Millisecond
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *connection) Connect(key, addr string, timeout time.Duration, state chan<- state) error {
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
	ctx, cancel := context.WithCancel(context.Background())
	this.cancel = cancel
	this.WaitGroup.Add(1)
	go func(ctx context.Context) {
		defer this.WaitGroup.Done()
		this.recv(ctx, state)
	}(ctx)

	// Send CONNECT message. Errors here require disconnect
	if _, data, err := this.channel.Connect(); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	} else {
		this.channel.ping = time.Time{}
	}

	// Success
	return nil
}

func (this *connection) Disconnect() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// If already disconnected, return
	if this.isConnected() == false {
		return nil
	}

	// Send CLOSE message, reset ping time
	var result error
	if _, data, err := this.channel.Disconnect(); err != nil {
		result = multierror.Append(result, err)
	} else if err := this.send(data); err != nil {
		result = multierror.Append(result, err)
	} else {
		this.channel.ping = time.Time{}
	}

	// Close connection
	if err := this.conn.Close(); err != nil {
		result = multierror.Append(result, err)
	} else {
		this.conn = nil
	}

	// End receive loop and wait
	this.cancel()
	this.WaitGroup.Wait()

	// Release resources
	this.cancel = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *connection) Addr() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if this.isConnected() {
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
	} else if this.isConnected() == false {
		return gopi.ErrOutOfOrder
	} else if err := binary.Write(this.conn, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	} else if _, err := this.conn.Write(data); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *connection) recv(ctx context.Context, state chan<- state) {
	var length uint32
	for {
		select {
		case <-ctx.Done():
			return
		default:
			timeout := time.Now().Add(recvTimeout)
			if this.isConnected() == false {
				// Do not read
			} else if err := this.conn.SetReadDeadline(timeout); err != nil {
				this.recverror(err)
			} else if err := binary.Read(this.conn, binary.BigEndian, &length); err != nil {
				if err == io.EOF || os.IsTimeout(err) {
					// Ignore error
				} else {
					this.recverror(err)
				}
			} else if err := this.decode(length); err != nil {
				this.recverror(err)
			}
		}
	}
}

func (this *connection) recverror(err error) {
	if strings.HasSuffix(err.Error(), "use of closed network connection") {
		// HACK! Ignore error
		return
	} else if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNABORTED) || errors.Is(err, syscall.ENETUNREACH) {
		// Broken pipe, disconnect
		this.ch <- Close(this.channel.key)
		return
	}

	// Report other errors in the background, unless channel is full
	go func() {
		select {
		case this.ch <- NewError(this.channel.key, err):
			break
		default:
			fmt.Println("[DEADLOCK]", err)
		}
	}()
}

func (this *connection) isConnected() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.conn != nil
}

func (this *connection) decode(length uint32) error {
	payload := make([]byte, length)

	// Ignore zero-sized data
	if length == 0 {
		return nil
	}

	// Receive, decode and send any follow-ups
	if this.isConnected() == false {
		// Ignore when no connection
	} else if size, err := io.ReadFull(this.conn, payload); err != nil {
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
