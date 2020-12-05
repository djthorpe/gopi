package mdns

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register *mdns.Discovery -> gopi.ServiceDiscovery
	graph.RegisterUnit(reflect.TypeOf(&Discovery{}), reflect.TypeOf((*gopi.ServiceDiscovery)(nil)))
}
