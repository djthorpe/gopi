package chromaprint

import (
	"reflect"

	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register display as gopi.Display
	graph.RegisterUnit(reflect.TypeOf(&manager{}), reflect.TypeOf((*Manager)(nil)))
}
