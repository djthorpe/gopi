package client

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register metrics.Server
	graph.RegisterUnit(reflect.TypeOf(&connpool{}), reflect.TypeOf((*gopi.ConnPool)(nil)))
}
