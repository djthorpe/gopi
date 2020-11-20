// +build linux

package file

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register argonone
	graph.RegisterUnit(reflect.TypeOf(&filepoll{}), reflect.TypeOf((*gopi.FilePoll)(nil)))
}
