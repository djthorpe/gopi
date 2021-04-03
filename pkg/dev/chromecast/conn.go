package chromecast

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

type Conn struct {
	sync.Mutex
	sync.WaitGroup
	Channel
	*tls.Conn

	cancel context.CancelFunc
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewConnWithTimeout(key string, addr string, timeout time.Duration) (*Conn, error) {
	this := new(Conn)

	// Connect
	if conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
	}, "tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	}); err != nil {
		return nil, err
	} else {
		this.Conn = conn
		this.Channel.Init(key)
	}

	// Start the receive loop, which will end on cancel()
	ctx, cancel := context.WithCancel(context.Background())
	this.cancel = cancel
	this.WaitGroup.Add(1)
	go func(ctx context.Context) {
		defer this.WaitGroup.Done()
		this.recv(ctx, timeout)
	}(ctx)

	return this, nil
}

func (this *Conn) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// End receive loop and wait
	if this.cancel != nil {
		this.cancel()
		this.WaitGroup.Wait()
		this.cancel = nil
	}

	// Close connection
	var result error
	if this.Conn != nil {
		if err := this.Conn.Close(); err != nil {
			result = multierror.Append(result, err)
		}
		this.Conn = nil
	}

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Conn) String() string {
	str := "<cast.conn"
	str += fmt.Sprintf(" key=%q", this.key)
	if addr := this.Addr(); addr != nil {
		str += fmt.Sprint(" addr=", addr)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Conn) Addr() net.Addr {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if this.Conn == nil {
		return nil
	} else {
		return this.Conn.RemoteAddr()
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Conn) send(data []byte) error {
	if len(data) == 0 {
		return nil
	} else if this.Addr() == nil {
		return gopi.ErrOutOfOrder
	} else if err := binary.Write(this.Conn, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	} else if _, err := this.Conn.Write(data); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Conn) recv(ctx context.Context, timeout time.Duration) {
	var length uint32
	for {
		select {
		case <-ctx.Done():
			return
		default:
			timeout := time.Now().Add(timeout)
			if this.Addr() == nil {
				// Do not read
			} else if err := this.Conn.SetReadDeadline(timeout); err != nil {
				fmt.Println("recv ERROR", err)
			} else if err := binary.Read(this.Conn, binary.BigEndian, &length); err != nil {
				if err == io.EOF || os.IsTimeout(err) {
					// Ignore error
				} else {
					fmt.Println("recv ERROR", err)
				}
			} else if err := this.recvdata(length); err != nil {
				fmt.Println("recv ERROR", err)
			}
		}
	}
}

func (this *Conn) recvdata(length uint32) error {
	payload := make([]byte, length)

	// Ignore zero-sized data
	if length == 0 {
		return nil
	}

	// Receive, decode and send any follow-ups
	if size, err := io.ReadFull(this.Conn, payload); err != nil {
		return err
	} else if uint32(size) != length {
		return fmt.Errorf("Received different number of bytes %v read, expected %v", size, length)
	} else if response, err := this.Channel.decode(payload); err != nil {
		return err
	} else if err := this.send(response); err != nil {
		return err
	}

	// Return success
	return nil
}
