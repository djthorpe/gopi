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
// STRINGIFY

func (p Point) String() string {
	return fmt.Sprintf("gopi.Point{ %v,%v }", p.X, p.Y)
}
