package input

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.InoutService and gopi.InputStub
	graph.RegisterUnit(reflect.TypeOf(&service{}), reflect.TypeOf((*gopi.InputService)(nil)))
	graph.RegisterServiceStub(Input_ServiceDesc.ServiceName, reflect.TypeOf(&stub{}))
}
