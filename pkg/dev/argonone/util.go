package argonone

import (
	"fmt"
	"sort"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type fanConfigArr []struct {
	celcius float32
	fan     uint8
}

type fanValue struct {
	cur   uint8
	when  time.Time
	delay time.Duration
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - FAN CONFIG

func (arr fanConfigArr) Len() int {
	return len(arr)
}

func (arr fanConfigArr) Swap(i, j int) {
	arr[i].celcius, arr[j].celcius = arr[j].celcius, arr[i].celcius
	arr[i].fan, arr[j].fan = arr[j].fan, arr[i].fan
}

func (arr fanConfigArr) Less(i, j int) bool {
	return arr[i].celcius < arr[j].celcius
}

func (arr fanConfigArr) fanForTemperature(celcius float32) uint8 {
	sort.Sort(arr)
	fan := uint8(0)
	for _, v := range arr {
		if celcius < v.celcius {
			return fan
		} else {
			fan = v.fan
		}
	}
	return fan
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - HYSTERISIS

func NewFanValue(d time.Duration) *fanValue {
	this := new(fanValue)
	this.delay = d
	return this
}

func (v *fanValue) Set(new uint8) bool {
	changed := v.cur != new || v.when.IsZero()
	ago := time.Now().Sub(v.when)

	// Where changed, set new value
	if changed {
		v.when = time.Now()
	}

	// Revoke changed if "when" has not been reached
	if changed {
		changed = ago >= v.delay
		fmt.Println("old value=", v.cur, "new value=", new, "changed=", changed, "when ago=", ago)
	}

	// If changed, set new value
	if changed {
		v.cur = new
	}

	return changed
}
