package main

import (
	"fmt"
	"os"
	"path"
	"flag"
	"strings"
    "github.com/djthorpe/gopi/rpi"
)


func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,"Usage: %s <command>\n",path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	pi := rpi.New()
	defer pi.Terminate()
	fmt.Println(rpi.VCGenCmd(strings.Join(args," ")))
}
