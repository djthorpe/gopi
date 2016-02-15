package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rjeczalik/notify"
)

type MediaItem struct {
	Path string
}

type MediaDatabase struct {
	fileRoot string
	files    map[string]MediaItem
}

func NewMediaDatabase(fileRoot string) *MediaDatabase {
	this := new(MediaDatabase)
	this.fileRoot = fileRoot
	this.files = make(map[string]MediaItem)
	return this
}

// TODO: process media items in the database

func processWorkItem(item WorkItem) {

	// skip hidden files
	filename := path.Base(item.Path)
	if strings.HasPrefix(filename, ".") {
		return
	}

	// check for file removal
	_, err := os.Stat(item.Path)
	if os.IsNotExist(err) && item.Event == notify.Remove {
		fmt.Println("Removed: ", item.Path)
		return
	}
	// check for file addition
	if item.Event == notify.Create && checkAllowedExtension(strings.ToLower(path.Ext(filename))) {
		fmt.Println("Created: ", item.Path)
		// TODO: add to the database
		return
	}
	// check for file rename
	if item.Event == notify.Rename && checkAllowedExtension(strings.ToLower(path.Ext(filename))) {
		fmt.Println("Renamed: ", item.Path)
		return
	}
	// check for file change
	if item.Event == notify.Write {
		fmt.Println("Modified: ", item.Path)
		return
	}

	fmt.Println("OTHER ", item.Path, item.Event)
}
