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

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Timer struct{}

type timer struct {
	log                 gopi.Logger
	subscribers         []chan gopi.Event
	channels            []reflect.SelectCase
	units               map[int]*unit
	done, reload, done2 chan struct{}
}

type unit struct {
	timer    *time.Timer
	ticker   *time.Ticker
	userInfo interface{}
}

type event struct {
	source    gopi.Driver
	userInfo  interface{}
	timestamp time.Time
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the driver
func (config Timer) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("sys.timer.Open{ }")

	this := new(timer)
	this.log = log
	this.subscribers = make([]chan gopi.Event, 0)
	this.units = make(map[int]*unit, 0)
	this.channels = make([]reflect.SelectCase, 2)
	this.done = make(chan struct{})
	this.reload = make(chan struct{})
	this.done2 = make(chan struct{})

	// Zero'th channel is the done signal which
	// is emitted when Close() is called
	this.channels[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(this.done),
	}
	// First channel is the reload signal when
	// a new ticker or timer is added
	this.channels[1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(this.reload),
	}

	// Background go routine - waits for done or reload
	go this.wait_for_timers()

	return this, nil
}

// Close the driver
func (this *timer) Close() error {
	this.log.Debug("sys.timer.Close{ }")

	// Send a done signal to quit background, and wait for end
	this.done <- gopi.DONE
	_ = <-this.done2

	// stop timer & ticker resources
	for _, unit := range this.units {
		if unit.ticker != nil {
			unit.ticker.Stop()
		}
		if unit.timer != nil {
			unit.timer.Stop()
		}
	}

	this.units = nil
	this.channels = nil
	this.subscribers = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - TIMERS

// Schedule a timeout (one shot)
func (this *timer) NewTimeout(duration time.Duration, userInfo interface{}) {
	timer := time.NewTimer(duration)
	this.channels = append(this.channels, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(timer.C),
	})
	this.units[len(this.channels)-1] = &unit{
		timer:    timer,
		userInfo: userInfo,
	}
	this.reload <- gopi.DONE
}

// Schedule an interval, which can fire immediately
func (this *timer) NewInterval(duration time.Duration, userInfo interface{}, immediately bool) {
	ticker := time.NewTicker(duration)
	this.channels = append(this.channels, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ticker.C),
	})
	unit := &unit{
		ticker:   ticker,
		userInfo: userInfo,
	}
	this.units[len(this.channels)-1] = unit
	this.reload <- gopi.DONE
	if immediately {
		this.emit(unit, time.Now())
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *timer) String() string {
	return fmt.Sprintf("<sys.timer>{ timers=%v }", this.units)
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - EVENTS

// Subscribe to events emitted. Returns unique subscriber
// identifier and channel on which events are emitted
func (this *timer) Subscribe() chan gopi.Event {
	this.log.Debug2("<sys.timer.Subscribe>{ }")

	// Create a new channel for emitting events
	subscriber := make(chan gopi.Event)
	this.subscribers = append(this.subscribers, subscriber)
	return subscriber
}

// Unsubscribe from events emitted
func (this *timer) Unsubscribe(subscriber chan gopi.Event) {
	this.log.Debug2("<sys.timer.Unsubscribe>{ }")

	for i := range this.subscribers {
		if this.subscribers[i] == subscriber {
			this.subscribers[i] = nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// EVENT INTERFACE

func (this *event) Name() string {
	return "TimerEvent"
}

func (this *event) Source() gopi.Driver {
	return this.source
}

func (this *event) Timestamp() time.Time {
	return this.timestamp
}

func (this *event) UserInfo() interface{} {
	return this.userInfo
}

func (this *event) String() string {
	return fmt.Sprintf("<sys.timer.event>{ ts=%v userInfo=%v }", this.timestamp, this.userInfo)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *timer) emit(userInfo *unit, value time.Time) {
	event := &event{source: this, userInfo: userInfo.userInfo, timestamp: value}
	for _, channel := range this.subscribers {
		if channel != nil {
			channel <- event
		}
	}
}

// wait_for_timers will wait for an event on any channel in the list of
// channels. The zero'th channel is the 'done' signal which indicates Close()
// has been called
func (this *timer) wait_for_timers() {

	for {
		if chosen, value, ok := reflect.Select(this.channels); ok && chosen == 0 {
			// Break out
			break
		} else if ok && chosen == 1 {
			// Reload
			continue
		} else if ok && chosen < len(this.channels) {
			if unit, exists := this.units[chosen]; exists {
				this.emit(unit, value.Interface().(time.Time))
			}
		}
	}
	this.done2 <- gopi.DONE
}
