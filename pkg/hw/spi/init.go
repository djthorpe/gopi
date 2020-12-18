package spi

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register spi.SPI as gopi.SPI
	graph.RegisterUnit(reflect.TypeOf(&spi{}), reflect.TypeOf((*gopi.SPI)(nil)))
}
