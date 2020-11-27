package broadcom

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register broadcom.GPIO as gopi.GPIO
	graph.RegisterUnit(reflect.TypeOf(&GPIO{}), reflect.TypeOf((*gopi.GPIO)(nil)))
}
