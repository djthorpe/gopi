package layout /* import "github.com/djthorpe/gopi/sys/default/layout" */

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/third_party/flex"
	"github.com/djthorpe/gopi/util"
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
	if root := this.newView(tag, class); root == nil {
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
// PUBLIC METHODS FOR VIEW

// Return view tag or zero if not defined
func (this *view) Tag() uint {
	return this.tag
}

// Return view class
func (this *view) Class() string {
	return this.class
}

func (this *view) Position() gopi.ViewPosition {
	switch this.node.Style.PositionType {
	case flex.PositionTypeRelative:
		return gopi.VIEW_POSITION_RELATIVE
	case flex.PositionTypeAbsolute:
		return gopi.VIEW_POSITION_ABSOLUTE
	default:
		panic("Invalid ViewPosition value")
	}
}

func (this *view) Direction() gopi.ViewDirection {
	switch this.node.Style.FlexDirection {
	case flex.FlexDirectionColumn:
		return gopi.VIEW_DIRECTION_COLUMN
	case flex.FlexDirectionColumnReverse:
		return gopi.VIEW_DIRECTION_COLUMN_REVERSE
	case flex.FlexDirectionRow:
		return gopi.VIEW_DIRECTION_ROW
	case flex.FlexDirectionRowReverse:
		return gopi.VIEW_DIRECTION_ROW_REVERSE
	default:
		panic("Invalid ViewDirection value")
	}
}

func (this *view) Justify() gopi.ViewJustify {
	switch this.node.Style.JustifyContent {
	case flex.JustifyCenter:
		return gopi.VIEW_JUSTIFY_CENTER
	case flex.JustifyFlexEnd:
		return gopi.VIEW_JUSTIFY_END
	case flex.JustifyFlexStart:
		return gopi.VIEW_JUSTIFY_START
	case flex.JustifySpaceAround:
		return gopi.VIEW_JUSTIFY_SPACE_AROUND
	case flex.JustifySpaceBetween:
		return gopi.VIEW_JUSTIFY_SPACE_BETWEEN
	default:
		panic("Invalid ViewJustify value")
	}
}

func (this *view) Wrap() gopi.ViewWrap {
	switch this.node.Style.FlexWrap {
	case flex.WrapWrap:
		return gopi.VIEW_WRAP_ON
	case flex.WrapNoWrap:
		return gopi.VIEW_WRAP_OFF
	default:
		panic("Invalid ViewWrap value")
	}
}

func (this *view) Align() gopi.ViewAlign {
	switch this.node.Style.AlignContent {
	case flex.WrapWrap:
		return gopi.VIEW_WRAP_ON
	case flex.WrapNoWrap:
		return gopi.VIEW_WRAP_OFF
	default:
		panic("Invalid ViewWrap value")
	}
}

func (this *view) SetDirection(value gopi.ViewDirection) {

}

func (this *view) SetJustify(value gopi.ViewJustify) {

}

func (this *view) SetWrap(value gopi.ViewWrap) {

}

func (this *view) SetAlign(value gopi.ViewAlign) {

}

func (this *view) SetPositionAbsolute(left, right, top, bottom, start, end float) {
	this.node.StyleSetPositionType(flex.PositionTypeAbsolute)
	this.node.StyleSetPosition(flex.EdgeLeft, left)
	this.node.StyleSetPosition(flex.EdgeRight, right)
	this.node.StyleSetPosition(flex.EdgeTop, top)
	this.node.StyleSetPosition(flex.EdgeBottom, bottom)
	this.node.StyleSetPosition(flex.EdgeStart, start)
	this.node.StyleSetPosition(flex.EdgeEnd, end)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *driver) newView(tag uint, class string) *view {
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

	// Set default style
	v.node.Style.PositionType = flex.PositionTypeRelative

	// TODO: Copy style across to node

	// Return view
	return v
}

func (this *driver) isValidViewClass(class string) bool {
	return reViewClassName.MatchString(class)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	return fmt.Sprintf("gopi.sys.default.Layout{ direction=%v }", this.direction)
}
