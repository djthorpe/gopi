// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap_test

import (
	"sync"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
)

func Test_RetainCount_000(t *testing.T) {
	t.Log("Test_RetainCount_000")
}

func Test_RetainCount_001(t *testing.T) {
	count := new(bitmap.RetainCount)
	count.Inc()
	if count.Dec() != true {
		t.Error("Expected true from Dec method")
	}
}

func Test_RetainCount_002(t *testing.T) {
	var wg sync.WaitGroup
	count := new(bitmap.RetainCount)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			count.Inc()
			wg.Done()
		}()
	}
	wg.Wait()
	if count.Inc() != 101 {
		t.Error("Expected Inc to return 101")
	}
}
