package lirc

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
	"github.com/djthorpe/gopi/v3/pkg/hw/lirc/keycode"
)

func init() {
	graph.RegisterUnit(reflect.TypeOf(&lirc{}), reflect.TypeOf((*gopi.LIRC)(nil)))
	graph.RegisterUnit(reflect.TypeOf(&keycode.Manager{}), reflect.TypeOf((*gopi.LIRCKeycodeManager)(nil)))
}
