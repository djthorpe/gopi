package base

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PubSub struct {
	sync.Mutex
	channels []chan interface{}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Subscribe returns a new channel on which emitting events can occur
func (this *PubSub) Subscribe() <-chan interface{} {
	this.Lock()
	defer this.Unlock()

	// Create channels with a capacity of one
	if this.channels == nil {
		this.channels = make([]chan interface{}, 0, 1)
	}
	// Return a new channel
	channel := make(chan interface{})
	this.channels = append(this.channels, channel)
	return channel
}

// Unsubscribe closes a channel and removes it from the list
// of channels which emitting can happen on
func (this *PubSub) Unsubscribe(subscriber <-chan interface{}) {
	this.Lock()
	defer this.Unlock()

	if this.channels != nil {
		for i := range this.channels {
			if this.channels[i] == subscriber {
				close(this.channels[i])
				this.channels[i] = nil
			}
		}
	}
}

// Close will unsubscribe all remaining channels
func (this *PubSub) Close() error {
	this.Lock()
	defer this.Unlock()

	if this.channels != nil {
		for _, subscriber := range this.channels {
			if subscriber != nil {
				close(subscriber)
			}
		}
		this.channels = nil
	}

	// Always return success
	return nil
}

// Emit an event onto all subscriber channels, this method
// will block until all subscribers receive the value
func (this *PubSub) Emit(value interface{}) {
	this.Lock()
	defer this.Unlock()

	if this.channels != nil {
		for _, channel := range this.channels {
			if channel != nil {
				channel <- value
			}
		}
	}
}

func (this *PubSub) String() string {
	return "<gopi.PubSub subscribers=" + fmt.Sprint(len(this.channels)) + ">"
}
