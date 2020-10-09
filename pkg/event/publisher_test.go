package event_test

import (
	"sync"
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/event"
)

func Test_Event_000(t *testing.T) {
	pub := &event.Publisher{}
	if ch := pub.Subscribe(); ch == nil {
		t.Error("Unexpected nil return value")
	} else {
		pub.Unsubscribe(ch)
	}
}

func Test_Event_001(t *testing.T) {
	var wg sync.WaitGroup

	pub := &event.Publisher{}
	evts := 0
	total := 100
	ch := pub.Subscribe()

	// Receive events
	go func() {
		wg.Add(1)
		for _ = range ch {
			evts += 1
		}
		wg.Done()
	}()

	// Emit events
	for i := 0; i < total; i++ {
		pub.Emit(nil)
	}

	// Unsubscribe channel
	pub.Unsubscribe(ch)

	// Wait for end of goroutine
	wg.Wait()

	// Check for number of events
	if evts != total {
		t.Error("Unexpected number of events,", evts, "!=", total)
	}
}

/*
func Test_Event_002(t *testing.T) {
	pub := &event.Publisher{}
	evts := 0
	total := 100

	// TODO

	// Receive events
	recv := func() {
		ch := pub.Subscribe()
		for _ = range ch {
			evts += 1
		}
		pub.Unsubscribe(ch)
	}

	// Receive events and increment counter
	go recv()
	go recv()

	// Emit events
	for i := 0; i < total; i++ {
		pub.Emit(nil)
	}

	// Check for number of events
	if evts != total {
		t.Error("Unexpected number of events,", evts, "!=", total)
	}
}
*/
