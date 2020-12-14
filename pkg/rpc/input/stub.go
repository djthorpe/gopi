package input

import (
	"context"
	"io"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type stub struct {
	gopi.Conn
	InputClient
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *stub) New(conn gopi.Conn) {
	this.Conn = conn
	this.InputClient = NewInputClient(conn.(grpc.ClientConnInterface))
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *stub) Stream(ctx context.Context, ch chan<- gopi.InputEvent) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	stream, err := this.InputClient.Stream(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if msg, err := stream.Recv(); err == io.EOF {
				return nil
			} else if err != nil {
				return this.Err(err)
			} else if evt := protoToInputEvent(msg); evt != nil {
				ch <- evt
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stub) String() string {
	str := "<rpc.stub.input"
	str += " addr=" + strconv.Quote(this.Addr())
	return str + ">"
}
