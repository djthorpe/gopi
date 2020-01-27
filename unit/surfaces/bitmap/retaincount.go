/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// TYPE

type RetainCount struct {
	count uint
	sync.RWMutex
}

////////////////////////////////////////////////////////////////////////////////
// INCREMENT AND DECREMENT

func (this *RetainCount) Inc() uint {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.count += 1
	return this.count
}

func (this *RetainCount) Dec() bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if this.count == 0 {
		panic("Dec() would make count < 0")
	}
	this.count -= 1
	return this.count == 0
}

func (this *RetainCount) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return "<RetainCount " + fmt.Sprint(this.count) + ">"
}
