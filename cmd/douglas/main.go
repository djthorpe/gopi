package main

import (
	"os"

	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

func main() {
	os.Exit(tool.CommandLine("douglas", os.Args[1:], new(app)))
}
