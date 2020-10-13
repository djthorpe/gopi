package freetype

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register FontManager as gopi.FontManager
	graph.RegisterUnit(reflect.TypeOf(&FontManager{}), reflect.TypeOf((*gopi.FontManager)(nil)))
}
