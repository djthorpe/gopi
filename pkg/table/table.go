package table

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// Table represents a data table with a header of columns
type Table struct {
	cols          []cell
	types         []Types
	fields        map[string]int
	rows          [][]interface{}
	header        bool
	offset, limit uint
	merge         bool
}

// Formatter interface converts internal format to
// string
type Formatter interface {
	Format() (string, Alignment, Color)
}

type Option func(*Table)

/////////////////////////////////////////////////////////////////////
// CONSTANTS

type cellFormat int

const (
	cellAscii cellFormat = iota
	cellCsv
)

var (
	formatterType = reflect.TypeOf((*Formatter)(nil)).Elem()
)

/////////////////////////////////////////////////////////////////////
// NEW

// New creates an empty table
func New(opts ...Option) *Table {
	t := new(Table)
	t.header = true
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// ReadCSV creates a table from a CSV file optionally with a header
func ReadCSV(r io.Reader, opts ...Option) (*Table, error) {
	t := New(opts...)
	decoder := csv.NewReader(r)
	for {
		if row, err := decoder.Read(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if t.header {
			if row_ := t.decodeLine(row); len(row_) > 0 {
				t.SetHeader(row_...)
			}
			t.header = false
		} else {
			if row_ := t.decodeLine(row); len(row_) > 0 {
				t.Append(row_...)
			}
		}
	}

	// Return success
	return t, nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Append a row of data into the table
func (t *Table) Append(row ...interface{}) {
	row_ := make([]interface{}, len(row))
	for i := range row {
		// Append a new column and type
		if len(t.cols) <= i {
			t.cols = append(t.cols, cell{fmt.Sprintf("col_%02d", i), Auto, None})
			t.types = append(t.types, make(Types))
		}

		// Set the value
		row_[i] = row[i]

		// Parse value to determine the column kind
		t.types[i].Parse(row[i])
	}

	// Append row
	t.rows = append(t.rows, row_)

	// Map column names to indexes
	t.fields = make(map[string]int, len(t.cols))
	for i, v := range t.cols {
		t.fields[v.value] = i
	}
}

// Add a record from a map into the table, adding new
// columns as necessary
func (t *Table) Add(row map[string]interface{}) {
	row_ := make([]interface{}, len(row)+len(t.cols))
	if t.fields == nil {
		t.fields = make(map[string]int, len(row))
	}
	for k, v := range row {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if _, exists := t.fields[k]; exists == false {
			t.fields[k] = len(t.cols)
			t.cols = append(t.cols, cell{k, Auto, None})
			t.types = append(t.types, make(Types))
		}
		// Set the value and determine the column kind
		i := t.fields[k]
		row_[i] = v
		t.types[i].Parse(v)
	}
	t.rows = append(t.rows, row_[:len(t.cols)])
}

// Render table as ASCII to io.Writer
// Options WithHeader, WithMergeCells and WithOffsetLimit
// will affect the rendering of the table
func (t *Table) Render(w io.Writer, opts ...Option) {
	// Ignore if no cols or rows
	if len(t.cols) == 0 || len(t.rows) == 0 {
		return
	}
	// Set options
	for _, opt := range opts {
		opt(t)
	}
	// Write table
	table := tablewriter.NewWriter(w)
	if t.header {
		table.SetHeader(stringArray(t.cols))
	}

	// Set merge
	table.SetAutoMergeCells(t.merge)

	// Output rows
	c := uint(0)
	for i, row := range t.rows {
		if t.offset > 0 && uint(i) < t.offset {
			continue
		}
		if t.limit > 0 && c >= t.limit {
			continue
		}
		table.Append(stringArray(t.formatRow(row, cellAscii)))
		c++
	}

	// Set footer
	foot := []string{}
	for _, col := range t.types {
		foot = append(foot, fmt.Sprint(col.Kind()))
	}
	table.SetFooter(foot)
	table.Render()
}

// RenderCSV renders a table as CSV. Options WithHeader and
// WithOffsetLimit will affect the rendering of the table
func (t *Table) RenderCSV(w io.Writer, opts ...Option) {
	// Ignore if no cols or rows
	if len(t.cols) == 0 || len(t.rows) == 0 {
		return
	}
	// Set options
	t.header = true
	for _, opt := range opts {
		opt(t)
	}
	// Write table
	enc := csv.NewWriter(w)
	if t.header {
		enc.Write(stringArray(t.cols))
	}
	c := uint(0)
	for i, row := range t.rows {
		if t.offset > 0 && uint(i) < t.offset {
			continue
		}
		if t.limit > 0 && c >= t.limit {
			continue
		}
		enc.Write(stringArray(t.formatRow(row, cellCsv)))
		c++
	}
	enc.Flush()
}

func (t *Table) SetHeader(r ...interface{}) {
	t.cols = t.formatRow(r, cellAscii)
	// TODO: Preseve existing types
	t.types = make([]Types, len(t.cols))
	for i := range t.types {
		t.types[i] = make(Types)
	}
	t.fields = make(map[string]int, len(t.cols))
	for i, v := range t.cols {
		t.fields[v.value] = i
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (t *Table) formatRow(r []interface{}, f cellFormat) []cell {
	n := maxInt(len(t.cols), len(r))
	result := make([]cell, n)
	for i := 0; i < n; i++ {
		if i < len(r) {
			result[i] = t.formatCell(r[i], f)
		} else {
			result[i] = t.formatCell(nil, f)
		}
	}
	return result
}

func (t *Table) formatCell(v interface{}, f cellFormat) cell {
	if v == nil || isNil(reflect.ValueOf(v)) {
		return cell{formatNil(f), Auto, None}
	} else if f, ok := v.(Formatter); ok {
		v, a, c := f.Format()
		return cell{v, a, c}
	} else {
		return cell{fmt.Sprint(v), Auto, None}
	}
}

func (t *Table) decodeLine(row []string) []interface{} {
	result := make([]interface{}, len(row))
	for i, cell := range row {
		result[i] = cell
	}
	return result
}

func stringArray(row []cell) []string {
	result := make([]string, len(row))
	for i, cell := range row {
		result[i] = valueWithColor(cell.value, cell.color)
	}
	return result
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func isNil(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return rv.IsNil()
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	default:
		return false
	}
}

func formatNil(f cellFormat) string {
	switch f {
	case cellCsv:
		return ""
	default:
		return "<nil>"
	}
}
