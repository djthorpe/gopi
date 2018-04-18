package event_test

import (
	"testing"
	"time"

	"github.com/djthorpe/gopi"
	evt "github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// EVENT IMPLEMENTATION

type event struct {
	value int
}

func (*event) Name() string {
	return "event"
}

func (*event) Source() gopi.Driver {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// CREATE CREATION

func Test_000(t *testing.T) {
	if pubsub := evt.NewPubSub(0); pubsub == nil {
		t.Fatal("Cannot create a pubsub object")
	}
}

func Test_001(t *testing.T) {
	if pubsub := evt.NewPubSub(0); pubsub == nil {
		t.Fatal("Cannot create a pubsub object")
	} else if c := pubsub.Subscribe(); c == nil {
		t.Fatal("Cannot subscribe")
	}
}

func Test_002(t *testing.T) {
	pubsub := evt.NewPubSub(0)
	c1 := pubsub.Subscribe()
	c2 := pubsub.Subscribe()
	if c1 == c2 || c1 == nil || c2 == nil {
		t.Error("c1==c2")
	}
}

func Test_003(t *testing.T) {
	pubsub := evt.NewPubSub(0)
	pubsub.Close()
	c := pubsub.Subscribe()
	if c != nil {
		t.Error("Expectung nil on subscribe for closed object")
	}
}

func Test_004(t *testing.T) {
	pubsub := evt.NewPubSub(0)
	c := pubsub.Subscribe()
	pubsub.Close()
	select {
	case evt := <-c:
		if evt != nil {
			t.Error("Expecting nil returned on closed pubsub object")
		}
	default:
		t.Error("Expecting nil returned on closed pubsub object")
	}
}

func Test_005(t *testing.T) {
	pubsub := evt.NewPubSub(0)
	c := pubsub.Subscribe()
	e := &event{}
	go func() {
		pubsub.Emit(e)
	}()
	select {
	case evt := <-c:
		if evt != e {
			t.Errorf("Expecting e to be emitted to channel, got %v", evt)
		}
	}
}

func Test_006(t *testing.T) {
	pubsub := evt.NewPubSub(0)
	c1 := pubsub.Subscribe()
	c2 := pubsub.Subscribe()
	e := &event{}
	go func() {
		pubsub.Emit(e)
	}()
	i := 0
FOR_LOOP:
	for {
		select {
		case evt := <-c1:
			if evt != e {
				t.Errorf("Expecting e to be emitted to channel C1, got %v", evt)
			} else {
				t.Log("Got C1 event")
				i = i + 1
			}
		case evt := <-c2:
			if evt != e {
				t.Errorf("Expecting e to be emitted to channel C2, got %v", evt)
			} else {
				t.Log("Got C2 event")
				i = i + 1
			}
		case <-time.After(100 * time.Millisecond):
			t.Log("Timeout")
			break FOR_LOOP
		}
	}
	if i != 2 {
		t.Errorf("Expected e to be emitted to both channels")
	}
}
