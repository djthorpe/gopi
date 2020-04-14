package base_test

import (
	"sync"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/base"
)

func TestPubSub_000(t *testing.T) {
	pub := &base.PubSub{}
	defer pub.Close()
	if ch := pub.Subscribe(); ch == nil {
		t.Error("Expected channel, got nil")
	} else {
		pub.Unsubscribe(ch)
	}
}

func TestPubSub_001(t *testing.T) {
	pub := &base.PubSub{}
	defer pub.Close()

	// Emitter with no subscriber
	var wg sync.WaitGroup
	wg.Add(1)
	go func(times int) {
		defer wg.Done()
		for i := 0; i < times; i++ {
			t.Log("Emit", i)
			pub.Emit(i)
		}
	}(100)
	wg.Wait()
}

func TestPubSub_002(t *testing.T) {
	pub := &base.PubSub{}
	times := 100

	// Emitter with one subscriber
	var wg1, wg2 sync.WaitGroup

	wg2.Add(1)
	wg1.Add(1)
	go func(times int) {
		defer wg2.Done()
		sub := pub.Subscribe()
		wg1.Done()
		i := 0
		// Continue until pub.CLose()
		for v := range sub {
			t.Log("Receive", v)
			i++
		}
		if times != i {
			t.Error("Expected to recive", times, "values, but got", i)
		}
		t.Log("Times", times)
	}(times)

	wg1.Wait() // Wait for subscribe before emitting
	wg2.Add(1)
	go func(times int) {
		defer wg2.Done()
		for i := 0; i < times; i++ {
			t.Log("Emit", i)
			pub.Emit(i)
		}
		pub.Close()
	}(times)
	wg2.Wait()
}

func TestPubSub_003(t *testing.T) {
	pub := &base.PubSub{}
	times := 100

	// Emitter with n subscribers
	var ready, done sync.WaitGroup

	for n := 0; n < 20; n++ {
		ready.Add(1)
		done.Add(1)
		go func(times int) {
			sub := pub.Subscribe()
			ready.Done()
			i := 0
			// Continue until pub.CLose()
			for v := range sub {
				t.Log("Receive", v)
				i++
			}
			if times != i {
				t.Error("Expected to recive", times, "values, but got", i)
			}
			done.Done()
		}(times)
	}

	ready.Wait() // Wait for all subscribes are subscribed

	for i := 0; i < times; i++ {
		t.Log("Emit", i)
		pub.Emit(i)
	}
	pub.Close()
	done.Wait() // Wait for all subscribers to end
}
