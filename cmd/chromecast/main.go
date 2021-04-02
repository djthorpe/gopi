package main

import (
	"os"

	// Modules
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

func main() {
	os.Exit(tool.CommandLine("chromecast", os.Args[1:], new(app)))
}
