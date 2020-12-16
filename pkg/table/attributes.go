package table

import (
	"fmt"
	"strings"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Color uint16
type Alignment int

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	None Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const (
	Bold      Color = 0x0100
	Underline Color = 0x0200
	Inverse   Color = 0x0400
)

const (
	Auto Alignment = iota
	Left
	Center
	Right
)

const cESC = "\u001b"
const cSEP = ";"

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func valueWithColor(v string, c Color) string {
	if c == None {
		return v
	} else {
		return startFormat(c) + v + stopFormat()
	}
}

func startFormat(c Color) string {
	seq := []string{}
	seq = append(seq, fmt.Sprintf("%s[%d", cESC, int(c)&0x0F+29))
	if c&Bold == Bold {
		seq = append(seq, "1")
	}
	if c&Underline == Underline {
		seq = append(seq, "4")
	}
	if c&Inverse == Inverse {
		seq = append(seq, "7")
	}
	return strings.Join(seq, cSEP) + "m"
}

func stopFormat() string {
	return fmt.Sprintf("%s[%dm", cESC, 0)
}
