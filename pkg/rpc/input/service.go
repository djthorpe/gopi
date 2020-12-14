package input

import (
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type service struct {
	gopi.Logger
	gopi.Unit
	gopi.Server
	gopi.Publisher
	sync.Mutex
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *service) New(cfg gopi.Config) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("RegisterService: ", "(Server == nil)")
	} else {
		return this.Server.RegisterService(RegisterInputServer, this)
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *service) mustEmbedUnimplementedInputServer() {}

/////////////////////////////////////////////////////////////////////
// RPC METHODS

// Stream will stream all InputEvents until the stream is closed or
// shutdown is requested
func (this *service) Stream(_ *empty.Empty, stream Input_StreamServer) error {
	this.Logger.Debug("<Stream>")

	// Send a null event once a second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Subscribe to input events
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Obtain server cancel context
	ctx := this.Server.NewStreamContext()

	// Loop which streams until server context cancels
	// or an error occurs sending a Ping
	for {
		select {
		case evt := <-ch:
			if evt_, ok := evt.(gopi.InputEvent); ok {
				if err := stream.Send(protoFromInputEvent(evt_)); err != nil {
					this.Print(err)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := stream.Send(protoFromInputEvent(nil)); err != nil {
				this.Logger.Debug("Error sending null event, ending stream")
				return err
			}
		}
	}
}
