package influxdb

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// *influxdb.Writer -> gopi.MetricWriter
	graph.RegisterUnit(reflect.TypeOf(&Writer{}), reflect.TypeOf((*gopi.MetricWriter)(nil)))
}
