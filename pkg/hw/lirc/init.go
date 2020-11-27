package lirc

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	graph.RegisterUnit(reflect.TypeOf(&lirc{}), reflect.TypeOf((*gopi.LIRC)(nil)))
}
