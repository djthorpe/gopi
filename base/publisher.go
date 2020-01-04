/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package base

import (
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Publisher struct {
	q    map[uint][]chan interface{}
	lock map[uint]*sync.Mutex

	sync.Mutex
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Publisher) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Wait until done
	this.WaitGroup.Wait()

	// close all queues
	for queue := range this.q {
		this.close(queue)
	}

	// Release resources
	this.q = nil
	this.lock = nil

	// Success
	return nil
}

func (this *Publisher) close(queue uint) {
	// Lock queue whilst closing
	this.lock[queue].Lock()
	defer this.lock[queue].Unlock()
	// Close channels
	for _, c := range this.q[queue] {
		close(c)
	}
	// Set nil array for queue
	this.q[queue] = nil
}

func (this *Publisher) Emit(queue uint, value interface{}) {
	if chans := this.newChannels(queue); chans != nil {
		// Lock queue whilst emitting
		this.lock[queue].Lock()
		defer this.lock[queue].Unlock()
		// Emit values
		for _, c := range chans {
			c <- value
		}
	}
}

// Subscribe receives emitted messages until Unsubscribe is called
func (this *Publisher) Subscribe(queue uint, capacity int, callback func(value interface{})) gopi.Channel {
	stop, start := make(chan struct{}), make(chan struct{})
	go func() {
		this.WaitGroup.Add(1)
		defer this.WaitGroup.Done()

		// Subscribe and indicate the subscription has occured by closing the
		// start channel
		evt := this.SubscribeInt(queue, capacity)
		close(start)
		// If error occured with subscribing, end go routine immediately
		if evt == nil {
			return
		}
		// Background go routine continues until stop signal is received
		go func() {
			<-stop
			this.UnsubscribeInt(evt)
		}()
		// Repeat accepting and handling events, and end when unsubscribe is called
	FOR_LOOP:
		for {
			select {
			case value := <-evt:
				if value == nil {
					break FOR_LOOP
				} else {
					callback(value)
				}
			}
		}
		// End end, close stop
		close(stop)
	}()
	// In main function, block until the start is received
	<-start
	// Return the stop channel
	return stop
}

// Unsubscribe waits until the receive loop has completed
func (this *Publisher) Unsubscribe(stop gopi.Channel) {
	stop <- struct{}{}
	<-stop
}

func (this *Publisher) Len(queue uint) int {
	return len(this.newChannels(queue))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Publisher) newChannels(queue uint) []chan interface{} {
	if this.q == nil || this.lock == nil {
		this.Mutex.Lock()
		this.q = make(map[uint][]chan interface{})
		this.lock = make(map[uint]*sync.Mutex)
		this.Mutex.Unlock()
	}
	chans, exists := this.q[queue]
	if exists == false {
		this.Mutex.Lock()
		this.q[queue] = make([]chan interface{}, 0, 1)
		this.lock[queue] = &sync.Mutex{}
		chans = this.q[queue]
		this.Mutex.Unlock()
	}
	return chans
}

func (this *Publisher) SubscribeInt(queue uint, capacity int) <-chan interface{} {
	if chans := this.newChannels(queue); chans == nil {
		return nil
	} else {
		new := make(chan interface{}, capacity)
		this.q[queue] = append(chans, new)
		return new
	}
}

func (this *Publisher) UnsubscribeInt(c <-chan interface{}) bool {
	for queue, chans := range this.q {
		for i, other := range chans {
			if other == c {
				close(other)
				this.lock[queue].Lock()
				this.q[queue] = append(chans[:i], chans[i+1:]...)
				this.lock[queue].Unlock()
				return true
			}
		}
	}
	return false
}
