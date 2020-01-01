/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package timer

import (
	"sync"
	"time"

	"github.com/djthorpe/gopi/v2/base"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Timer struct {
	Bus gopi.Bus
}

type timer struct {
	eventId gopi.EventId                   // Current EventId
	stop    map[gopi.EventId]chan struct{} // Map of stop channels
	bus     gopi.Bus                       // Event bus

	base.Unit
	sync.Mutex
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Timer) Name() string { return "gopi.Timer" }

func (config Timer) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(timer)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if config.Bus == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("Missing Bus")
	} else {
		this.bus = config.Bus
		this.stop = make(map[gopi.EventId]chan struct{})
	}
	return this, nil
}

func (this *timer) Close() error {
	// Send stop signals
	for _, stop := range this.stop {
		close(stop)
	}

	// Wait until all stopped
	this.Wait()

	// Release resources
	this.stop = nil

	// Call other close
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Timer

func (this *timer) NewTicker(duration time.Duration) gopi.EventId {
	eventId := this.nextId()
	go func(timerId gopi.EventId, duration time.Duration) {
		stop := this.makeStop(eventId)
		ticker := time.NewTicker(duration)
		this.Add(1)
	FOR_LOOP:
		for {
			select {
			case <-ticker.C:
				this.bus.Emit(newTimerEvent(this, eventId))
			case <-stop:
				ticker.Stop()
				break FOR_LOOP
			}
		}
		this.Done()
	}(eventId, duration)
	return eventId
}

func (this *timer) NewTimer(duration time.Duration) gopi.EventId {
	eventId := this.nextId()
	go func(timerId gopi.EventId, duration time.Duration) {
		stop := this.makeStop(eventId)
		timer := time.NewTimer(duration)
		this.Add(1)
	FOR_LOOP:
		for {
			select {
			case <-timer.C:
				this.bus.Emit(newTimerEvent(this, eventId))
			case <-stop:
				timer.Stop()
				break FOR_LOOP
			}
		}
		this.Done()
	}(eventId, duration)
	return eventId
}

func (this *timer) Cancel(eventId gopi.EventId) error {
	this.Lock()
	defer this.Unlock()
	if stop, exists := this.stop[eventId]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("eventId")
	} else {
		delete(this.stop, eventId)
		close(stop)
	}
	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *timer) nextId() gopi.EventId {
	this.Lock()
	defer this.Unlock()
	this.eventId = this.eventId + 1
	return this.eventId
}

func (this *timer) makeStop(eventId gopi.EventId) chan struct{} {
	this.Lock()
	defer this.Unlock()
	stop := make(chan struct{})
	this.stop[eventId] = stop
	return stop
}
