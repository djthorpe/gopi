package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

type fields struct {
	order  []gopi.MediaKey
	fields map[gopi.MediaKey]string
}

func NewFields() *fields {
	this := new(fields)
	this.fields = make(map[gopi.MediaKey]string)
	return this
}

func (this *fields) Add(key gopi.MediaKey) {
	if _, exists := this.fields[key]; exists == false {
		this.fields[key] = fmt.Sprint(key)
		this.order = append(this.order, key)
		sort.Sort(this)
	}
}

func (this *fields) Keys() []gopi.MediaKey {
	return this.order
}

func (this *fields) Len() int {
	return len(this.order)
}

func (this *fields) Less(i, j int) bool {
	ti := this.fields[this.order[i]]
	tj := this.fields[this.order[j]]
	return strings.Compare(ti, tj) < 0
}

func (this *fields) Swap(i, j int) {
	this.order[i], this.order[j] = this.order[j], this.order[i]
}

func (this *fields) Names() []string {
	result := []string{}
	for _, key := range this.order {
		result = append(result, this.fields[key])
	}
	return result
}
