package log

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/djthorpe/gopi/v3"
)

type Log struct {
	sync.Mutex
	gopi.Unit

	// Flags
	debug *bool
	t     *testing.T
}

///////////////////////////////////////////////////////////////////////////////
// Implement gopi.Unit

func (this *Log) Define(cfg gopi.Config) error {
	this.debug = cfg.FlagBool("debug", false, "Set debugging flag")
	return nil
}

func (*Log) New(gopi.Config) error {
	log.SetFlags(log.Ltime)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Log) Print(args ...interface{}) {
	this.Lock()
	defer this.Unlock()
	if this.t != nil {
		this.t.Log(args...)
	} else {
		log.Print(args...)
	}
}

func (this *Log) IsDebug() bool {
	if this.debug == nil {
		return false
	} else {
		return *this.debug
	}
}

func (this *Log) Debug(args ...interface{}) {
	if this.IsDebug() {
		this.Lock()
		defer this.Unlock()
		if this.t != nil {
			this.t.Log(args...)
		} else {
			log.Print(args...)
		}
	}
}

func (this *Log) Printf(fmt string, args ...interface{}) {
	this.Lock()
	defer this.Unlock()
	if this.t != nil {
		this.t.Logf(fmt, args...)
	} else {
		log.Printf(fmt, args...)
	}
}

func (this *Log) Debugf(fmt string, args ...interface{}) {
	if this.IsDebug() {
		this.Lock()
		defer this.Unlock()
		if this.t != nil {
			this.t.Logf(fmt, args...)
		} else {
			log.Printf(fmt, args...)
		}
	}
}

func (this *Log) T() *testing.T {
	return this.t
}

func (this *Log) SetT(t *testing.T) {
	this.Lock()
	defer this.Unlock()
	this.t = t
	*this.debug = true
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
	}
	return str + ">"
}
