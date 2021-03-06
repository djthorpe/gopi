package rotel

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.RotelService and gopi.PingStub
	graph.RegisterUnit(reflect.TypeOf(&service{}), reflect.TypeOf((*gopi.RotelService)(nil)))
	graph.RegisterServiceStub(Manager_ServiceDesc.ServiceName, reflect.TypeOf(&stub{}))
}
