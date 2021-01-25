// +build dvb

package dvb

import (
	"reflect"

	graph "github.com/djthorpe/gopi/v3/pkg/graph"
	gopi "github.com/djthorpe/gopi/v3"
)

func init() {
	graph.RegisterUnit(reflect.TypeOf(&Manager{}), reflect.TypeOf((*gopi.DVBManager)(nil)))
}
