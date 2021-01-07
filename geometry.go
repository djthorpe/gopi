package gopi

import (
	"fmt"
)

/*
	This file contains interface defininitons for graphics geometry:

	* Points
	* Rectangles

	Dimensions are defined in float32 rather than uint to support
	vector graphics.
*/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Point struct {
	X, Y float32
}

type Size struct {
	W, H float32
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ZeroPoint = Point{0, 0}
	ZeroSize  = Size{0, 0}
)

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS

func (p1 Point) Equals(p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p Point) String() string {
	return fmt.Sprintf("gopi.Point{ %.1f,%.1f }", p.X, p.Y)
}

func (s Size) String() string {
	return fmt.Sprintf("gopi.Size{ %.1f,%.1f }", s.W, s.H)
}
