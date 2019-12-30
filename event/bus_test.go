/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package event_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	event "github.com/djthorpe/gopi/v2/event"
)

func Test_Bus_000(t *testing.T) {
	t.Log("Test_000")
}

func Test_Bus_001(t *testing.T) {
	if bus, err := gopi.New(event.Bus{}, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(bus)
	}
}
