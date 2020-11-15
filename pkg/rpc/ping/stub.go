package ping

import (
	"context"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type stub struct {
	gopi.Conn
	PingClient
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *stub) New(conn gopi.Conn) {
	this.Conn = conn
	this.PingClient = NewPingClient(conn.(grpc.ClientConnInterface))
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *stub) Ping(ctx context.Context) error {
	// Ensure one call per connection
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.PingClient.Ping(ctx, &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stub) String() string {
	str := "<pingstub"
	str += " addr=" + strconv.Quote(this.Addr())
	return str + ">"
}
