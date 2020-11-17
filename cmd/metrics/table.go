package main

import (
	"sync"
)

type table struct {
	sync.Mutex

	headers bool
	columns []string
	types   []types
}

type types struct {
	t map[string]bool
}

func NewTable(headers bool) *table {
	this := new(table)
	this.headers = headers
	return this
}

func NewTypes() *types {
	this := new(types)
	this.t = make(map[string]bool)
	return this
}

func (this *table) Scan(row []string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	//fmt.Println(row)

}
