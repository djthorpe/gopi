package http

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
	handler "github.com/djthorpe/gopi/v3/pkg/http/handler"
)

func init() {
	// Register server and services
	graph.RegisterUnit(reflect.TypeOf(&Server{}), reflect.TypeOf((*gopi.Server)(nil)))
	graph.RegisterUnit(reflect.TypeOf(&handler.Static{}), reflect.TypeOf((*gopi.HttpStatic)(nil)))
	graph.RegisterUnit(reflect.TypeOf(&handler.Logger{}), reflect.TypeOf((*gopi.HttpLogger)(nil)))
	graph.RegisterUnit(reflect.TypeOf(&handler.Templates{}), reflect.TypeOf((*gopi.HttpTemplate)(nil)))
}
