package event

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/djthorpe/gopi/v3"
)

func Test_Event_000(t *testing.T) {
	pub := &publisher{}
	if ch := pub.Subscribe(); ch == nil {
		t.Error("Unexpected nil return value")
	} else {
		pub.Unsubscribe(ch)
	}
}

func Test_Event_001(t *testing.T) {
	var wg sync.WaitGroup

	pub := &publisher{}
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

func Test_Event_002(t *testing.T) {
	pub := &publisher{}
	evts := 0
	rcvs := rand.Int() % 20
	total := 100

	var wg sync.WaitGroup

	// Receive events
	recv := func(ch <-chan gopi.Event) {
		for _ = range ch {
			t.Log("got", evts)
			evts += 1
		}
		wg.Done()
	}

	// Receive events in the background
	for i := 0; i < rcvs; i++ {
		wg.Add(1)
		go recv(pub.Subscribe())
	}

	// Emit events
	for i := 0; i < total; i++ {
		pub.Emit(nil)
	}

	// Wait for all goroutinnes completed
	wg.Wait()

	// Check for number of events
	if evts != total*rcvs {
		t.Error("Unexpected number of events,", evts, "!=", total*rcvs)
	}
}
