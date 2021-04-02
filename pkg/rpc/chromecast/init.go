package chromecast

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.CastService and gopi.CastStub
	graph.RegisterUnit(reflect.TypeOf(&Service{}), reflect.TypeOf((*gopi.CastService)(nil)))
	graph.RegisterServiceStub(Manager_ServiceDesc.ServiceName, reflect.TypeOf(&Stub{}))
}
