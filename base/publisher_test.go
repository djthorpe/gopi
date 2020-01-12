/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package base_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/base"
)

func Test_Publisher_000(t *testing.T) {
	t.Log("Test_Publisher_000")
}

func Test_Publisher_001(t *testing.T) {
	var this base.Publisher

	queue := uint(123)
	count1 := uint(0)
	count2 := uint(0)

	stop1 := this.Subscribe(queue, func(value interface{}) {
		if value.(uint) == queue {
			count1++
		}
		t.Log("RECEIVED1 = ", count1)
	})
	t.Log("SUB1 = ", stop1)

	stop2 := this.Subscribe(queue, func(value interface{}) {
		if value.(uint) == queue {
			count2++
		}
		t.Log("RECEIVED2 = ", count2)
	})
	t.Log("SUB2 = ", stop2)

	t.Log("EMIT x 10")
	for i := 0; i < 10; i++ {
		this.Emit(queue, queue)
	}

	t.Log("CLOSE1")
	this.Unsubscribe(stop1)

	t.Log("EMIT x 15")
	for i := 0; i < 15; i++ {
		this.Emit(queue, queue)
	}

	t.Log("CLOSE2")
	this.Unsubscribe(stop2)

	t.Log("EMIT x 20")
	for i := 0; i < 20; i++ {
		this.Emit(queue, queue)
	}

	t.Log("CLOSE")
	this.Close()

	if count1 != 10 {
		t.Error("Unexpected value for count1", count1)
	}
	if count2 != 25 {
		t.Error("Unexpected value for count2", count2)
	}
}
