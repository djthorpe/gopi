package rotel

import (
	"context"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type service struct {
	sync.Mutex
	gopi.Unit
	gopi.Logger
	gopi.Server
	gopi.RotelManager
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *service) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.Server, this.RotelManager)
	return this.Server.RegisterService(RegisterManagerServer, this)
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *service) CancelStreams() {}

func (this *service) mustEmbedUnimplementedManagerServer() {}

/////////////////////////////////////////////////////////////////////
// RPC METHODS

// Change Power State
func (this *service) SetPower(_ context.Context, req *Bool) (*State, error) {
	this.Logger.Debug("<SetPower ", req, ">")

	if err := this.RotelManager.SetPower(req.Value); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

// Change Input Source
func (this *service) SetSource(_ context.Context, req *String) (*State, error) {
	this.Logger.Debug("<SetSource ", req, ">")

	if err := this.RotelManager.SetSource(req.Value); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

// Change Volume
func (this *service) SetVolume(_ context.Context, req *Uint) (*State, error) {
	this.Logger.Debug("<SetVolume ", req, ">")

	if err := this.RotelManager.SetVolume(uint(req.Value)); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

func (this *service) SetMute(_ context.Context, req *Bool) (*State, error) {
	this.Logger.Debug("<SetMute ", req, ">")

	if err := this.RotelManager.SetMute(req.Value); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

func (this *service) SetBypass(_ context.Context, req *Bool) (*State, error) {
	this.Logger.Debug("<SetBypass ", req, ">")

	if err := this.RotelManager.SetBypass(req.Value); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

func (this *service) SetTreble(_ context.Context, req *Int) (*State, error) {
	this.Logger.Debug("<SetTreble ", req, ">")

	if err := this.RotelManager.SetTreble(int(req.Value)); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

func (this *service) SetBass(_ context.Context, req *Int) (*State, error) {
	this.Logger.Debug("<SetBass ", req, ">")

	if err := this.RotelManager.SetBass(int(req.Value)); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

func (this *service) SetBalance(_ context.Context, req *String) (*State, error) {
	this.Logger.Debug("<SetBalance ", req, ">")

	if err := this.RotelManager.SetBalance(req.Value); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

func (this *service) SetDimmer(_ context.Context, req *Uint) (*State, error) {
	this.Logger.Debug("<SetDimmer ", req, ">")

	if err := this.RotelManager.SetDimmer(uint(req.Value)); err != nil {
		return nil, err
	} else {
		return toProtoState(this.RotelManager), nil
	}
}

func (this *service) Play(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.Logger.Debug("<Play>")

	if err := this.RotelManager.Play(); err != nil {
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *service) Stop(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.Logger.Debug("<Stop>")

	if err := this.RotelManager.Stop(); err != nil {
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *service) Pause(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.Logger.Debug("<Pause>")

	if err := this.RotelManager.Pause(); err != nil {
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *service) NextTrack(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.Logger.Debug("<NextTrack>")

	if err := this.RotelManager.NextTrack(); err != nil {
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *service) PrevTrack(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.Logger.Debug("<PrevTrack>")

	if err := this.RotelManager.PrevTrack(); err != nil {
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

// Stream Rotel events to client
func (this *service) Stream(_ *empty.Empty, stream Manager_StreamServer) error {
	this.Logger.Debug("<Stream>")

	// Send a null event once a second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Subscribe to input events
	ch := this.RotelManager.Subscribe()
	defer this.RotelManager.Unsubscribe(ch)

	// Obtain server cancel context
	ctx := this.Server.NewStreamContext()

	// Loop which streams until server context cancels or an error occurs sending a Ping
	for {
		select {
		case evt := <-ch:
			// Send event
			if evt_, ok := evt.(gopi.RotelEvent); ok {
				this.Debug("Stream: ", evt_)
				if err := stream.Send(toProtoEvent(evt_)); err != nil {
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
