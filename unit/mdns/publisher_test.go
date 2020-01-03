/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns_test

import (
	"fmt"
	"sync"
	"testing"

	// Frameworks
	mdns "github.com/djthorpe/gopi/v2/unit/mdns"
)

func Test_Publisher_000(t *testing.T) {
	t.Log("Test_Publisher_000")
}

func Test_Publisher_001(t *testing.T) {
	this := new(mdns.Publisher)
	if c := this.Subscribe(0, 0); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else if cap(c) != 0 {
		t.Error("Unexpected return from cap()")
	} else {
		t.Log(c)
	}

	if c := this.Subscribe(0, 1); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else if cap(c) != 1 {
		t.Error("Unexpected return from cap()")
	} else {
		t.Log(c)
	}
}

func Test_Publisher_002(t *testing.T) {
	this := new(mdns.Publisher)
	if this.Len(0) != 0 {
		t.Error("Unexpected return from Len")
	} else if c := this.Subscribe(0, 0); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else if this.Len(0) != 1 {
		t.Error("Unexpected return from Len")
	} else if this.Unsubscribe(c) != true {
		t.Error("Unexpected return from Unsubscribe")
	} else if this.Len(0) != 0 {
		t.Error("Unexpected return from Len")
	}
}

func Test_Publisher_003(t *testing.T) {
	this := new(mdns.Publisher)
	queue := uint(100)
	chans := make([]<-chan interface{}, queue)

	for i := 0; i < len(chans); i++ {
		if chans[i] = this.Subscribe(queue, 0); chans[i] == nil {
			t.Error("Unexpected return from Subscribe()")
		}
	}
	if this.Len(queue) != int(queue) {
		t.Error("Unexpected return from Len")
	}
	for i := queue; i > 0; i-- {
		if this.Unsubscribe(chans[i-1]) != true {
			t.Error("Unexpected return from Unsubscribe()")
		}
	}
}

func Test_Publisher_004(t *testing.T) {
	this := new(mdns.Publisher)
	queue := uint(123)
	chans := make([]<-chan interface{}, queue)

	for i := 0; i < len(chans); i++ {
		if chans[i] = this.Subscribe(uint(i), 0); chans[i] == nil {
			t.Error("Unexpected return from Subscribe()")
		}
	}
	for i := queue; i > 0; i-- {
		if this.Unsubscribe(chans[i-1]) != true {
			t.Error("Unexpected return from Unsubscribe()")
		}
	}
}

func Test_Publisher_005(t *testing.T) {
	this := new(mdns.Publisher)
	queue := uint(123)

	if c := this.Subscribe(queue, 0); c == nil {
		t.Error("Unexpected return from Subscribe()")
	} else {
		var wait sync.WaitGroup
		wait.Add(1)
		go func() {
			value := <-c
			fmt.Println("Receive", value)
			wait.Done()
			if value.(uint) != queue {
				t.Error("Unexpected value received from channel")
			}
		}()
		fmt.Println("Emit")
		this.Emit(queue, queue)
		fmt.Println("Wait")
		wait.Wait()
		fmt.Println("Done")
	}
}

func Test_Publisher_006(t *testing.T) {
	this := new(mdns.Publisher)
	queue := uint(123)

	c1 := this.Subscribe(queue, 0)
	c2 := this.Subscribe(queue, 0)

	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		value := <-c1
		fmt.Println("C1 Receive", value)
		wait.Done()
		if value.(uint) != queue {
			t.Error("C1 Unexpected value received from channel")
		}
	}()
	go func() {
		value := <-c2
		fmt.Println("C2 Receive", value)
		wait.Done()
		if value.(uint) != queue {
			t.Error("C2 Unexpected value received from channel")
		}
	}()

	fmt.Println("Emit")
	this.Emit(queue, queue)
	fmt.Println("Wait")
	wait.Wait()
	fmt.Println("Done")
}
