package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type table struct {
	headers bool
	columns []string
	types   []*types
	rows    []*row
}

type types struct {
	t map[string]bool
}

type row struct {
	values []string
}

var (
	order = []string{"bool", "uint", "int", "float", "string"}
)

func NewTable(headers bool) *table {
	this := new(table)
	this.headers = headers
	return this
}

func NewTypes() *types {
	this := new(types)
	this.t = make(map[string]bool, len(order))
	for _, t := range order {
		this.t[t] = true
	}
	return this
}

func NewRow(values []string) *row {
	this := new(row)
	this.values = values
	return this
}

func (this *table) Scan(row []string) bool {
	// Set up types
	if this.types == nil {
		for _ = range row {
			this.types = append(this.types, NewTypes())
		}
	}

	// Set up column headers
	if this.columns == nil {
		if this.headers {
			this.columns = row
			return false
		} else {
			for i := range row {
				this.columns = append(this.columns, fmt.Sprintf("column%02d", i))
			}
		}
	}

	// Scan cells and detect allowed types
	for i, cell := range row {
		if i < len(this.types) {
			this.types[i].Scan(cell)
		}
	}

	return true
}

func (this *table) Append(row []string) {
	// scan returns true if the row should be appended
	if this.Scan(row) {
		this.rows = append(this.rows, NewRow(row))
	}
}

func (this *table) Schema() *table {
	schema := &table{}
	schema.headers = this.headers
	schema.columns = this.columns
	schema.types = this.types

	types := make([]string, len(schema.types))
	for i, t := range schema.types {
		types[i] = fmt.Sprint(t.Type())
	}
	schema.Append(types)

	return schema
}

func (this *table) Write(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader(this.columns)
	table.SetAutoFormatHeaders(false)

	for _, row := range this.rows {
		table.Append(row.values)
	}

	table.Render()
	return nil
}

func (this *types) Scan(value string) {
	// Don't bother checking when only one type remains
	if len(this.t) <= 1 {
		return
	}
	// Trim space to check and remove invalid types
	value = strings.TrimSpace(value)
	if value != "" {
		for _, t := range order {
			if _, exists := this.t[t]; exists {
				if this.valid(value, t) == false {
					delete(this.t, t)
				}
			}
		}
	}
}

func (this *types) Type() string {
	for _, t := range order {
		if _, exists := this.t[t]; exists {
			return t
		}
	}
	return ""
}

func (this *types) valid(value, t string) bool {
	switch t {
	case "bool":
		if _, err := strconv.ParseBool(value); err != nil {
			return false
		}
	case "uint":
		if _, err := strconv.ParseUint(value, 0, 64); err != nil {
			return false
		}
	case "int":
		if _, err := strconv.ParseInt(value, 0, 64); err != nil {
			return false
		}
	case "float":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return false
		}
	}
	return true
}
