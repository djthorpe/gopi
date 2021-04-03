package chromecast

import (
	context "context"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	"github.com/golang/protobuf/ptypes"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type Service struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger
	gopi.Server
	gopi.CastManager
	gopi.Publisher
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *Service) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.Server, this.CastManager, this.Publisher)

	// Register gRPC service
	if err := this.RegisterService(RegisterManagerServer, this); err != nil {
		return err
	}

	// Do an initial scan for chromecasts
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := this.CastManager.Devices(ctx); err != nil {
		return err
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// RPC METHODS

func (this *Service) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	this.Logger.Debug("<List ", req, ">")

	timeout, err := ptypes.Duration(req.Timeout)
	if err != nil {
		return nil, err
	} else if timeout == 0 {
		timeout = time.Second
	}

	timeoutctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	devices, err := this.CastManager.Devices(timeoutctx)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Cast: toProtoCastList(devices),
	}, nil
}

func (this *Service) Stream(_ *empty.Empty, stream Manager_StreamServer) error {
	this.Logger.Debug("<Stream>")

	// Send a null event once a second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Subscribe to input events
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Obtain server cancel context
	ctx := this.Server.NewStreamContext()

	// Loop which streams until server context cancels or an error occurs sending a Ping
	for {
		select {
		case evt := <-ch:
			// Send event
			if evt, ok := evt.(gopi.CastEvent); ok {
				this.Debug("Stream: ", evt)
				if err := stream.Send(toProtoEvent(evt)); err != nil {
					this.Print("Stream: ", err)
				}
			}
		case <-ctx.Done():
			// Context done
			return ctx.Err()
		case <-ticker.C:
			// Send a ping
			if err := stream.Send(toProtoNull()); err != nil {
				this.Debug("Stream: ", "Error sending null event, ending stream")
				return err
			}
		}
	}
}

func (this *Service) Connect(ctx context.Context, req *CastRequest) (*Cast, error) {
	this.Logger.Debug("<Connect ", req, ">")

	if cast := this.CastManager.Get(req.Key); cast == nil {
		return nil, gopi.ErrNotFound
	} else if err := this.CastManager.Connect(cast); err != nil {
		return nil, err
	} else {
		return toProtoCast(cast), nil
	}
}

func (this *Service) Disconnect(ctx context.Context, req *CastRequest) (*empty.Empty, error) {
	this.Logger.Debug("<Disconnect ", req, ">")

	if cast := this.CastManager.Get(req.Key); cast == nil {
		return nil, gopi.ErrNotFound
	} else if err := this.CastManager.Disconnect(cast); err != nil {
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *Service) SetVolume(ctx context.Context, req *VolumeRequest) (*Cast, error) {
	this.Logger.Debug("<SetVolume ", req, ">")

	if cast := this.CastManager.Get(req.Key); cast == nil {
		return nil, gopi.ErrNotFound
	} else if err := this.CastManager.SetVolume(ctx, cast, req.Volume); err != nil {
		return nil, err
	} else {
		return toProtoCast(cast), nil
	}
}

func (this *Service) SetMuted(ctx context.Context, req *MutedRequest) (*Cast, error) {
	this.Logger.Debug("<SetMuted ", req, ">")

	if cast := this.CastManager.Get(req.Key); cast == nil {
		return nil, gopi.ErrNotFound
	} else if err := this.CastManager.SetMuted(ctx, cast, req.Muted); err != nil {
		return nil, err
	} else {
		return toProtoCast(cast), nil
	}
}

func (this *Service) SetApp(ctx context.Context, req *AppRequest) (*Cast, error) {
	this.Logger.Debug("<SetApp ", req, ">")

	if cast := this.CastManager.Get(req.Key); cast == nil {
		return nil, gopi.ErrNotFound
	} else if err := this.CastManager.LaunchAppWithId(ctx, cast, req.App); err != nil {
		return nil, err
	} else {
		return toProtoCast(cast), nil
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Service) mustEmbedUnimplementedManagerServer() {
	// NOOP
}
