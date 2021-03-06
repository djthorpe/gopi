package googlecast

import (
	context "context"
	"net/url"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type service struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger
	gopi.Server
	gopi.CastManager
	gopi.Publisher

	// Map of ID to Chromecast
	casts map[string]gopi.Cast
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *service) New(cfg gopi.Config) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("RegisterService: ", "(Server == nil)")
	} else if this.CastManager == nil {
		return gopi.ErrInternalAppError.WithPrefix("RegisterService: ", "(CastManager == nil)")
	} else if this.Logger == nil {
		return gopi.ErrInternalAppError.WithPrefix("RegisterService: ", "(Logger == nil)")
	} else if this.Publisher == nil {
		return gopi.ErrInternalAppError.WithPrefix("RegisterService: ", "(Publisher == nil)")
	} else if err := this.Server.RegisterService(RegisterManagerServer, this); err != nil {
		return err
	}

	// Set up mapping for chromecasts
	this.casts = make(map[string]gopi.Cast)

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *service) mustEmbedUnimplementedManagerServer() {}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// setCasts updates the list of known Chromecasts
func (this *service) setCasts(casts []gopi.Cast) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	for _, cast := range casts {
		key := cast.Id()
		this.casts[key] = cast
	}
}

// lookupCasts using service discovery, cancel after one second
func (this *service) lookupCasts(ctx context.Context) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// Perform discovery
	if casts, err := this.CastManager.Devices(timeout); err != nil && err != context.DeadlineExceeded {
		return err
	} else {
		this.setCasts(casts)
	}

	// Return success
	return nil
}

// getCast returns a known chromecast or nil
func (this *service) getCast(key string) gopi.Cast {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if cast, exists := this.casts[key]; exists {
		return cast
	} else {
		return nil
	}
}

// listCast returns all known chromecasts
func (this *service) listCasts() []gopi.Cast {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	casts := make([]gopi.Cast, 0, len(this.casts))
	for _, cast := range this.casts {
		casts = append(casts, cast)
	}
	return casts
}

// getCastEx returns a known chromecast or does discovery before returning
func (this *service) getCastEx(ctx context.Context, key string) gopi.Cast {
	// Read first
	if cast := this.getCast(key); cast != nil {
		return cast
	}

	// Discover
	this.lookupCasts(ctx)

	// Return discovered cast
	return this.getCast(key)
}

func (this *service) toProtoCastState(cast gopi.Cast) (*CastState, error) {
	var err error

	// Retrieve state information
	level, muted, err := this.CastManager.Volume(cast)
	if err != nil {
		return nil, err
	}
	app, err := this.CastManager.App(cast)
	if err != nil {
		return nil, err
	}
	// Return cast state
	return &CastState{
		Cast:   toProtoCast(cast),
		Volume: toProtoVolume(level, muted),
		App:    toProtoApp(app),
	}, nil
}

/////////////////////////////////////////////////////////////////////
// RPC METHODS

func (this *service) ListCasts(ctx context.Context, _ *empty.Empty) (*ListResponse, error) {
	this.Debug("<ListCasts>")

	// Perform discovery
	if err := this.lookupCasts(ctx); err != nil {
		return nil, err
	}

	// Construct the reply
	reply := &ListResponse{}
	for _, cast := range this.listCasts() {
		reply.Cast = append(reply.Cast, toProtoCast(cast))
	}

	// Return success
	return reply, nil
}

func (this *service) SetApp(ctx context.Context, req *AppRequest) (*CastState, error) {
	this.Debug("<SetApp ", req, ">")

	// Retrieve Chromecast, LaunchApp
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.LaunchAppWithId(cast, req.Appid); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) LoadURL(ctx context.Context, req *LoadRequest) (*CastState, error) {
	this.Debug("<LoadURL ", req, ">")

	// Retrieve Chromecast, LoadURL
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if url, err := url.Parse(req.Url); err != nil {
		return nil, err
	} else if err := this.CastManager.LoadURL(cast, url, true); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) SetVolume(ctx context.Context, req *VolumeRequest) (*CastState, error) {
	this.Debug("<SetVolume ", req, ">")

	// Retrieve Chromecast, SetVolume
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.SetVolume(cast, req.Volume); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) SetMute(ctx context.Context, req *MuteRequest) (*CastState, error) {
	this.Debug("<SetMute ", req, ">")

	// Retrieve Chromecast, SetMuted
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.SetMuted(cast, req.Muted); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) Get(ctx context.Context, req *CastRequest) (*CastState, error) {
	this.Debug("<Get ", req, ">")

	// Retrieve Chromecast
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) Stop(ctx context.Context, req *CastRequest) (*CastState, error) {
	this.Debug("<Stop ", req, ">")

	// Retrieve Chromecast, Play
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.SetPlay(cast, false); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) Play(ctx context.Context, req *CastRequest) (*CastState, error) {
	this.Debug("<Play ", req, ">")

	// Retrieve Chromecast, Play
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.SetPlay(cast, true); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) Pause(ctx context.Context, req *CastRequest) (*CastState, error) {
	this.Debug("<Pause ", req, ">")

	// Retrieve Chromecast, Play
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.SetPause(cast, true); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) SeekAbs(ctx context.Context, req *SeekRequest) (*CastState, error) {
	this.Debug("<SeekAbs ", req, ">")

	// Retrieve Chromecast, Seek
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.SeekAbs(cast, req.Position.AsDuration()); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

func (this *service) SeekRel(ctx context.Context, req *SeekRequest) (*CastState, error) {
	this.Debug("<SeekRel ", req, ">")

	// Retrieve Chromecast, Seek
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.SeekRel(cast, req.Position.AsDuration()); err != nil {
		return nil, err
	}

	// Return cast state
	return this.toProtoCastState(cast)
}

// Stream measuremets to client
func (this *service) Stream(req *CastRequest, stream Manager_StreamServer) error {
	this.Logger.Debug("<Stream", req, ">")

	// Send a null event once a second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Subscribe to events
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Obtain server cancel context
	ctx := this.Server.NewStreamContext()

	// Loop which streams until server context cancels
	// or an error occurs sending a Ping
	for {
		select {
		case evt := <-ch:
			if castevt, ok := evt.(gopi.CastEvent); ok {
				if req.Id == "" || req.Id == castevt.Cast().Id() {
					if err := stream.Send(toProtoEvent(castevt)); err != nil {
						this.Print(err)
					}
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := stream.Send(toProtoNull()); err != nil {
				this.Logger.Debug("Error sending null event, ending stream")
				return err
			}
		}
	}
}
