package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

type publisher struct {
	gopi.Unit
	sync.RWMutex

	q  chan gopi.Event
	ch []chan gopi.Event
}

const (
	// queuesize defines the buffer of events, in case the receiver is not
	// quick at picking up events compared to sender
	queuesize = 100
)

func (this *publisher) New(gopi.Config) error {
	this.q = make(chan gopi.Event, queuesize)
	return nil
}

func (this *publisher) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Close queue
	close(this.q)

	// Unsubscribe channels
	for _, ch := range this.ch {
		if ch != nil {
			close(ch)
		}
	}

	// Dispose
	this.q = nil
	this.ch = nil

	// Return success
	return nil
}

func (this *publisher) Run(ctx context.Context) error {
	for {
		select {
		case evt := <-this.q:
			this.RWMutex.RLock()
			for _, ch := range this.ch {
				if ch != nil {
					ch <- evt
				}
			}
			this.RWMutex.RUnlock()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (this *publisher) Subscribe() <-chan gopi.Event {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	ch := make(chan gopi.Event)
	this.ch = append(this.ch, ch)
	return ch
}

func (this *publisher) Unsubscribe(ch <-chan gopi.Event) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	for i, other := range this.ch {
		if other == ch {
			close(other)
			this.ch[i] = nil
		}
	}
}

func (this *publisher) Emit(evt gopi.Event, block bool) error {
	// Use NullEvent when evt is nil
	if evt == nil {
		evt = NewNullEvent()
	}

	// Blocking case
	if block {
		this.q <- evt
		return nil
	}

	// Non-blocking case
	select {
	case this.q <- evt:
		return nil
	default:
		return gopi.ErrChannelFull.WithPrefix(evt.Name())
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *publisher) String() string {
	str := "<publisher"
	if this == nil {
		str += " nil"
	} else {
		str += " ch=" + fmt.Sprint(this.ch)
	}
	return str + ">"
}
