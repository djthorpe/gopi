/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package timer

import (
	"fmt"
	"reflect"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Timer struct{}

type timer struct {
	log      gopi.Logger
	channels []reflect.SelectCase
	units    map[int]*unit

	event.Publisher
	event.Tasks
}

type unit struct {
	timer        *time.Timer
	ticker       *time.Ticker
	userInfo     interface{}
	counter      uint
	duration     time.Duration
	max_duration time.Duration
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the timer
func (config Timer) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("sys.timer.Open{ }")

	this := new(timer)
	this.log = log
	this.channels = make([]reflect.SelectCase, 1)
	this.channels[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(make(chan event.Signal)),
	}
	this.units = make(map[int]*unit)

	// Background go routine - waits for reload, stop or ticker events
	this.Tasks.Start(this.wait_for_timers)

	return this, nil
}

// Close the timer
func (this *timer) Close() error {
	this.log.Debug("sys.timer.Close{ }")

	// Cancel all the timers
	for _, unit := range this.units {
		unit.Cancel()
	}

	// End the task
	if err := this.Tasks.Close(); err != nil {
		this.log.Warn("sys.timer.Close: %v", err)
	}

	// Unsubscribe and close
	this.Publisher.Close()

	// Blank out instance variables
	this.channels = nil
	this.units = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - TIMERS

// NewTimeout schedules a on-shot timer
func (this *timer) NewTimeout(duration time.Duration, userInfo interface{}) error {
	this.log.Debug2("sys.timer.NewTimeout{ duration=%v userInfo=%v }", duration, userInfo)

	// Create the timeout, append channel and reload
	timer := time.NewTimer(duration)
	this.append(reflect.ValueOf(timer.C), &unit{
		timer:    timer,
		userInfo: userInfo,
	})
	// Success
	return nil
}

// NewInterval schedules a periodic firing, which can fire immediately
func (this *timer) NewInterval(duration time.Duration, userInfo interface{}, immediately bool) error {
	this.log.Debug2("sys.timer.NewInterval{ duration=%v userInfo=%v immediately=%v }", duration, userInfo, immediately)

	// Create the ticker, append channel and reload
	ticker := time.NewTicker(duration)
	idx := this.append(reflect.ValueOf(ticker.C), &unit{
		ticker:   ticker,
		userInfo: userInfo,
	})
	if immediately {
		this.emit(idx, time.Now())
	}
	// Success
	return nil
}

// NewBackoff schedules a backoff timer with maximum backoff duration
func (this *timer) NewBackoff(duration time.Duration, max_duration time.Duration, userInfo interface{}) error {
	this.log.Debug2("sys.timer.NewBackoff{ duration=%v max_duration=%v userInfo=%v }", duration, max_duration, userInfo)

	// Check for zero durations
	if duration == 0 || max_duration == 0 {
		return gopi.ErrBadParameter
	}
	if max_duration <= duration {
		return gopi.ErrBadParameter
	}

	// Create the timeout, append channel and reload
	timer := time.NewTimer(duration)
	idx := this.append(reflect.ValueOf(timer.C), &unit{
		timer:        timer,
		userInfo:     userInfo,
		duration:     duration,
		max_duration: max_duration,
	})

	// Emit the event immediately
	this.emit(idx, time.Now())

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *timer) String() string {
	return fmt.Sprintf("<sys.timer>{ timers=%v }", this.units)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *timer) emit(idx int, ts time.Time) {
	if u, ok := this.units[idx]; ok == false {
		this.log.Warn("sys.timer.emit: Invalid index, %v", idx)
	} else {
		// Increment the counter (number of times fired)
		u.counter = u.counter + 1
		// Emit the event
		this.Emit(NewTimerEvent(this, u, ts))
		// If this is a backoff timeout and the counter is above 1
		// (it's not the immediate firing) then reset the backoff to double
		// the current interval, up to a maximum of max_duration
		if u.max_duration > 0 && u.counter > 1 && u.timer != nil {
			u.duration *= 2
			if u.duration > u.max_duration {
				u.duration = u.max_duration
			}
			this.log.Debug2("sys.timer.emit: backoff interval=%v for %v", u.duration, u)
			u.timer.Reset(u.duration)
		}
	}
}

func (this *timer) append(c reflect.Value, u *unit) int {
	// append channel and reload
	this.channels = append(this.channels, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: c,
	})
	idx := len(this.channels) - 1
	if u != nil {
		this.units[idx] = u
		// send a signal to the zero'th channel
		this.channels[0].Chan.Send(reflect.ValueOf(gopi.DONE))
	}
	// return the index of the channel
	return idx
}

// wait_for_timers will wait for an event on any channel in the list of
// channels
func (this *timer) wait_for_timers(start chan<- event.Signal, stop <-chan event.Signal) error {
	// Indicate we are now running
	start <- gopi.DONE

	// Append stop signal but don't reload
	stop_idx := this.append(reflect.ValueOf(stop), nil)

FOR_LOOP:
	for {
		// Wait for reload, stop or a timer maturing
		chosen, _, ok := reflect.Select(this.channels)
		// We received a reload signal
		if chosen == 0 && ok {
			continue
		}
		// We received a stop signal so break out of the loop
		if chosen == stop_idx {
			break FOR_LOOP
		}
		// Otherwise, assume we received a timer maturing
		if ok {
			this.emit(chosen, time.Now())
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *unit) Cancel() {
	if this.timer != nil {
		this.timer.Stop()
	}
	if this.ticker != nil {
		this.ticker.Stop()
	}
}

func (this *unit) String() string {
	return fmt.Sprintf("<sys.timer>{ userInfo=%v }", this.userInfo)
}
