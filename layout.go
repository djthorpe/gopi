/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"

/////////////////////////////////////////////////////////////////////
// TYPES

// LayoutDirection is either left to right, right to left, or inherited
// from parent
type LayoutDirection uint

// EncodeIndentOptions defines XML indent style when writing out
type EncodeIndentOptions uint

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	LAYOUT_DIRECTION_NONE LayoutDirection = iota
	LAYOUT_DIRECTION_LEFTRIGHT
	LAYOUT_DIRECTION_RIGHTLEFT
	LAYOUT_DIRECTION_INHERIT = LAYOUT_DIRECTION_NONE
)

/////////////////////////////////////////////////////////////////////
// STRUCTS

type ViewStyle struct {
	Top    string
	Left   string
	Bottom string
	Right  string
	Width  string
	Height string
}

/////////////////////////////////////////////////////////////////////
// INTERFACES

// View defines a 2D rectangular drawable element
type View interface {
	// Return tag for this view
	Tag() uint

	// Return class for this view
	Class() string

	// Return style for this view
	Style() *ViewStyle
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
