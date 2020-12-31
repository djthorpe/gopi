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

func (this *service) SetApp(ctx context.Context, req *AppRequest) (*CastResponse, error) {
	this.Debug("<SetApp ", req, ">")

	// Retrieve Chromecast, LaunchApp
	cast := this.getCastEx(ctx, req.Id)
	if cast == nil {
		return nil, gopi.ErrNotFound.WithPrefix(req.Id)
	} else if err := this.CastManager.LaunchAppWithId(cast, req.Appid); err != nil {
		return nil, err
	}

	// Return success
	return &CastResponse{
		Cast: toProtoCast(cast),
	}, nil
}

func (this *service) LoadURL(ctx context.Context, req *LoadRequest) (*CastResponse, error) {
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

	// Return success
	return &CastResponse{
		Cast: toProtoCast(cast),
	}, nil
}

func (this *service) SetVolume(context.Context, *VolumeRequest) (*CastResponse, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *service) SetMute(context.Context, *MuteRequest) (*CastResponse, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *service) Stop(context.Context, *CastRequest) (*CastResponse, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *service) Play(context.Context, *CastRequest) (*CastResponse, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *service) Pause(context.Context, *CastRequest) (*CastResponse, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *service) Seek(context.Context, *SeekRequest) (*CastResponse, error) {
	return nil, gopi.ErrNotImplemented
}
