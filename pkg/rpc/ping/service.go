package ping

import (
	"context"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type service struct {
	gopi.Logger
	gopi.Unit
	gopi.Server
	sync.Mutex
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *service) New(cfg gopi.Config) error {
	return this.Server.RegisterService(RegisterPingServer, this)
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *service) CancelStreams() {}

/////////////////////////////////////////////////////////////////////
// RPC METHODS

// Ping simply returns an empty data structure when called
func (this *service) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.Logger.Debug("<Ping>")
	return &empty.Empty{}, nil
}

// Version returns information about the running process
func (this *service) Version(context.Context, *empty.Empty) (*VersionResponse, error) {
	this.Logger.Debug("<Version>")
	return &VersionResponse{}, nil
}
