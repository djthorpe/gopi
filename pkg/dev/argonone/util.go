package argonone

import (
	"sort"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type fanConfigArr []struct {
	celcius float32
	fan     uint8
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
