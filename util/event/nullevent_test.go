package event_test

import (
	"testing"

	// Frameworks

	"github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// GET NULL EVENT

func TestNullEvent_000(t *testing.T) {
	if nullevent := event.NullEvent; nullevent == nil {
		t.Fatal("nullevent == nil")
	}
}
