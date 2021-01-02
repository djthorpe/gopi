package googlecast

import (
	context "context"
	"io"
	"net/url"
	"strconv"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type stub struct {
	gopi.Conn
	ManagerClient
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *stub) New(conn gopi.Conn) {
	this.Conn = conn
	this.ManagerClient = NewManagerClient(conn.(grpc.ClientConnInterface))
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *stub) ListCasts(ctx context.Context) ([]gopi.Cast, error) {
	response, err := this.ManagerClient.ListCasts(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	// Translate protobuf to gopi.Cast interface
	result := make([]gopi.Cast, len(response.Cast))
	for i, cast := range response.Cast {
		result[i] = fromProtoCast(cast)
	}

	// Return success
	return result, nil
}

func (this *stub) SetApp(ctx context.Context, castId, appId string) error {
	if _, err := this.ManagerClient.SetApp(ctx, &AppRequest{
		Id:    castId,
		Appid: appId,
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) LoadURL(ctx context.Context, castId string, url *url.URL) error {
	if _, err := this.ManagerClient.LoadURL(ctx, &LoadRequest{
		Id:  castId,
		Url: url.String(),
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) SetVolume(ctx context.Context, castId string, value float32) error {
	if _, err := this.ManagerClient.SetVolume(ctx, &VolumeRequest{
		Id:     castId,
		Volume: value,
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) SetMute(ctx context.Context, castId string, value bool) error {
	if _, err := this.ManagerClient.SetMute(ctx, &MuteRequest{
		Id:    castId,
		Muted: value,
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) Stop(ctx context.Context, castId string) error {
	if _, err := this.ManagerClient.Stop(ctx, &CastRequest{
		Id: castId,
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) Play(ctx context.Context, castId string) error {
	if _, err := this.ManagerClient.Play(ctx, &CastRequest{
		Id: castId,
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) Pause(ctx context.Context, castId string) error {
	if _, err := this.ManagerClient.Pause(ctx, &CastRequest{
		Id: castId,
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) SeekAbs(ctx context.Context, castId string, value time.Duration) error {
	if _, err := this.ManagerClient.SeekAbs(ctx, &SeekRequest{
		Id:       castId,
		Position: toProtoDuration(value),
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) SeekRel(ctx context.Context, castId string, value time.Duration) error {
	if _, err := this.ManagerClient.SeekRel(ctx, &SeekRequest{
		Id:       castId,
		Position: toProtoDuration(value),
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *stub) Stream(ctx context.Context, id string, ch chan<- gopi.CastEvent) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	stream, err := this.ManagerClient.Stream(ctx, &CastRequest{Id: id})
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
			} else if evt := fromProtoEvent(msg); evt != nil {
				ch <- evt
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stub) String() string {
	str := "<rpc.stub.castmanager"
	str += " addr=" + strconv.Quote(this.Addr())
	return str + ">"
}
