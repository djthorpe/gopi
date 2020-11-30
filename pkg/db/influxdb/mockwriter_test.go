package influxdb_test

import (
	"context"
	"fmt"
	"net"

	gopi "github.com/djthorpe/gopi/v3"
	influxdb "github.com/djthorpe/gopi/v3/pkg/db/influxdb"
)

type MockWriter struct {
	gopi.Unit
	gopi.Publisher
	gopi.Logger
	*server
	*influxdb.Writer
}

func (this *MockWriter) New(gopi.Config) error {
	port := UnusedPort()
	addr := fmt.Sprint("localhost:", port)
	if server, err := NewMockServer(this.Logger.T(), addr); err != nil {
		this.Logger.T().Error(err)
		return err
	} else {
		this.server = server
		this.Writer.URL.Host = addr
	}
	// return success
	return nil
}

func (this *MockWriter) Run(ctx context.Context) error {
	ch := this.Subscribe()
	defer this.Unsubscribe(ch)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-ch:
			if measurement, ok := evt.(gopi.Measurement); ok && measurement != nil {
				if err := this.Writer.Write(evt.(gopi.Measurement)); err != nil {
					this.Logger.T().Error(err)
				}
			}
		}
	}
}

func (this *MockWriter) Dispose() error {
	if err := this.server.Close(); err != nil {
		this.Logger.T().Error(err)
		return err
	}

	// Release resources
	this.server = nil

	// Return success
	return nil
}

func UnusedPort() int {
	if addr, err := net.ResolveTCPAddr("tcp", "localhost:0"); err != nil {
		return 0
	} else if listener, err := net.ListenTCP("tcp", addr); err != nil {
		return 0
	} else {
		listener.Close()
		return listener.Addr().(*net.TCPAddr).Port
	}
}
