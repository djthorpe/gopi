/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package dispmanx

import "sync"

type RetainCount struct {
	count uint
	sync.Mutex
}

func (this *RetainCount) Inc() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.count += 1
}

func (this *RetainCount) Dec() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.count -= 1
	return this.count == 0
}
