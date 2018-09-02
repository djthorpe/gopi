/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type (
	MetricRate uint
	MetricType uint
)

/////////////////////////////////////////////////////////////////////
// INTERFACE

// Metric is an abstract method to store values associated with
// a measurement
type Metric interface {
	// Return the metric rate (store values over a period)
	Rate() MetricRate

	// Return the metric type (the units used for the metric)
	Type() MetricType

	// Return the name of the metric
	Name() string

	// Return the unit for the metric (°C for example)
	Unit() string

	// Return the last metric value as a uint
	UintValue() uint

	// Return the last metric value as a float64
	FloatValue() float64
}

// Metrics returns various metrics for host and
// custom metrics
type Metrics interface {
	Driver

	// Uptimes for host and for application
	UptimeHost() time.Duration
	UptimeApp() time.Duration

	// Load Average (1, 5 and 15 minutes)
	LoadAverage() (float64, float64, float64)

	// Return metric channel which records uint values
	NewMetricUint(MetricType, MetricRate, string) (chan<- uint, error)

	// Return metric channel which records float64 values
	NewMetricFloat64(MetricType, MetricRate, string) (chan<- float64, error)

	// Return all metrics of a particular type, or METRIC_TYPE_NONE
	// for all metrics
	Metrics(MetricType) []Metric
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	METRIC_RATE_NONE MetricRate = iota
	METRIC_RATE_MINUTE
	METRIC_RATE_HOUR
	METRIC_RATE_DAY
)

const (
	METRIC_TYPE_NONE    MetricType = iota
	METRIC_TYPE_PURE               // Pure number
	METRIC_TYPE_CELCIUS            // Temperature
)

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v MetricRate) String() string {
	switch v {
	case METRIC_RATE_MINUTE:
		return "METRIC_RATE_MINUTE"
	case METRIC_RATE_HOUR:
		return "METRIC_RATE_HOUR"
	case METRIC_RATE_DAY:
		return "METRIC_RATE_DAY"
	default:
		return "[?? Invalid MetricRate value]"
	}
}

func (t MetricType) String() string {
	switch t {
	case METRIC_TYPE_PURE:
		return "METRIC_TYPE_PURE"
	case METRIC_TYPE_CELCIUS:
		return "METRIC_TYPE_CELCIUS"
	default:
		return "[?? Invalid MetricType value]"
	}
}
