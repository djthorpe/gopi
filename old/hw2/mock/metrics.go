/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mock

import (
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	counter "github.com/djthorpe/gopi/util/metrics"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Metrics struct{}

type metrics struct {
	log      gopi.Logger
	counters map[chan<- uint]*metric
}

type metric struct {
	Counter *counter.Counter
	Name    string
	Type    gopi.MetricType
	Rate    gopi.MetricRate
	Chan    chan uint
	Done    chan struct{}
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// Timestamp for module creation
	ts = time.Now()
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open creates a new metrics object, returns error if not possible
func (config Metrics) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.mock.Metrics>Open{}")

	// create new driver
	this := new(metrics)
	this.log = log
	this.counters = make(map[chan<- uint]*metric, 0)

	// return driver
	return this, nil
}

// Close connection
func (this *metrics) Close() error {
	this.log.Debug("<sys.hw.mock.Metrics>Close{}")

	// Close all counter channels
	for _, metric := range this.counters {
		metric.Done <- gopi.DONE
		<-metric.Done
		close(metric.Chan)
	}

	// Release resources
	this.counters = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SYSTEM METRICS INTERFACE IMPLEMENTATION

func (this *metrics) UptimeHost() time.Duration {
	return time.Since(ts)
}

func (this *metrics) UptimeApp() time.Duration {
	return time.Since(ts) + time.Second
}

func (this *metrics) LoadAverage() (float64, float64, float64) {
	// Output some fake numbers
	return 0.1, 0.2, 0.3
}

////////////////////////////////////////////////////////////////////////////////
// RATE METRICS INTERFACE IMPLEMENTATION

func (this *metrics) NewCounter(metric_type gopi.MetricType, metric_rate gopi.MetricRate, name string) (chan<- uint, error) {
	this.log.Debug2("<sys.hw.mock.Metrics>NewCounter{ type=%v rate=%v name='%v' }", metric_type, metric_rate, name)

	// Create a new counter, append to list of existing counters
	if c := counter.NewCounter(metric_rate); c == nil {
		return nil, gopi.ErrBadParameter
	} else {
		m := &metric{
			Counter: c,
			Name:    name,
			Type:    metric_type,
			Rate:    metric_rate,
			Chan:    make(chan uint),
			Done:    make(chan struct{}),
		}
		this.counters[m.Chan] = m
		// Consume the value
		go this.consume(m)
		// Return the counter
		return m.Chan, nil
	}
}

func (this *metrics) Metric(counter chan<- uint) *gopi.Metric {
	if m, exists := this.counters[counter]; exists == false {
		return nil
	} else {
		sum, samples, length := m.Counter.Sum()
		return &gopi.Metric{
			Rate:  m.Rate,
			Type:  m.Type,
			Name:  m.Name,
			Mean:  float64(sum) / float64(samples),
			Total: uint(float64(sum) * float64(length) / float64(samples)),
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *metrics) consume(m *metric) {
FOR_LOOP:
	for {
		select {
		case value := <-m.Chan:
			if value > 0 {
				m.Counter.Increment(time.Now(), value)
			}
		case <-m.Done:
			break FOR_LOOP
		}
	}
	close(m.Done)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *metrics) String() string {
	return fmt.Sprintf("<sys.hw.mock.Metrics>{}")
}
