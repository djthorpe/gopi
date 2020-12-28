package main

import (
	"github.com/djthorpe/gopi/v3/pkg/table"
)

type header struct {
	string
}

func (h header) Format() (string, table.Alignment, table.Color) {
	return h.string, table.Auto, table.Bold
}
