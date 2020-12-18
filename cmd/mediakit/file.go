package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type media struct {
	URL     *url.URL
	Flags   gopi.MediaFlag
	Info    os.FileInfo
	Meta    map[gopi.MediaKey]interface{}
	Streams []*stream
}

type stream struct {
	Flags gopi.MediaFlag
	Index int
}

type flag struct {
	gopi.MediaFlag
}

type filesize struct {
	int64
}

type name struct {
	string
}

/////////////////////////////////////////////////////////////////////
// MEDIA

func NewMedia(m gopi.Media, info os.FileInfo) *media {
	this := new(media)
	this.URL = m.URL()
	this.Info = info
	this.Flags = m.Flags()
	this.Meta = make(map[gopi.MediaKey]interface{})
	return this
}

func (this *media) Dict() map[string]interface{} {
	dict := map[string]interface{}{
		"Name": name{this.Info.Name()},
		"Type": flag{this.Flags},
		"Size": filesize{this.Info.Size()},
	}
	for k, v := range this.Meta {
		dict[string(k)] = v
	}
	return dict
}

/////////////////////////////////////////////////////////////////////
// STREAM

func NewStream(s gopi.MediaStream) *stream {
	this := new(stream)
	this.Index = s.Index()
	this.Flags = s.Flags()
	return this
}

/////////////////////////////////////////////////////////////////////
// FORMATS

func (v flag) Format() (string, table.Alignment, table.Color) {
	str := ""
	for _, flag := range strings.Split(fmt.Sprint(v), "|") {
		flag := strings.TrimPrefix(flag, "MEDIA_FLAG_")
		str += strings.ToLower(flag) + "|"
	}
	return strings.TrimSuffix(str, "|"), table.Auto, table.Bold
}

func (v filesize) Format() (string, table.Alignment, table.Color) {
	str := ""

	const unit = 1024
	if v.int64 < unit {
		str = fmt.Sprintf("%vB", v.int64)
	} else {
		div, exp := int64(unit), 0
		for n := v.int64 / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
		str = fmt.Sprintf("%.1f%ciB", float64(v.int64)/float64(div), "KMGTPE"[exp])
	}
	return str, table.Auto, table.Bold
}

func (v name) Format() (string, table.Alignment, table.Color) {
	return v.string, table.Auto, table.Bold
}
