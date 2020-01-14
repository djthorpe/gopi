/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "gopi.LIRCEvent", Handler: LIRCHandler},
	}
	mutex  sync.Mutex
	header bool
)

func LIRCHandler(_ context.Context, _ gopi.App, evt_ gopi.Event) {
	mutex.Lock()
	defer mutex.Unlock()
	if header == false {
		fmt.Printf("\n%-10s %-10s %s\n", "MODE", "TYPE", "VALUE")
		fmt.Printf("%-10s %-10s %s\n", strings.Repeat("-", 10), strings.Repeat("-", 10), strings.Repeat("-", 20))
		header = true
	}
	evt := evt_.(gopi.LIRCEvent)
	fmt.Printf("%-10s %-10s %v\n",
		strings.ToLower(strings.TrimPrefix(fmt.Sprint(evt.Mode()), "LIRC_MODE_")),
		strings.ToLower(strings.TrimPrefix(fmt.Sprint(evt.Type()), "LIRC_TYPE_")),
		evt.Value())

	// Print header again after timeout
	if evt.Type() == gopi.LIRC_TYPE_TIMEOUT {
		header = false
	}
}
