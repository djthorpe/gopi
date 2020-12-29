package metrics

import (
	"context"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type service struct {
	gopi.Logger
	sync.Mutex
	gopi.Publisher
	gopi.Server
	gopi.Unit
	gopi.Metrics
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *service) New(cfg gopi.Config) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("RegisterService: ", "(Server == nil)")
	} else if this.Metrics == nil {
		return gopi.ErrInternalAppError.WithPrefix("RegisterService: ", "(Metrics == nil)")
	} else {
		return this.Server.RegisterService(RegisterMetricsServer, this)
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *service) mustEmbedUnimplementedMetricsServer() {}

/////////////////////////////////////////////////////////////////////
// RPC METHODS

// List registered measurements
func (this *service) List(context.Context, *empty.Empty) (*Measurements, error) {
	this.Logger.Debug("<List>")

	measurements := this.Metrics.Measurements()
	response := &Measurements{
		Metric: make([]*Measurement, 0, len(measurements)),
	}
	for _, measurement := range measurements {
		if pb := toProtoMeasurement(measurement); pb != nil {
			response.Metric = append(response.Metric, pb)
		}
	}

	return response, nil
}

// Stream measuremets to client
func (this *service) Stream(req *Name, stream Metrics_StreamServer) error {
	this.Logger.Debug("<Stream", req, ">")

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
			if measurement, ok := evt.(gopi.Measurement); ok {
				// TODO FILTER MEAUSREMENTS
				if err := stream.Send(toProtoMeasurement(measurement)); err != nil {
					this.Print(err)
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
