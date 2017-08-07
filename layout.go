/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"
import (
	"io"
)

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
	LAYOUT_DIRECTION_NONE    Direction = iota
	LAYOUT_DIRECTION_INHERIT           = LAYOUT_DIRECTION_NONE
	LAYOUT_DIRECTION_LEFTRIGHT
	LAYOUT_DIRECTION_RIGHTLEFT
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

// View defines a 2D rectangular drawable element
type View interface {
}

// Layout defines the methods of calculating layout of views within
// a rectangular surface (window)
type Layout interface {
	Driver

	// Return the root view for this layout
	View() View

	// Calculate layout using the root node as the correct size
	CalculateLayout() bool

	// Calculate layout with a new root size
	CalculateLayoutWithSize(w, h float32) bool

	// Return XML encoded version of the layout
	Encode(w io.Writer, indent EncodeIndentOptions) error
}

/*
type LayoutConfig struct {
	Direction Direction
	Width     float32
	Height    float32
}
*/
