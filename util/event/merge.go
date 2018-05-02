/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Merge multiple incoming events into one, and fan out to subscribers
package event

import (
	"reflect"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// EventMerger represents a way to merge events
// and subscribe the emitted events
type EventMerger interface {
	gopi.Publisher

	// Add a channel for merging
	Add(<-chan gopi.Event)

	// Close the event merger
	Close()
}

type merger struct {
	// change indicates that the cases need to be reloaded and done
	// indicates the background task is done
	change chan struct{}
	done   chan struct{}

	// all the channels which are being merged
	in []<-chan gopi.Event

	// the pubsub object for fanning out emitted events
	pubsub *PubSub
}

////////////////////////////////////////////////////////////////////////////////
// NEW AND CLOSE

// Create an event merger object and start listening on incoming channels
func NewEventMerger(channels ...<-chan gopi.Event) EventMerger {
	this := new(merger)
	this.change = make(chan struct{})
	this.done = make(chan struct{})
	this.in = make([]<-chan gopi.Event, len(channels))
	this.pubsub = NewPubSub(len(channels))

	// Obtain the channels
	for i := range channels {
		this.in[i] = channels[i]
	}

	// Start merger in background
	go this.mergeInBackground()

	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLISHER INTERFACE IMPLEMENTATION

func (this *merger) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

func (this *merger) Unsubscribe(subscriber <-chan gopi.Event) {
	this.pubsub.Unsubscribe(subscriber)
}

func (this *merger) Emit(evt gopi.Event) {
	go func() {
		// Perform in separate goroutine to prevent blocking
		this.pubsub.Emit(evt)
	}()
}

////////////////////////////////////////////////////////////////////////////////
// EVENTMERGER INTERFACE IMPLEMENTATION

// Add an input channel
func (this *merger) Add(new_channel <-chan gopi.Event) {
	this.in = append(this.in, new_channel)
	// Signal change
	this.change <- gopi.DONE
}

// Close channels and release resources
func (this *merger) Close() {
	// Close the pubsub object
	this.pubsub.Close()
	this.pubsub = nil
	// Empty channels array
	// It's the responsibility of the
	// adder to close those channels
	this.in = nil
	// Close change channel to indicate done
	close(this.change)
	// Wait for done signal
	<-this.done
	// Release done
	close(this.done)
	this.done = nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *merger) cases() []reflect.SelectCase {
	cases := make([]reflect.SelectCase, 1, len(this.in)+1)
	// Add the change channel
	cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(this.change),
	}
	// Add the remaining channels - ignoring nil channels
	// which have been closed
	if this.in != nil {
		for i := range this.in {
			if this.in[i] != nil {
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(this.in[i]),
				})
			}
		}
	}
	// return all cases
	return cases
}

func (this *merger) mergeInBackground() {
	// Continue loop until chanhe channel is closed
	cases := this.cases()
FOR_LOOP:
	for {
		// Deal with zero cases condition
		if len(cases) == 0 {
			break FOR_LOOP
		}
		// select cases
		i, v, ok := reflect.Select(cases)
		if i == 0 && ok == false {
			// We need to reload the cases. If zero then end
			break FOR_LOOP
		} else if i == 0 {
			// Reload cases
			cases = this.cases()
		} else if ok {
			this.Emit(v.Interface().(gopi.Event))
		} else if i > 0 {
			// Set channel to nil to remove from cases
			this.in[i-1] = nil
			// Rebuild cases
			cases = this.cases()
		}
	}
	// Indicate the background thread is done
	this.done <- gopi.DONE
}
