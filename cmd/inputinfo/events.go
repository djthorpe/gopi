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
	"sync"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

var (
	Format = "%-20v %-15v %v\n"
	Header sync.Once
)

func EventHandler(_ context.Context, _ gopi.App, evt_ gopi.Event) {
	evt := evt_.(gopi.InputEvent)
	Header.Do(func() {
		fmt.Printf(Format,"DEVICE","TYPE","INFO")
		fmt.Printf(Format,strings.Repeat("-", 20),strings.Repeat("-", 15),strings.Repeat("-", 40))
	})	
	n := TruncateString(evt.Device().Name(),20)
	t := strings.TrimPrefix(fmt.Sprint(evt.Type()),"INPUT_EVENT_")
	switch evt.Type() {
	case gopi.INPUT_EVENT_KEYPRESS,gopi.INPUT_EVENT_KEYRELEASE,gopi.INPUT_EVENT_KEYREPEAT:
		s := strings.TrimPrefix(fmt.Sprint(evt.KeyCode()),"KEYCODE_")
		if evt.ScanCode() != 0 {
			s += fmt.Sprintf(" (scancode=0x%08X)",evt.ScanCode())
		}
		fmt.Printf(Format,n,t,s)
	case gopi.INPUT_EVENT_ABSPOSITION:
		s := fmt.Sprintf("X=%.1f Y=%.1f",evt.Abs().X,evt.Abs().Y)
		if evt.KeyCode() != gopi.KEYCODE_NONE {
			s += fmt.Sprintf(" (keycode=%v)",evt.KeyCode())
		}
		fmt.Printf(Format,n,t,s)
	case gopi.INPUT_EVENT_RELPOSITION:
		fmt.Printf(Format,n,t,fmt.Sprintf("X=%.1f Y=%.1f",evt.Rel().X,evt.Rel().Y))
	}
	
}

func TruncateString(value string, l int) string {
	if len(value) > l {
		value = value[0:l-4] + "..."
	}
	return value
}
