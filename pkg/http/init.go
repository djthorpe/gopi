package http

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register as gopi.Server
	graph.RegisterUnit(reflect.TypeOf(&Server{}), reflect.TypeOf((*gopi.Server)(nil)))
}
