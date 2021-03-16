package input

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.InputManager
	graph.RegisterUnit(reflect.TypeOf(&Manager{}), reflect.TypeOf((*gopi.InputManager)(nil)))
}
