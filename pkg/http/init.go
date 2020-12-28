package http

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register server and services
	graph.RegisterUnit(reflect.TypeOf(&Server{}), reflect.TypeOf((*gopi.Server)(nil)))
	graph.RegisterUnit(reflect.TypeOf(&Static{}), reflect.TypeOf((*gopi.HttpStatic)(nil)))
}
