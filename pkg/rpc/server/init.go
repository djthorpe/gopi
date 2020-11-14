package server

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register metrics.Server
	graph.RegisterUnit(reflect.TypeOf(&server{}), reflect.TypeOf((*gopi.Server)(nil)))
}
