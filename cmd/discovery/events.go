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
	gopi "github.com/djthorpe/gopi/v2"
)

var (
	WriteHeader sync.Once
	WriteFormat = "%-8s %-25s %-20s %-35s\n"
)

////////////////////////////////////////////////////////////////////////////////
// EVENTS

func EventHandler(_ context.Context, _ gopi.App, evt gopi.Event) {
	WriteHeader.Do(func() {
		fmt.Printf(WriteFormat, "TYPE", "NAME", "SERVICE", "HOST")
		fmt.Printf(WriteFormat,
			strings.Repeat("-", 8),
			strings.Repeat("-", 25),
			strings.Repeat("-", 20),
			strings.Repeat("-", 35),
		)
	})
	evt_ := evt.(gopi.RPCEvent)
	type_ := strings.TrimPrefix(fmt.Sprint(evt_.Type()), "RPC_EVENT_SERVICE_")
	host_ := fmt.Sprintf("%s:%v", evt_.Service().Host, evt_.Service().Port)
	if evt_.Service().Port == 0 {
		host_ = ""
	}
	fmt.Printf(WriteFormat,
		TruncateString(strings.ToLower(type_), 8),
		TruncateString(evt_.Service().Name, 25),
		TruncateString(evt_.Service().Service, 20),
		TruncateString(host_, 35),
	)
}

func TruncateString(value string, l int) string {
	if len(value) > l {
		value = value[0:l-4] + "..."
	}
	return value
}
