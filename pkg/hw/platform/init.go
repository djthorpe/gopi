package platform

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register Platform as gopi.Platform
	graph.RegisterUnit(reflect.TypeOf(&Platform{}), reflect.TypeOf((*gopi.Platform)(nil)))
}
