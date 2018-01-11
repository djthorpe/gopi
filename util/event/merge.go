/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

// Merge multiple incoming events into a single outgoing event channel
package event

import (
	"fmt"
	"reflect"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type EventMerger struct {
	change chan struct{}
	in     []<-chan gopi.Event
	out    chan<- gopi.Event
}

////////////////////////////////////////////////////////////////////////////////
// NEW

// Create an event merger object and start listening on
// incoming channels
func NewEventMerger(channels ...<-chan gopi.Event) *EventMerger {
	this := new(EventMerger)
	this.change = make(chan struct{})
	this.in = make([]<-chan gopi.Event, len(channels))
	this.out = make(chan<- gopi.Event)
	for i := range channels {
		this.in[i] = channels[i]
	}

	//  Start merger in background
	go this.mergeInBackground()

	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Channel returns the output channel
func (this *EventMerger) Channel() chan<- gopi.Event {
	return this.out
}

// Add an input channel
func (this *EventMerger) Add(new_channel <-chan gopi.Event) {
	this.in = append(this.in, new_channel)
	// Signal change
	this.change <- struct{}{}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *EventMerger) cases() []reflect.SelectCase {
	cases := make([]reflect.SelectCase, 1, len(this.in)+1)
	// Add the change channel
	cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(this.change),
	}
	// Add the remaining channels - ignoring nil channels
	// which have been closed
	for i := range this.in {
		if this.in[i] != nil {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(this.in[i]),
			})
		}
	}
	// TODO: return nil if all channels closed
	return cases
}

func (this *EventMerger) mergeInBackground() {
	// Continue loop until all in channels are closed
	cases := this.cases()
	for {
		// select - TODO: Cache the cases structure, we don't want to rebuild
		// it everytime a message is received
		chosen, v, ok := reflect.Select(cases)
		if chosen == 0 {
			// We need to reload the cases. If zero then end
			if cases = this.cases(); len(cases) == 0 {
				fmt.Println("FINISHED")
				close(this.out)
				close(this.change)
				return
			} else {
				fmt.Println("RELOAD")
				continue
			}
		} else if ok {
			this.out <- v.Interface().(gopi.Event)
		} else {
			// TODO: Remove channel from list of cases
			fmt.Printf("reflect.Select failed on channel %v\n", chosen-1)
		}
	}
}
