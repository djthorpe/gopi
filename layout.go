/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"

import (
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
type ViewPosition uint

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	LAYOUT_DIRECTION_NONE LayoutDirection = iota
	LAYOUT_DIRECTION_LEFTRIGHT
	LAYOUT_DIRECTION_RIGHTLEFT
	LAYOUT_DIRECTION_INHERIT = LAYOUT_DIRECTION_NONE
)

const (
	VIEW_DIRECTION_COLUMN ViewDirection = iota
	VIEW_DIRECTION_COLUMN_REVERSE
	VIEW_DIRECTION_ROW
	VIEW_DIRECTION_ROW_REVERSE
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
	VIEW_POSITION_RELATIVE ViewPosition = iota
	VIEW_POSITION_ABSOLUTE
)

/////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// Undefined is used to set a position as "not defined" or "auto"
	Undefined float32 = float32(math.NaN())
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

// View defines a 2D rectangular drawable element
type View interface {
	// Return tag for this view
	Tag() uint

	// Return class for this view
	Class() string

	// Get Style Attributes
	Position() ViewPosition
	Direction() ViewDirection
	Justify() ViewJustify
	Wrap() ViewWrap
	Align() ViewAlign

	// Set Style Attributes
	SetDirection(value ViewDirection)
	SetJustify(value ViewJustify)
	SetWrap(value ViewWrap)
	SetAlign(value ViewAlign)

	// Set Absolute positioning
	SetPositionAbsolute()
	SetPositionPixel(ViewEdge, float32)
	SetPositionPercent(ViewEdge, float32)
	SetPositionAuto(ViewEdge)

	// Determine if view changes require layout
	IsDirty() bool
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

func (v ViewPosition) String() string {
	switch v {
	case VIEW_POSITION_RELATIVE:
		return "VIEW_POSITION_RELATIVE"
	case VIEW_POSITION_ABSOLUTE:
		return "VIEW_POSITION_ABSOLUTE"
	default:
		return "[?? Invalid ViewPosition value]"
	}
}
