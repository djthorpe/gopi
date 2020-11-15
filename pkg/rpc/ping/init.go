package ping

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.PingService and gopi.PingStub
	graph.RegisterUnit(reflect.TypeOf(&service{}), reflect.TypeOf((*gopi.PingService)(nil)))
	graph.RegisterServiceStub(_Ping_serviceDesc.ServiceName, reflect.TypeOf(&stub{}))
}
