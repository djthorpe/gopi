package gopi

import (
	"time"
)

/*
	This file contains definitions for transmission of measurement
	data and querying data
*/

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Metrics provides a mechanism for defining measurements
// and emitting data, which may be time-series based and include
// tags (which are indexed and can be used for grouping) and
// metrics (which are not, and can be aggregated).
type Metrics interface {
	// Define a measurement with metric definitions and optional tag fields
	NewMeasurement(string, string, ...Field) (Measurement, error)

	// Emit tags and metrics for a named measurement, omitting timestamp
	Emit(string, []Field, ...interface{}) error

	// EmitTS emits tags and metrics for a named measurement, with defined timestamp
	EmitTS(string, time.Time, []Field, ...interface{}) error

	// Measurements returns array of all defined measurements
	Measurements() []Measurement

	// Field creates a field or nil if invalid
	Field(string, ...interface{}) Field

	// HostTag returns a field with the current hostname
	HostTag() Field

	// UserTag returns a field with the current username
	UserTag() Field

	// EnvTag returns a field with the value of an environment variable
	EnvTag(string) Field
}

// MetricWriter implements a database writing object
type MetricWriter interface {
	Ping() (time.Duration, error) // Ping the database and return latency
	Write(...Measurement) error   // Write one or more measurements to the database
}

// MetricReader implements a database query
type MetricReader interface {
	Ping() (time.Duration, error) // Ping the database and return latency

	// NewQuery constructs a query with measurement name and the
	// names of tags which are grouped. By default, all metrics are returned
	NewQuery(string, ...string) MetricQuery
}

// MetricQuery constructs queries which can be executed by the MetricReader,
// TODO methods which modify the basic query
type MetricQuery interface{}

// Measurement is a single data point
type Measurement interface {
	Event

	Time() time.Time  // Time returns the timestamp for the data point or time.Time{}
	Tags() []Field    // Return the dimensions/tags for the data point
	Metrics() []Field // Return the metrics for the data point

	Get(string) interface{}        // Return a field value
	Set(string, interface{}) error // Set a field value
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
