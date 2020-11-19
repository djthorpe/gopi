package i2c

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register i2c
	graph.RegisterUnit(reflect.TypeOf(&i2c{}), reflect.TypeOf((*gopi.I2C)(nil)))
}
