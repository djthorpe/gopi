package rotel

import (
	"fmt"
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

// Stream rotel events to client
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

	// Loop which streams until server context cancels
	// or an error occurs sending a Ping
	for {
		select {
		case evt := <-ch:
			fmt.Println("TODO: Emit: ", evt)
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
