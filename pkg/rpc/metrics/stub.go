package metrics

import (
	"context"
	"io"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type stub struct {
	gopi.Conn
	MetricsClient
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *stub) New(conn gopi.Conn) {
	this.Conn = conn
	this.MetricsClient = NewMetricsClient(conn.(grpc.ClientConnInterface))
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *stub) List(ctx context.Context) ([]gopi.Measurement, error) {
	// Ensure one call per connection
	this.Conn.Lock()
	defer this.Conn.Unlock()

	metrics, err := this.MetricsClient.List(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	results := make([]gopi.Measurement, len(metrics.Metric))
	for i, metric := range metrics.Metric {
		results[i] = fromProtoMeasurement(metric)
	}
	return results, nil
}

func (this *stub) Stream(ctx context.Context, name string, ch chan<- gopi.Measurement) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	stream, err := this.MetricsClient.Stream(ctx, &Name{Name: name})
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
			} else if evt := fromProtoMeasurement(msg); evt != nil && evt.Name() != "" {
				ch <- evt
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stub) String() string {
	str := "<rpc.stub.metrics"
	str += " addr=" + strconv.Quote(this.Addr())
	return str + ">"
}
