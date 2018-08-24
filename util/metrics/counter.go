/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Provide a way to count metrics and provide total and rates
package metrics

import (
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

type Counter struct {
	sync.Mutex
	values map[int]uint
	rate   gopi.MetricRate
	len    int
}

func NewCounter(rate gopi.MetricRate) *Counter {
	this := new(Counter)
	this.len = numberOfBuckets(rate)
	if this.len == 0 {
		// Unsupported rate
		return nil
	}
	this.values = make(map[int]uint, this.len)
	this.rate = rate
	return this
}

func (this *Counter) Increment(ts time.Time, value uint) {
	this.Lock()
	defer this.Unlock()

	bucket := bucketForTimestamp(ts, this.rate)
	if _, exists := this.values[bucket]; exists == false {
		// If the bucket doesn't exist, then create it and
		// potentially delete the next bucket
		this.values[bucket] = 0
		next_bucket := (bucket + 1) % this.len
		if _, next_exists := this.values[next_bucket]; next_exists {
			delete(this.values, next_bucket)
		}
	}
	this.values[bucket] += value
}

// Return the sum and the number of samples, and total number
func (this *Counter) Sum() (uint, int, int) {
	sum := uint(0)
	for _, v := range this.values {
		sum += v
	}
	return sum, len(this.values), this.len
}

func bucketForTimestamp(ts time.Time, rate gopi.MetricRate) int {
	switch rate {
	case gopi.METRIC_RATE_SECOND:
		return int(ts.Second())
	case gopi.METRIC_RATE_MINUTE:
		return int(ts.Minute())
	case gopi.METRIC_RATE_HOUR:
		return int(ts.Hour())
	default:
		return 0
	}
}

func numberOfBuckets(rate gopi.MetricRate) int {
	switch rate {
	case gopi.METRIC_RATE_SECOND:
		return 60
	case gopi.METRIC_RATE_MINUTE:
		return 60
	case gopi.METRIC_RATE_HOUR:
		return 24
	default:
		return 0
	}
}
