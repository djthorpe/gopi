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

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Timer struct{}

type timer struct {
	gopi.UnitBase
	sync.Mutex
	sync.WaitGroup

	// Current Id
	timerId gopi.TimerId
	// Map of stop channels
	stop map[gopi.TimerId]chan struct{}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Timer) Name() string { return "gopi.Timer" }

func (config Timer) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(timer)
	if err := this.UnitBase.Init(log); err != nil {
		return nil, err
	} else {
		this.stop = make(map[gopi.TimerId]chan struct{})
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
	return this.UnitBase.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Timer

func (this *timer) NewTicker(duration time.Duration) gopi.TimerId {
	timerId := this.nextId()
	go func(timerId gopi.TimerId, duration time.Duration) {
		stop := this.makeStop(timerId)
		ticker := time.NewTicker(duration)
		this.Add(1)
	FOR_LOOP:
		for {
			select {
			case <-ticker.C:
				this.Log.Debug("Ticker", timerId)
			case <-stop:
				ticker.Stop()
				break FOR_LOOP
			}
		}
		this.Done()
	}(timerId, duration)
	return timerId
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *timer) nextId() gopi.TimerId {
	this.Lock()
	defer this.Unlock()
	this.timerId = this.timerId + 1
	return this.timerId
}

func (this *timer) makeStop(timerId gopi.TimerId) chan struct{} {
	this.Lock()
	defer this.Unlock()
	stop := make(chan struct{})
	this.stop[timerId] = stop
	return stop
}
