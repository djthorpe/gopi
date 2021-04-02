package event

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	graph.RegisterUnit(reflect.TypeOf(&publisher{}), reflect.TypeOf((*gopi.Publisher)(nil)))
	graph.RegisterUnit(reflect.TypeOf(&Promises{}), reflect.TypeOf((*gopi.Promises)(nil)))
}
