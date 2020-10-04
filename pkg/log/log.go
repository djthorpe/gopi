package log

import (
	"fmt"
	"log"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

type Log struct {
	sync.Mutex
	gopi.Unit

	// Flags
	debug, verbose *bool
}

///////////////////////////////////////////////////////////////////////////////
// Implement gopi.Unit

func (this *Log) Define(cfg gopi.Config) error {
	this.debug = cfg.Bool("debug", false, "Set debugging flag")
	this.verbose = cfg.Bool("verbose", true, "Set verbose logging flag")
	return nil
}

func (*Log) New(gopi.Config) error {
	log.SetFlags(log.Ltime)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Implement gopi.Logger

func (this *Log) Print(args ...interface{}) {
	if this != nil {
		this.Lock()
		defer this.Unlock()
		log.Print(args...)
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Log) String() string {
	str := "<log"
	if this == nil {
		str += " nil"
	} else {
		if this.debug != nil {
			str += " debug=" + fmt.Sprint(*this.debug)
		}
		if this.verbose != nil {
			str += " verbose=" + fmt.Sprint(*this.verbose)
		}
	}
	return str + ">"
}
