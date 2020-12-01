package keycode

import (
	"path/filepath"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type keycodedb struct {
	path string
	name string
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewKeycodeDatabase(path, name string) (*keycodedb, error) {
	this := new(keycodedb)
	this.name = name
	this.path = path
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *keycodedb) String() string {
	str := "<keycodedb"
	if this.path != "" {
		str += " filename=" + strconv.Quote(filepath.Base(this.path))
	}
	if this.name != "" {
		str += " name=" + strconv.Quote(this.name)
	}
	return str + ">"
}
