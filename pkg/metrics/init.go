package metrics

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/event"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	graph.RegisterUnit(reflect.TypeOf(&metrics{}), reflect.TypeOf((*gopi.Metrics)(nil)))
}
