/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi_test

import (
	"fmt"
	"strings"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

func Test_Event_000(t *testing.T) {
	t.Log("Test_Event_000")
}

func Test_Event_001(t *testing.T) {
	for ns := gopi.EVENT_NS_DEFAULT; ns <= gopi.EVENT_NS_MAX; ns++ {
		str := fmt.Sprint(ns)
		if strings.HasPrefix(str, "EVENT_NS_") == false {
			t.Error("Unexpected prefix", str)
		} else {
			t.Log(uint(ns), "=>", str)
		}
	}
}
