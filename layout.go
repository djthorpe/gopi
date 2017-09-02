/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"

import (
	"encoding/xml"
	"math"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// LayoutDirection is either left to right, right to left, or inherited
// from parent
type LayoutDirection uint

// ViewDirection is how the view children are laid out, column,
// columnereverse, row, rowreverse
type ViewDirection uint

// ViewDisplay is either flex (default) or none
type ViewDisplay uint

// ViewOverflow is either visible, hidden or scroll
type ViewOverflow uint

// ViewJustify is how view children are aligned within a parent
// START is default
type ViewJustify uint

// ViewWrap is how view children move to the next line within a
// parent (default off)
type ViewWrap uint

// ViewAlign determines alignment of children within the parent
// by default it's STRETCH
type ViewAlign uint

// ViewPosition is by default relative. When absolute, uses
// only left, right, top, bottom, start and end in order to
// set position
type ViewPositioning uint

// ViewEdge defines an edge
type ViewEdge uint

// ViewDimension defines width or height
type ViewDimension uint

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	LAYOUT_DIRECTION_NONE LayoutDirection = iota
	LAYOUT_DIRECTION_LEFTRIGHT
	LAYOUT_DIRECTION_RIGHTLEFT
	LAYOUT_DIRECTION_INHERIT = LAYOUT_DIRECTION_NONE
)

const (
	VIEW_DIRECTION_ROW ViewDirection = iota
	VIEW_DIRECTION_COLUMN
	VIEW_DIRECTION_ROW_REVERSE
	VIEW_DIRECTION_COLUMN_REVERSE
)

const (
	VIEW_DISPLAY_FLEX ViewDisplay = iota
	VIEW_DISPLAY_NONE
)

const (
	VIEW_OVERFLOW_VISIBLE ViewOverflow = iota
	VIEW_OVERFLOW_HIDDEN
	VIEW_OVERFLOW_SCROLL
)

const (
	VIEW_JUSTIFY_FLEX_START ViewJustify = iota
	VIEW_JUSTIFY_FLEX_END
	VIEW_JUSTIFY_CENTER
	VIEW_JUSTIFY_SPACE_BETWEEN
	VIEW_JUSTIFY_SPACE_AROUND
	VIEW_JUSTIFY_CENTRE = VIEW_JUSTIFY_CENTER
)

const (
	VIEW_WRAP_ON ViewWrap = iota
	VIEW_WRAP_OFF
	VIEW_WRAP_REVERSE
)

const (
	VIEW_ALIGN_AUTO ViewAlign = iota
	VIEW_ALIGN_FLEX_START
	VIEW_ALIGN_CENTER
	VIEW_ALIGN_FLEX_END
	VIEW_ALIGN_STRETCH
	VIEW_ALIGN_BASELINE
	VIEW_ALIGN_SPACE_BETWEEN
	VIEW_ALIGN_SPACE_AROUND
	VIEW_ALIGN_CENTRE = VIEW_ALIGN_CENTER
)

const (
	VIEW_POSITIONING_RELATIVE ViewPositioning = iota
	VIEW_POSITIONING_ABSOLUTE
)

const (
	VIEW_EDGE_NONE ViewEdge = iota
	VIEW_EDGE_TOP
	VIEW_EDGE_BOTTOM
	VIEW_EDGE_LEFT
	VIEW_EDGE_RIGHT
	VIEW_EDGE_ALL
)

const (
	VIEW_DIMENSION_NONE ViewDimension = iota
	VIEW_DIMENSION_WIDTH
	VIEW_DIMENSION_HEIGHT
	VIEW_DIMENSION_ALL
)

/////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// EdgeUndefined is used to set a position as "not defined" or "auto"
	EdgeUndefined float32 = float32(math.NaN())

	// ValueAuto is used to set a value to "auto"
	ValueAuto float32 = float32(math.NaN())

	// BasisAuto is the 'auto' setting for basis
	BasisAuto float32 = float32(math.NaN())

	// TagNone is when there is no tag associated with a view
	TagNone uint = 0
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

