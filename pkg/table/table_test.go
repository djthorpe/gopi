package table_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/table"
)

const (
	CSV_PATH = "../../etc/csv/time_series_covid19_confirmed_global.csv"
)

func Test_Table_001(t *testing.T) {
	table := table.New()
	if table == nil {
		t.Fatal("Unexpected nil value")
	}
	var buf bytes.Buffer
	table.Render(&buf)
	if buf.Len() != 0 {
		t.Fatal("Expected no output")
	}
}

func Test_Table_002(t *testing.T) {
	table := table.New()
	table.Append("a", "b")
	table.Append("c")
	table.Append("d", "e", "f")

	var buf bytes.Buffer
	table.Render(&buf)
	t.Log(buf.String())
}

func Test_Table_003(t *testing.T) {
	table := table.New()
	table.Append("a", "b")
	table.Append("c")
	table.Append("d", "e", "f")

	var buf bytes.Buffer
	table.RenderCSV(&buf)
	t.Log(buf.String())
}

type cell struct{ string }

func (c cell) Format() (string, table.Alignment, table.Color) {
	return "[" + c.string + "]", table.Auto, table.Red
}

func Test_Table_004(t *testing.T) {
	table := table.New()
	table.Append(cell{"a"}, cell{"b"})
	table.Append(&cell{"c"})
	table.Append(&cell{"d"})

	table.SetHeader("COL A", "COL B", "COL C")

	var buf bytes.Buffer
	table.Render(&buf)
	t.Log(buf.String())
}

func Test_Table_005(t *testing.T) {
	fh, err := os.Open(CSV_PATH)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if table, err := table.ReadCSV(fh); err != nil {
		t.Fatal(err)
	} else {
		var buf bytes.Buffer
		table.Render(&buf)
		t.Log(buf.String())
	}

}

func Test_Table_006(t *testing.T) {
	v := table.New()
	v.Add(map[string]interface{}{
		"Col A": "A",
	})
	v.Add(map[string]interface{}{
		"Col B": "B",
	})
	v.Add(map[string]interface{}{
		"Col A": "A",
		"Col B": "B",
	})
	v.Add(map[string]interface{}{
		"Col C": "C",
	})

	var buf bytes.Buffer
	v.Render(&buf, table.WithHeader(false), table.WithOffsetLimit(1, 2))
	t.Log(buf.String())
}
