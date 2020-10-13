package main

import (
	"os"

	_ "github.com/djthorpe/gopi/v3/pkg/dev/waveshare"
	_ "github.com/djthorpe/gopi/v3/pkg/log"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

func main() {
	os.Exit(tool.CommandLine("douglas", os.Args[1:], new(app)))
}
