package ping

import (
	"fmt"
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register metrics.PingService and metrics.PingStub
	graph.RegisterUnit(reflect.TypeOf(&service{}), reflect.TypeOf((*gopi.PingService)(nil)))
	fmt.Println("TODO REGISTER", _Ping_serviceDesc.ServiceName, reflect.TypeOf(&stub{}))
}
