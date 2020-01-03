/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Publisher struct {
	q map[uint][]chan interface{}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Publisher) NewChannels(queue uint) []chan interface{} {
	if this.q == nil {
		this.q = make(map[uint][]chan interface{})
	}
	if _, exists := this.q[queue]; exists == false {
		this.q[queue] = make([]chan interface{}, 0, 1)
	}
	return this.q[queue]
}

func (this *Publisher) Emit(queue uint, value interface{}) {
	if chans := this.NewChannels(queue); chans != nil {
		for _, c := range chans {
			c <- value
		}
	}
}

func (this *Publisher) Subscribe(queue uint, capacity int) <-chan interface{} {
	if chans := this.NewChannels(queue); chans == nil {
		return nil
	} else {
		new := make(chan interface{}, capacity)
		this.q[queue] = append(chans, new)
		return new
	}
}

func (this *Publisher) Unsubscribe(c <-chan interface{}) bool {
	for queue, chans := range this.q {
		for i, other := range chans {
			if other == c {
				this.q[queue] = append(chans[:i], chans[i+1:]...)
				close(other)
				return true
			}
		}
	}
	return false
}

func (this *Publisher) Len(queue uint) int {
	return len(this.NewChannels(queue))
}
