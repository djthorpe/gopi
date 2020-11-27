package config

import (
	"runtime"
	"strconv"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	GitTag      string
	GitBranch   string
	GitHash     string
	GoBuildTime string
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type version struct {
	name string
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func NewVersion(name string) gopi.Version {
	return &version{name}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *version) Name() string {
	return this.name
}

func (*version) Version() (string, string, string) {
	return GitTag, GitBranch, GitHash
}

func (*version) GoVersion() string {
	return runtime.Version()
}

func (*version) BuildTime() time.Time {
	if GoBuildTime != "" {
		if t, err := time.Parse(time.RFC3339, GoBuildTime); err == nil {
			return t
		}
	}
	return time.Time{}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *version) String() string {
	str := "<version"
	str += " name=" + strconv.Quote(this.Name())
	tag, branch, hash := this.Version()
	if tag != "" {
		str += " tag=" + strconv.Quote(tag)
	}
	if branch != "" {
		str += " branch=" + strconv.Quote(branch)
	}
	if hash != "" {
		str += " hash=" + strconv.Quote(hash)
	}
	if goversion := this.GoVersion(); goversion != "" {
		str += " goversion=" + strconv.Quote(goversion)
	}
	if buildtime := this.BuildTime(); buildtime.IsZero() == false {
		str += " buildtime=" + strconv.Quote(buildtime.Format(time.RFC3339))
	}
	return str + ">"
}