// View defines a 2D rectangular drawable element
type View interface {
	// Return tag for this view
	Tag() uint

	// Return class for this view
	Class() string

	// Return positioning
	Positioning() ViewPositioning

	// Get Style Attributes
	Display() ViewDisplay
	Overflow() ViewOverflow
	Direction() ViewDirection
	Wrap() ViewWrap
	JustifyContent() ViewJustify
	AlignItems() ViewAlign
	AlignContent() ViewAlign
	AlignSelf() ViewAlign
	Grow() float32
	Shrink() float32
	BasisString() string

	// Set Style Attributes
	SetDisplay(value ViewDisplay)
	SetOverflow(value ViewOverflow)
	SetDirection(value ViewDirection)
	SetWrap(value ViewWrap)
	SetJustifyContent(value ViewJustify)
	SetAlignItems(value ViewAlign)
	SetAlignContent(value ViewAlign)
	SetAlignSelf(value ViewAlign)
	SetGrow(value float32)
	SetShrink(value float32)
	SetBasisValue(value float32)
	SetBasisPercent(value float32)
	SetBasisAuto()

	// Set position
	SetPositionValue(value float32, edges ...ViewEdge)
	SetPositionPercent(percent float32, edges ...ViewEdge)

	// Set padding
	SetPaddingValue(value float32, edges ...ViewEdge)
	SetPaddingPercent(percent float32, edges ...ViewEdge)

	// Set margin
	SetMarginValue(value float32, edges ...ViewEdge)
	SetMarginPercent(percent float32, edges ...ViewEdge)
	SetMarginAuto(edges ...ViewEdge)

	// Set width and height
	SetDimensionValue(value float32, dimension ViewDimension)
	SetDimensionPercent(percent float32, dimension ViewDimension)
	SetDimensionAuto(dimension ViewDimension)

	// Minimum and maximum dimensions
	SetDimensionMinValue(value float32, dimension ViewDimension)
	SetDimensionMinPercent(percent float32, dimension ViewDimension)
	SetDimensionMinAuto(dimension ViewDimension)
	SetDimensionMaxValue(value float32, dimension ViewDimension)
	SetDimensionMaxPercent(percent float32, dimension ViewDimension)
	SetDimensionMaxAuto(dimension ViewDimension)

	// Get strings for position, margin and padding, each edge is separated by a space
	PositionString(edges ...ViewEdge) string
	MarginString(edges ...ViewEdge) string
	PaddingString(edges ...ViewEdge) string
	DimensionString(dimension ViewDimension) string
	DimensionMinString(dimension ViewDimension) string
	DimensionMaxString(dimension ViewDimension) string

	// Determine if view changes on this element require layout
	IsDirty() bool

	// Get layout values which are provided once 'CalculateLayout'
	// has been called in the Layout object
	LayoutValue(edge ViewEdge) float32
	LayoutWidth() float32
	LayoutHeight() float32

	// Require ability to marshall XML
	xml.Marshaler
}

// Layout defines the methods of calculating layout of views within
// a rectangular surface (window)
type Layout interface {
	Driver

	// Return default view direction
	Direction() LayoutDirection

	// Create a root view with a particular tag and view class,
	// returns a nil object if the view could not be created
	// due to invalid class or pre-existing tag
	NewRootView(tag uint, class string) View

	// Return the root view for a particular tag, returns
	// a nil object if the view could not be found
	RootViewForTag(tag uint) View

	// Calculate layout using the root node as the correct size
	//CalculateLayoutForTag(tag uint) bool

	// Calculate layout with a new root size
	//CalculateLayoutForTagWithSize(tag uint, w, h float32) bool

	// Return XML encoded version of the layout
	// Encode(w io.Writer, indent EncodeIndentOptions) error
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v ViewPositioning) String() string {
	switch v {
	case VIEW_POSITIONING_RELATIVE:
		return "VIEW_POSITIONING_RELATIVE"
	case VIEW_POSITIONING_ABSOLUTE:
		return "VIEW_POSITIONING_ABSOLUTE"
	default:
		return "[?? Invalid ViewPositioning value]"
	}
}

func (v ViewDisplay) String() string {
	switch v {
	case VIEW_DISPLAY_FLEX:
		return "VIEW_DISPLAY_FLEX"
	case VIEW_DISPLAY_NONE:
		return "VIEW_DISPLAY_NONE"
	default:
		return "[?? Invalid ViewDisplay value]"
	}
}

func (v ViewOverflow) String() string {
	switch v {
	case VIEW_OVERFLOW_VISIBLE:
		return "VIEW_OVERFLOW_VISIBLE"
	case VIEW_OVERFLOW_HIDDEN:
		return "VIEW_OVERFLOW_HIDDEN"
	case VIEW_OVERFLOW_SCROLL:
		return "VIEW_OVERFLOW_SCROLL"
	default:
		return "[?? Invalid ViewOverflow value]"
	}
}

