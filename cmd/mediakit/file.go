package main

import (
	"path/filepath"

	"github.com/djthorpe/gopi/v3"
)

type file struct {
	Path  string
	Name  string
	Flags gopi.MediaFlag
	//	Metadata map[gopi.MediaKey]interface{}
	Streams []gopi.MediaStream
}

func NewFile(path string) *file {
	this := new(file)
	this.Path = path
	this.Name = filepath.Base(path)
	return this
}
