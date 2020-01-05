/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation https://gopi.mutablelogic.com/
  For Licensing and Usage information, please see LICENSE.md
*/

package timer

import (
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type evt struct {
	source    gopi.Driver
	info      *unit
	timestamp time.Time
}

////////////////////////////////////////////////////////////////////////////////
// EVENT INTERFACE

func NewTimerEvent(source gopi.Timer, u *unit, ts time.Time) gopi.Event {
	return &evt{source, u, ts}
}

func (this *evt) Name() string {
	return "TimerEvent"
}

func (this *evt) Source() gopi.Driver {
	return this.source
}

func (this *evt) Timestamp() time.Time {
	return this.timestamp
}

func (this *evt) UserInfo() interface{} {
	return this.info.userInfo
}

func (this *evt) String() string {
	return fmt.Sprintf("<sys.timer.event>{ ts=%v counter=%v userInfo=%v }", this.timestamp.Format(time.Kitchen), this.info.counter, this.info.userInfo)
}

func (this *evt) Counter() uint {
	return this.info.counter
}

func (this *evt) Cancel() {
	this.info.Cancel()
}
