package metrics

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.MetricsService and gopi.MetricsStub
	graph.RegisterUnit(reflect.TypeOf(&service{}), reflect.TypeOf((*gopi.MetricsService)(nil)))
	graph.RegisterServiceStub(_Metrics_serviceDesc.ServiceName, reflect.TypeOf(&stub{}))
}
