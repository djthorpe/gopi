package chromecast

import (
	context "context"
	"io"
	"strconv"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Stub struct {
	gopi.Conn
	ManagerClient
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *Stub) New(conn gopi.Conn) {
	this.Conn = conn
	this.ManagerClient = NewManagerClient(conn.(grpc.ClientConnInterface))
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Stub) List(ctx context.Context) ([]gopi.Cast, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Stub) Stream(ctx context.Context, ch chan<- gopi.CastEvent) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	stream, err := this.ManagerClient.Stream(ctx, &empty.Empty{})
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
			} else if evt := fromProtoEvent(msg); evt != nil && evt.Flags() != gopi.CAST_FLAG_NONE {
				ch <- evt
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Stub) String() string {
	str := "<chromecast.rpc.stub"
	str += " addr=" + strconv.Quote(this.Addr())
	return str + ">"
}
