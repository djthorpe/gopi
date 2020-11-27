package display

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register display as gopi.Display
	graph.RegisterUnit(reflect.TypeOf(&display{}), reflect.TypeOf((*gopi.Display)(nil)))
}
