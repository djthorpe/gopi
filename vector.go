/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
)

type Point struct {
	X, Y float32
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ZeroPoint = Point{0, 0}
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
