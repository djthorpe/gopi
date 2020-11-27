<<<<<<< HEAD
=======
/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

>>>>>>> master
package gopi

import (
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Metrics provides a mechanism for defining measurements
// and emitting data, which may be time-series based and include
// tags/dimensions (which are indexed) and metrics (which are not)
type Metrics interface {
	// Define a measurement with metric definitions and optional tag fields
	NewMeasurement(string, string, ...Field) (Measurement, error)

	// NewFields returns array of fields
	NewFields(...string) []Field

	// Emit metrics for a named measurement, omitting timestamp
	Emit(string, ...interface{}) error

	// EmitTS emits metrics for a named measurement, with defined timestamp
	EmitTS(string, time.Time, ...interface{}) error

	// Measurements returns array of all defined measurements
	Measurements() []Measurement
}

// Measurement is a single data point
type Measurement interface {
	Event

	Time() time.Time  // Time returns the timestamp for the data point or time.Time{}
	Tags() []Field    // Return the dimenstions/tags for the data point
	Metrics() []Field // Return the metrics for the data point
}

type Field interface {
	// Name returns field name
	Name() string

	// Kind returns kind of field or nil
	Kind() string

	// IsNil returns true if value is nil
	IsNil() bool

	// Value returns field value, or nil
	Value() interface{}

	// SetValue sets specific value and returns error if unsupported
	SetValue(interface{}) error

	// Copy returns a copy of the field
	Copy() Field
}
