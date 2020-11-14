package ping

import (
	"context"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type service struct {
	gopi.Unit
	gopi.Server
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
	return &empty.Empty{}, nil
}
