package event_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// GET PUBLISHER OBJECT

func TestPublisher_000(t *testing.T) {
	publisher := &event.Publisher{}
	defer publisher.Close()
	if ch := publisher.Subscribe(); ch == nil {
		t.Error("Expected channel, got nil")
	} else {
		publisher.Unsubscribe(ch)
	}
}
