package event

import (
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi"
)

type Publisher struct {
	sync.Mutex
	channels []chan gopi.Event
}

// Subscribe returns a new channel on which emitting events can occur
func (this *Publisher) Subscribe() <-chan gopi.Event {
	this.Lock()
	defer this.Unlock()

	// Create channels with a capacity of one
	if this.channels == nil {
		this.channels = make([]chan gopi.Event, 0, 1)
	}
	// Return a new channel
	channel := make(chan gopi.Event)
	this.channels = append(this.channels, channel)
	return channel
}

// Unsubscribe closes a channel and removes it from the list
// of channels which emitting can happen on
func (this *Publisher) Unsubscribe(subscriber <-chan gopi.Event) {
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
func (this *Publisher) Close() {
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
}

// Emit an event onto all subscriber channels, this method
// will block if the subscribers are not processing incoming
// events
func (this *Publisher) Emit(evt gopi.Event) {
	this.Lock()
	defer this.Unlock()

	if this.channels != nil {
		for _, channel := range this.channels {
			if channel != nil {
				channel <- evt
			}
		}
	}
}
