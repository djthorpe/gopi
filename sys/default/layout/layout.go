package layout /* import "github.com/djthorpe/gopi/sys/default/layout" */

import (
	"errors"
	"fmt"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config defines a default configuration for a layout
type Config struct {
	Direction gopi.LayoutDirection
}

type view struct {
	root  bool
	tag   uint
	class string
}

type driver struct {
	log       gopi.Logger
	direction gopi.LayoutDirection
	root      map[uint]*view
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ErrNotImplemented = errors.New("Not Implemented")
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

func newLayout(config *gopi.AppConfig, logger gopi.Logger) (gopi.Driver, error) {
	var err gopi.Error
	if layout, ok := gopi.Open2(Config{}, logger, &err).(gopi.Layout); !ok {
		return nil, err
	} else {
		return layout, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Config) Open(logger gopi.Logger) (gopi.Driver, error) {
	this := new(driver)
	this.log = logger

	this.log.Debug2("gopi.sys.default.layout.Open()")

	// If direction is INHERIT or NONE, then choose LEFTRIGHT
	if config.Direction == gopi.LAYOUT_DIRECTION_NONE {
		this.direction = gopi.LAYOUT_DIRECTION_LEFTRIGHT
	} else {
		this.direction = config.Direction
	}

	// Create a map for tag to root view, with an initial
	// capacity of 1
	this.root = make(map[uint]*view, 1)

	return this, nil
}

func (this *driver) Close() error {
	this.log.Debug2("gopi.sys.default.layout.Close()")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *driver) Direction() gopi.LayoutDirection {
	return this.direction
}

func (this *driver) NewRootViewWithTag(tag uint) gopi.View {
	if _, exists := this.root[tag]; exists {
		this.log.Error("Tag %v already exists", tag)
		return nil
	}

	// create the new view
	root := newViewWithTag(tag)
	this.root[tag] = root

	// return the view
	return root
}

func (this *driver) RootViewForTag(tag uint) gopi.View {
	root := this.root[tag]
	return root
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func newViewWithTag(tag uint) *view {
	v := new(view)
	v.tag = tag
	return v
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	return fmt.Sprintf("gopi.sys.default.Layout{ direction=%v }", this.direction)
}