func (d LayoutDirection) String() string {
	switch d {
	case LAYOUT_DIRECTION_INHERIT:
		return "LAYOUT_DIRECTION_INHERIT"
	case LAYOUT_DIRECTION_LEFTRIGHT:
		return "LAYOUT_DIRECTION_LEFTRIGHT"
	case LAYOUT_DIRECTION_RIGHTLEFT:
		return "LAYOUT_DIRECTION_RIGHTLEFT"
	default:
		return "[?? Invalid LayoutDirection value]"
	}
}

func (d ViewDirection) String() string {
	switch d {
	case VIEW_DIRECTION_COLUMN:
		return "VIEW_DIRECTION_COLUMN"
	case VIEW_DIRECTION_COLUMN_REVERSE:
		return "VIEW_DIRECTION_COLUMN_REVERSE"
	case VIEW_DIRECTION_ROW:
		return "VIEW_DIRECTION_ROW"
	case VIEW_DIRECTION_ROW_REVERSE:
		return "VIEW_DIRECTION_ROW_REVERSE"
	default:
		return "[?? Invalid ViewDirection value]"
	}
}

func (v ViewJustify) String() string {
	switch v {
	case VIEW_JUSTIFY_FLEX_START:
		return "VIEW_JUSTIFY_FLEX_START"
	case VIEW_JUSTIFY_FLEX_END:
		return "VIEW_JUSTIFY_FLEX_END"
	case VIEW_JUSTIFY_CENTER:
		return "VIEW_JUSTIFY_CENTER"
	case VIEW_JUSTIFY_SPACE_BETWEEN:
		return "VIEW_JUSTIFY_SPACE_BETWEEN"
	case VIEW_JUSTIFY_SPACE_AROUND:
		return "VIEW_JUSTIFY_SPACE_AROUND"
	default:
		return "[?? Invalid ViewJustify value]"
	}
}

func (v ViewWrap) String() string {
	switch v {
	case VIEW_WRAP_ON:
		return "VIEW_WRAP_ON"
	case VIEW_WRAP_OFF:
		return "VIEW_WRAP_OFF"
	case VIEW_WRAP_REVERSE:
		return "VIEW_WRAP_REVERSE"
	default:
		return "[?? Invalid ViewWrap value]"
	}
}

func (v ViewAlign) String() string {
	switch v {
	case VIEW_ALIGN_AUTO:
		return "VIEW_ALIGN_AUTO"
	case VIEW_ALIGN_FLEX_START:
		return "VIEW_ALIGN_FLEX_START"
	case VIEW_ALIGN_CENTER:
		return "VIEW_ALIGN_CENTER"
	case VIEW_ALIGN_FLEX_END:
		return "VIEW_ALIGN_FLEX_END"
	case VIEW_ALIGN_STRETCH:
		return "VIEW_ALIGN_STRETCH"
	case VIEW_ALIGN_BASELINE:
		return "VIEW_ALIGN_BASELINE"
	case VIEW_ALIGN_SPACE_AROUND:
		return "VIEW_ALIGN_SPACE_AROUND"
	case VIEW_ALIGN_SPACE_BETWEEN:
		return "VIEW_ALIGN_SPACE_BETWEEN"
	default:
		return "[?? Invalid ViewAlign value]"
	}
}

func (v ViewEdge) String() string {
	switch v {
	case VIEW_EDGE_NONE:
		return "VIEW_EDGE_NONE"
	case VIEW_EDGE_TOP:
		return "VIEW_EDGE_TOP"
	case VIEW_EDGE_BOTTOM:
		return "VIEW_EDGE_BOTTOM"
	case VIEW_EDGE_LEFT:
		return "VIEW_EDGE_LEFT"
	case VIEW_EDGE_RIGHT:
		return "VIEW_EDGE_RIGHT"
	case VIEW_EDGE_ALL:
		return "VIEW_EDGE_ALL"
	default:
		return "[?? Invalid ViewEdge value]"
	}
}

func (v ViewDimension) String() string {
	switch v {
	case VIEW_DIMENSION_NONE:
		return "VIEW_DIMENSION_NONE"
	case VIEW_DIMENSION_WIDTH:
		return "VIEW_DIMENSION_WIDTH"
	case VIEW_DIMENSION_HEIGHT:
		return "VIEW_DIMENSION_HEIGHT"
	case VIEW_DIMENSION_ALL:
		return "VIEW_DIMENSION_ALL"
	default:
		return "[?? Invalid ViewDimension value]"
	}
}
