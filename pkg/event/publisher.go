package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

type publisher struct {
	gopi.Unit
	sync.Mutex

	q  chan gopi.Event
	ch []chan gopi.Event
}

const (
	queuesize = 10
)

func (this *publisher) New(gopi.Config) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.q = make(chan gopi.Event, queuesize)
	return nil
}

func (this *publisher) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Unsubscribe channels
	for _, ch := range this.ch {
		if ch != nil {
			close(ch)
		}
	}

	// Dispose
	close(this.q)
	this.q = nil
	this.ch = nil

	// Return success
	return nil
}

func (this *publisher) Run(ctx context.Context) error {
	fmt.Println("RUN")
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		case evt := <-this.q:
			for _, ch := range this.ch {
				if ch != nil && evt != nil {
					ch <- evt
				}
			}
		}
	}
	fmt.Println("RUN FINISHED")
	return nil
}

func (this *publisher) Subscribe() <-chan gopi.Event {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	ch := make(chan gopi.Event)
	this.ch = append(this.ch, ch)
	return ch
}

func (this *publisher) Unsubscribe(ch <-chan gopi.Event) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	for i, other := range this.ch {
		if other == ch {
			close(other)
			this.ch[i] = nil
		}
	}
}

func (this *publisher) Emit(evt gopi.Event) error {
	// Use NullEvent when evt is nil
	if evt == nil {
		evt = NewNullEvent()
	}

	// Emit and return error if cannot emit
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
