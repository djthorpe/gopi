// +build egl,gbm

package surface

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.SurfaceManager
	graph.RegisterUnit(reflect.TypeOf(&Manager{}), reflect.TypeOf((*gopi.SurfaceManager)(nil)))
}
