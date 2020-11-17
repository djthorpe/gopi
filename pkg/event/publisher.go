package event

import (
	"fmt"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

type publisher struct {
	gopi.Unit
	sync.Mutex

	ch []chan gopi.Event
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

func (this *publisher) Emit(evt gopi.Event) {
	for _, ch := range this.ch {
		if ch != nil {
			ch <- evt
		}
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
