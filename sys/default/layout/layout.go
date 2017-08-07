package layout /* import "github.com/djthorpe/gopi/sys/default/layout" */

import (
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register layout module
	registerFlags(gopi.RegisterModule(gopi.Module{
		Name: "default/layout",
		Type: gopi.MODULE_TYPE_LAYOUT,
		New:  newLayout,
	}))
}

func registerFlags(flags *util.Flags) {
	/* no flags for the layout module */
}

func newLayout(config *gopi.AppConfig) (gopi.Driver, error) {
	layout_factory, ok := gopi.Open2(Config{})
}
