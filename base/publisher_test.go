/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package base_test

import (
	"sync"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/base"
)

func Test_Publisher_000(t *testing.T) {
	t.Log("Test_Publisher_000")
}

func Test_Publisher_001(t *testing.T) {
	this := new(base.Publisher)
	if c := this.SubscribeInt(0, 0); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else if cap(c) != 0 {
		t.Error("Unexpected return from cap()")
	} else {
		t.Log(c)
	}

	if c := this.SubscribeInt(0, 1); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else if cap(c) != 1 {
		t.Error("Unexpected return from cap()")
	} else {
		t.Log(c)
	}
}

func Test_Publisher_002(t *testing.T) {
	this := new(base.Publisher)
	if this.Len(0) != 0 {
		t.Error("Unexpected return from Len")
	} else if c := this.SubscribeInt(0, 0); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else if this.Len(0) != 1 {
		t.Error("Unexpected return from Len")
	} else if this.UnsubscribeInt(c) != true {
		t.Error("Unexpected return from Unsubscribe")
	} else if this.Len(0) != 0 {
		t.Error("Unexpected return from Len")
	}
}

func Test_Publisher_003(t *testing.T) {
	this := new(base.Publisher)
	queue := uint(100)
	chans := make([]<-chan interface{}, queue)

	for i := 0; i < len(chans); i++ {
		if chans[i] = this.SubscribeInt(queue, 0); chans[i] == nil {
			t.Error("Unexpected return from Subscribe()")
		}
	}
	if this.Len(queue) != int(queue) {
		t.Error("Unexpected return from Len")
	}
	for i := queue; i > 0; i-- {
		if this.UnsubscribeInt(chans[i-1]) != true {
			t.Error("Unexpected return from Unsubscribe()")
		}
	}
}

func Test_Publisher_004(t *testing.T) {
	this := new(base.Publisher)
	queue := uint(123)
	chans := make([]<-chan interface{}, queue)

	for i := 0; i < len(chans); i++ {
		if chans[i] = this.SubscribeInt(uint(i), 0); chans[i] == nil {
			t.Error("Unexpected return from Subscribe()")
		}
	}
	for i := queue; i > 0; i-- {
		if this.UnsubscribeInt(chans[i-1]) != true {
			t.Error("Unexpected return from Unsubscribe()")
		}
	}
}

func Test_Publisher_005(t *testing.T) {
	this := new(base.Publisher)
	queue := uint(123)

	if c := this.SubscribeInt(queue, 0); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else {
		var wait sync.WaitGroup
		wait.Add(1)
		go func() {
			value := <-c
			t.Log("Receive", value)
			wait.Done()
			if value.(uint) != queue {
				t.Error("Unexpected value received from channel")
			}
		}()
		t.Log("Emit")
		this.Emit(queue, queue)
		t.Log("Wait")
		wait.Wait()
		t.Log("Done")
	}
}

func Test_Publisher_006(t *testing.T) {
	this := new(base.Publisher)
	queue := uint(123)

	c1 := this.SubscribeInt(queue, 0)
	c2 := this.SubscribeInt(queue, 0)

	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		value := <-c1
		t.Log("C1 Receive", value)
		wait.Done()
		if value.(uint) != queue {
			t.Error("C1 Unexpected value received from channel")
		}
	}()
	go func() {
		value := <-c2
		t.Log("C2 Receive", value)
		wait.Done()
		if value.(uint) != queue {
			t.Error("C2 Unexpected value received from channel")
		}
	}()

	t.Log("Emit")
	this.Emit(queue, queue)
	t.Log("Wait")
	wait.Wait()
	t.Log("Done")
}

func Test_Publisher_007(t *testing.T) {
	this := new(base.Publisher)

	queue := uint(123)
	count1 := uint(0)
	count2 := uint(0)

	stop1 := this.Subscribe(queue, 0, func(value interface{}) {
		if value.(uint) == queue {
			count1++
		}
		t.Log("RECEIVED1 = ", count1)
	})

	stop2 := this.Subscribe(queue, 0, func(value interface{}) {
		if value.(uint) == queue {
			count2++
		}
		t.Log("RECEIVED2 = ", count2)
	})

	t.Log("EMIT x 10")
	for i := 0; i < 10; i++ {
		t.Log(queue, queue)
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
