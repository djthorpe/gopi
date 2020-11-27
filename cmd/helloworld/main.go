package main

import (
	"os"

	"github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/log"
)

func main() {
	os.Exit(tool.CommandLine("gaffer", os.Args[1:], new(app)))
}
