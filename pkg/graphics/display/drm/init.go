// +build drm

package display

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.DisplayManager
	graph.RegisterUnit(reflect.TypeOf(&Manager{}), reflect.TypeOf((*gopi.DisplayManager)(nil)))
}
