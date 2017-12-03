/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package layout /* import "github.com/djthorpe/gopi/sys/layout" */

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/third_party/flex"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config defines a default configuration for a layout
type Config struct {
	Direction gopi.LayoutDirection
	Width     float32
	Height    float32
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

var (
	// class names are similar to identifiers, must start with alpha followed by alphanumeric (plus - and _)
	reViewClassName = regexp.MustCompile("^[A-Za-z][A-Za-z0-9\\-\\_]*$")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register layout module
	gopi.RegisterModule(gopi.Module{
		Name: "layout/flex",
		Type: gopi.MODULE_TYPE_LAYOUT,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Config{}, app.Logger)
		},
	})
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
// PUBLIC METHODS FOR DRIVER

// Direction returns default direction for new root views
func (this *driver) Direction() gopi.LayoutDirection {
	return this.direction
}

// NewRootView creates a new root view (which is usually the backing
// element of a drawable surface) with a view class and unique tag
func (this *driver) NewRootView(tag uint, class string) gopi.View {
	// Check for pre-existing view with this tag
	if _, exists := this.root[tag]; exists {
		this.log.Error("View tag %v: already exists", tag)
		return nil
	}

	// Ensure tag is non-zero for root views
	if tag == 0 {
		this.log.Error("View tag %v: cannot be zero", tag)
		return nil
	}

	// create the new view
	if root := this.newView(tag, class, gopi.VIEW_POSITIONING_ABSOLUTE); root == nil {
		this.log.Error("View tag %v: Unable to create a new root view", tag)
		return nil
	} else {
		// Append the view
		this.root[tag] = root

		// return the view
		return root
	}
}

// Returns root view for a particular tag or returns nil
func (this *driver) RootViewForTag(tag uint) gopi.View {
	root := this.root[tag]
	return root
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *driver) newView(tag uint, class string, positioning gopi.ViewPositioning) *view {
	// Check for valid view class
	if this.isValidViewClass(class) == false {
		return nil
	}

	// Create a configuration
	flexconfig := flex.NewConfig()
	flexconfig.Logger = func(config *flex.Config, node *flex.Node, level flex.LogLevel, format string, args ...interface{}) int {
		switch level {
		case flex.LogLevelError:
			this.log.Error(format, args)
		case flex.LogLevelFatal:
			this.log.Fatal(format, args)
		case flex.LogLevelWarn:
			this.log.Warn(format, args)
		case flex.LogLevelInfo:
			this.log.Info(format, args)
		case flex.LogLevelDebug:
			this.log.Debug(format, args)
		case flex.LogLevelVerbose:
			this.log.Debug2(format, args)
		}
		return 0
	}

	// Create view with tag and class
	v := new(view)
	v.tag = tag
	v.class = class
	v.node = flex.NewNodeWithConfig(flexconfig)

	// Set positioning with all edges as auto
	switch positioning {
	case gopi.VIEW_POSITIONING_ABSOLUTE:
		v.node.StyleSetPositionType(flex.PositionTypeAbsolute)
	case gopi.VIEW_POSITIONING_RELATIVE:
		v.node.StyleSetPositionType(flex.PositionTypeRelative)
	default:
		return nil
	}
	v.node.StyleSetPosition(flex.EdgeAll, gopi.EdgeUndefined)

	// Set default attributes
	v.SetDisplay(gopi.VIEW_DISPLAY_FLEX)
	v.SetOverflow(gopi.VIEW_OVERFLOW_VISIBLE)
	v.SetDirection(gopi.VIEW_DIRECTION_ROW)
	v.SetWrap(gopi.VIEW_WRAP_OFF)
	v.SetJustifyContent(gopi.VIEW_JUSTIFY_FLEX_START)
	v.SetAlignItems(gopi.VIEW_ALIGN_STRETCH)
	v.SetAlignContent(gopi.VIEW_ALIGN_STRETCH)
	v.SetAlignSelf(gopi.VIEW_ALIGN_AUTO)
	v.SetGrow(0.0)
	v.SetShrink(1.0)
	v.SetDimensionAuto(gopi.VIEW_DIMENSION_ALL)
	v.SetDimensionMinAuto(gopi.VIEW_DIMENSION_ALL)
	v.SetDimensionMaxAuto(gopi.VIEW_DIMENSION_ALL)

	// Return view
	return v
}

func (this *driver) isValidViewClass(class string) bool {
	return reViewClassName.MatchString(class)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	return fmt.Sprintf("gopi.sys.default.Layout{ direction=%v views=%v }", this.direction, this.root)
}
