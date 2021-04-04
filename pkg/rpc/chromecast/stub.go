package chromecast

import (
	context "context"
	"io"
	"net/url"
	"strconv"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	"github.com/golang/protobuf/ptypes"
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

func (this *Stub) List(ctx context.Context, timeout time.Duration) ([]gopi.Cast, error) {
	if response, err := this.ManagerClient.List(ctx, &ListRequest{
		Timeout: ptypes.DurationProto(timeout),
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCastList(response.Cast), nil
	}
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

func (this *Stub) Connect(ctx context.Context, key string) (gopi.Cast, error) {
	if cast, err := this.ManagerClient.Connect(ctx, &CastRequest{
		Key: key,
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCast(cast), nil
	}
}

func (this *Stub) Disconnect(ctx context.Context, key string) error {
	if _, err := this.ManagerClient.Disconnect(ctx, &CastRequest{
		Key: key,
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Stub) ConnectMedia(ctx context.Context, key string) (gopi.Cast, error) {
	if cast, err := this.ManagerClient.ConnectMedia(ctx, &CastRequest{
		Key: key,
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCast(cast), nil
	}
}

func (this *Stub) DisconnectMedia(ctx context.Context, key string) (gopi.Cast, error) {
	if cast, err := this.ManagerClient.DisconnectMedia(ctx, &CastRequest{
		Key: key,
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCast(cast), nil
	}
}

func (this *Stub) SetVolume(ctx context.Context, key string, level float32) (gopi.Cast, error) {
	if cast, err := this.ManagerClient.SetVolume(ctx, &VolumeRequest{
		Key:    key,
		Volume: level,
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCast(cast), nil
	}
}

func (this *Stub) SetMuted(ctx context.Context, key string, muted bool) (gopi.Cast, error) {
	if cast, err := this.ManagerClient.SetMuted(ctx, &MutedRequest{
		Key:   key,
		Muted: muted,
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCast(cast), nil
	}
}

func (this *Stub) LaunchAppWithId(ctx context.Context, key, app string) (gopi.Cast, error) {
	if cast, err := this.ManagerClient.SetApp(ctx, &AppRequest{
		Key: key,
		App: app,
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCast(cast), nil
	}
}

func (this *Stub) LoadMedia(ctx context.Context, key string, url *url.URL, autoplay bool) (gopi.Cast, error) {
	if url == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("LoadMedia")
	} else if cast, err := this.ManagerClient.LoadMedia(ctx, &MediaRequest{
		Key:      key,
		Url:      url.String(),
		Autoplay: autoplay,
	}); err != nil {
		return nil, err
	} else {
		return fromProtoCast(cast), nil
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Stub) String() string {
	str := "<chromecast.rpc.stub"
	str += " addr=" + strconv.Quote(this.Addr())
	return str + ">"
}
