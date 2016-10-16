package util

import (
	"os"
	"fmt"
)

type Logger interface {
	Print(v string)
}

type StderrLogger struct { }

func (this *StderrLogger) Print(v string) {
	fmt.Fprintln(os.Stderr,v)
}
